/*
Copyright 2021-2022 Terway Authors.

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
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	wh "sigs.k8s.io/controller-runtime/pkg/webhook"

	aliyun "github.com/AliyunContainerService/terway/pkg/aliyun/client"
	"github.com/AliyunContainerService/terway/pkg/aliyun/credential"
	"github.com/AliyunContainerService/terway/pkg/apis/crds"
	networkv1beta1 "github.com/AliyunContainerService/terway/pkg/apis/network.alibabacloud.com/v1beta1"
	"github.com/AliyunContainerService/terway/pkg/backoff"
	"github.com/AliyunContainerService/terway/pkg/cert"
	register "github.com/AliyunContainerService/terway/pkg/controller"
	_ "github.com/AliyunContainerService/terway/pkg/controller/all"
	"github.com/AliyunContainerService/terway/pkg/controller/webhook"
	"github.com/AliyunContainerService/terway/pkg/metric"
	"github.com/AliyunContainerService/terway/pkg/utils"
	"github.com/AliyunContainerService/terway/pkg/utils/k8sclient"
	"github.com/AliyunContainerService/terway/pkg/version"
	"github.com/AliyunContainerService/terway/pkg/vswitch"
	"github.com/AliyunContainerService/terway/types/controlplane"
)

var (
	scheme = runtime.NewScheme()
	log    = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
	utilruntime.Must(networkv1beta1.AddToScheme(scheme))

	metrics.Registry.MustRegister(metric.OpenAPILatency)
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var (
		configFilePath     string
		credentialFilePath string
	)
	flag.StringVar(&configFilePath, "config", "/etc/config/ctrl-config.yaml", "config file for controlplane")
	flag.StringVar(&credentialFilePath, "credential", "/etc/credential/ctrl-secret.yaml", "secret file for controlplane")
	klog.InitFlags(nil)
	defer klog.Flush()

	flag.Parse()

	ctrl.SetLogger(klogr.New())
	log.Info(version.Version)

	ctx := ctrl.SetupSignalHandler()

	cfg, err := controlplane.ParseAndValidate(configFilePath, credentialFilePath)
	if err != nil {
		panic(err)
	}
	backoff.OverrideBackoff(cfg.BackoffOverride)
	utils.SetStsKinds(cfg.CustomStatefulWorkloadKinds)

	log.Info("using config", "config", cfg)

	restConfig := ctrl.GetConfigOrDie()
	restConfig.QPS = cfg.KubeClientQPS
	restConfig.Burst = cfg.KubeClientBurst
	restConfig.UserAgent = version.UA
	k8sclient.RegisterClients(restConfig)

	directClient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		log.Error(err, "unable to init k8s client")
		os.Exit(1)
	}

	lo.ForEach([]string{crds.CRDPodENI, crds.CRDPodNetworking, crds.CRDNode}, func(item string, index int) {
		err = crds.CreateOrUpdateCRD(ctx, directClient, item)
		if err != nil {
			log.Error(err, "unable sync crd")
			os.Exit(1)
		}
	})

	ws := wh.NewServer(wh.Options{
		Port:    cfg.WebhookPort,
		CertDir: cfg.CertDir,
	})
	options := ctrl.Options{
		Scheme:                     scheme,
		HealthProbeBindAddress:     cfg.HealthzBindAddress,
		WebhookServer:              ws,
		LeaderElection:             cfg.LeaderElection,
		LeaderElectionID:           cfg.ControllerName,
		LeaderElectionNamespace:    cfg.ControllerNamespace,
		LeaderElectionResourceLock: "leases",
		MetricsBindAddress:         cfg.MetricsBindAddress,
		Cache: cache.Options{
			ByObject: map[client.Object]cache.ByObject{
				&corev1.Node{}: {
					Transform: func(i interface{}) (interface{}, error) {
						if node, ok := i.(*corev1.Node); ok {
							node.Status.Images = nil
							node.Status.VolumesInUse = nil
							node.Status.VolumesAttached = nil
							return node, nil
						}
						return nil, fmt.Errorf("unexpected type %T", i)
					},
				},
				&corev1.Pod{}: {
					Transform: func(i interface{}) (interface{}, error) {
						if pod, ok := i.(*corev1.Pod); ok {
							pod.Spec.Volumes = nil
							pod.Spec.EphemeralContainers = nil
							pod.Spec.SecurityContext = nil
							pod.Spec.ImagePullSecrets = nil
							pod.Spec.Tolerations = nil
							pod.Spec.ReadinessGates = nil
							pod.Spec.PreemptionPolicy = nil
							pod.Status.InitContainerStatuses = nil
							pod.Status.ContainerStatuses = nil
							pod.Status.EphemeralContainerStatuses = nil
							return pod, nil
						}
						return nil, fmt.Errorf("unexpected type %T", i)
					},
				},
			},
		},
	}

	if !cfg.DisableWebhook {
		err = cert.SyncCert(ctx, directClient, cfg.ControllerNamespace, cfg.ControllerName, cfg.ClusterDomain, cfg.CertDir)
		if err != nil {
			panic(err)
		}
	}

	var providers []credential.Interface
	if string(cfg.Credential.AccessKey) != "" && string(cfg.Credential.AccessSecret) != "" {
		providers = append(providers, credential.NewAKPairProvider(string(cfg.Credential.AccessKey), string(cfg.Credential.AccessSecret)))
	}
	providers = append(providers, credential.NewEncryptedCredentialProvider(cfg.CredentialPath, cfg.SecretNamespace, cfg.SecretName))
	providers = append(providers, credential.NewMetadataProvider())

	clientSet, err := credential.NewClientMgr(cfg.RegionID, providers...)
	if err != nil {
		panic(err)
	}

	aliyunClient, err := aliyun.New(clientSet, flowcontrol.NewTokenBucketRateLimiter(cfg.ReadOnlyQPS, cfg.ReadOnlyBurst), flowcontrol.NewTokenBucketRateLimiter(cfg.MutatingQPS, cfg.MutatingBurst))
	if err != nil {
		panic(err)
	}

	mgr, err := ctrl.NewManager(restConfig, options)
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

	if !cfg.DisableWebhook {
		mgr.GetWebhookServer().Register("/mutating", webhook.MutatingHook(mgr.GetClient()))
		mgr.GetWebhookServer().Register("/validate", webhook.ValidateHook())
	}

	vSwitchCtrl, err := vswitch.NewSwitchPool(cfg.VSwitchPoolSize, cfg.VSwitchCacheTTL)
	if err != nil {
		panic(err)
	}

	tp := oteltrace.NewNoopTracerProvider()
	if cfg.EnableTrace {
		grpcTP, err := initOpenTelemetry(ctx, "terway-controlplane", version.Version, cfg)
		if err != nil {
			panic(err)
		}
		defer func(grpcTP *trace.TracerProvider, ctx context.Context) {
			err := grpcTP.Shutdown(ctx)
			if err != nil {
				log.Error(err, "failed to shutdown grpc tracer")
			}
		}(grpcTP, ctx)
		tp = grpcTP
	}
	wg := &wait.Group{}
	ctrlCtx := &register.ControllerCtx{
		Context:        ctx,
		Config:         cfg,
		VSwitchPool:    vSwitchCtrl,
		AliyunClient:   aliyunClient,
		Wg:             wg,
		TracerProvider: tp,
	}

	for name := range register.Controllers {
		if controlplane.IsControllerEnabled(name, register.Controllers[name].Enable, cfg.Controllers) {
			err = register.Controllers[name].Creator(mgr, ctrlCtx)
			if err != nil {
				panic(err)
			}
			log.Info("register controller", "controller", name)
		}
	}

	log.Info("controller started")
	err = mgr.Start(ctx)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}

// initOpenTelemetry bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func initOpenTelemetry(ctx context.Context, serviceName, serviceVersion string, cfg *controlplane.Config) (*trace.TracerProvider, error) {
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			semconv.K8SNodeNameKey.String(os.Getenv("K8S_NODE_NAME")),
		))
	if err != nil {
		return nil, err
	}

	// Set up trace provider.
	headers := map[string]string{"Authentication": string(cfg.OtelToken)}
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.OtelEndpoint),
		otlptracegrpc.WithHeaders(headers))

	traceExporter, err := otlptrace.New(ctx, traceClient)

	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(5*time.Second)),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(traceProvider)

	return traceProvider, nil
}
