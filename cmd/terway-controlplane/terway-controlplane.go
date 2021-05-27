/*
Copyright 2021 Terway Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	goflag "flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/AliyunContainerService/terway/pkg/aliyun"
	"github.com/AliyunContainerService/terway/pkg/apis/crds"
	networkv1beta1 "github.com/AliyunContainerService/terway/pkg/apis/network.alibabacloud.com/v1beta1"
	"github.com/AliyunContainerService/terway/pkg/cert"
	register "github.com/AliyunContainerService/terway/pkg/controller"
	_ "github.com/AliyunContainerService/terway/pkg/controller/all"
	cfg "github.com/AliyunContainerService/terway/pkg/controller/dynamic"
	"github.com/AliyunContainerService/terway/pkg/controller/vswitch"
	"github.com/AliyunContainerService/terway/pkg/controller/webhook"
	"github.com/AliyunContainerService/terway/pkg/utils"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/pkg/version"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

var (
	leaseLockName      string
	leaseLockNamespace string
	healthzBindAddress string

	scheme = runtime.NewScheme()
	log    = ctrl.Log.WithName("setup")
)

func init() {
	_ = flag.Set("v", "4")

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	flag.StringVar(&leaseLockName, "lease-lock-name", "terway-controller-lock", "the lease lock resource name")
	flag.StringVar(&leaseLockNamespace, "lease-lock-namespace", "kube-system", "the lease lock resource namespace")
	flag.StringVar(&healthzBindAddress, "healthzBindAddress", "0.0.0.0:80", "for health check")

	flag.String("cert-dir", "/var/run/webhook-cert", "webhook cert dir")
	flag.String("controller-namespace", "kube-system", "specific controller run namespace")
	flag.Int("webhook-port", 443, "port for webhook")

	_ = viper.BindPFlag("cert-dir", flag.CommandLine.Lookup("cert-dir"))
	_ = viper.BindPFlag("controller-namespace", flag.CommandLine.Lookup("controller-namespace"))
	_ = viper.BindPFlag("webhook-port", flag.CommandLine.Lookup("webhook-port"))

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(networkv1beta1.AddToScheme(scheme))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ctrl.SetLogger(klogr.New())
	log.Info(fmt.Sprintf("GitCommit %s BuildDate %s Platform %s",
		version.Get().GitCommit, version.Get().BuildDate, version.Get().Platform))

	flag.Parse()

	utils.RegisterClients()
	err := crds.RegisterCRDs()
	if err != nil {
		panic(err)
	}

	err = cert.SyncCert()
	if err != nil {
		panic(err)
	}

	ak := os.Getenv("ALICLOUD_ACCESS_KEY")
	sk := os.Getenv("ALICLOUD_SECRET_KEY")
	region := os.Getenv("ALICLOUD_REGION")
	region = "cn-hangzhou"

	aliyunClient, err := aliyun.NewAliyun(ak, sk, region, "")
	if err != nil {
		panic(err)
	}

	ctx := ctrl.SetupSignalHandler()
	restConfig := ctrl.GetConfigOrDie()

	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:                     scheme,
		HealthProbeBindAddress:     healthzBindAddress,
		Host:                       "0.0.0.0",
		Port:                       viper.GetInt("webhook-port"),
		CertDir:                    viper.GetString("cert-dir"),
		LeaderElection:             false,
		LeaderElectionID:           uuid.New().String() + "-" + os.Getenv("POD_NAME"),
		LeaderElectionNamespace:    viper.GetString("controller-namespace"),
		LeaderElectionResourceLock: "leases",
		MetricsBindAddress:         "0",
	})
	if err != nil {
		panic(err)
	}
	err = mgr.AddHealthzCheck("healthz", healthz.Ping)
	if err != nil {
		panic(err)
	}
	err = mgr.AddReadyzCheck("readyz", healthz.Ping)
	if err != nil {
		panic(err)
	}

	mgr.GetWebhookServer().Register("/mutating", webhook.MutatingHook(mgr.GetClient()))

	vSwitchCtrl, err := vswitch.NewSwitchPool(aliyunClient)
	if err != nil {
		panic(err)
	}
	err = mgr.Add(vSwitchCtrl)
	if err != nil {
		panic(err)
	}

	c, err := cfg.NewReconcileConfig(mgr, aliyunClient, vSwitchCtrl)
	if err != nil {
		panic(err)
	}

	err = mgr.Add(c)
	if err != nil {
		panic(err)
	}

	for name := range register.Controllers {
		err = register.Controllers[name](mgr, aliyunClient, vSwitchCtrl, c)
		if err != nil {
			panic(err)
		}
		log.Info("register controller", "controller", name)
	}

	log.Info("controller started")
	err = mgr.Start(ctx)
	if err != nil {
		panic(err)
	}
}
