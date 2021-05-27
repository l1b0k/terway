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

package dynamic

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/AliyunContainerService/terway/pkg/aliyun"
	"github.com/AliyunContainerService/terway/pkg/apis/network.alibabacloud.com/v1beta1"
	"github.com/AliyunContainerService/terway/pkg/config"
	"github.com/AliyunContainerService/terway/pkg/controller/vswitch"
	"github.com/AliyunContainerService/terway/pkg/utils"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	terwayConfigmapNamespace string
	terwayConfigmapName      string
)

func init() {
	flag.StringVar(&terwayConfigmapName, "terway-configmap-name", "eni-config", "read terway config from configmap")
	flag.StringVar(&terwayConfigmapNamespace, "terway-configmap-namespace", "kube-system", "read terway config from configmap")
}

// ReconcileConfig reconciles cfg
type ReconcileConfig struct {
	logger       logr.Logger
	k8scs        kubernetes.Interface
	client       client.Client
	aliyunClient *aliyun.OpenAPI
	swPool       *vswitch.SwitchPool

	terwayConfigmapName      string
	terwayConfigmapNamespace string

	// muetx is protect fields below
	sync.RWMutex

	terwayConfig *config.Configure

	clusterID string
	vpcID     string
}

// NewReconcileConfig sync terway config
func NewReconcileConfig(mgr manager.Manager, aliyunClient *aliyun.OpenAPI, swPool *vswitch.SwitchPool) (*ReconcileConfig, error) {
	r := &ReconcileConfig{
		logger:       mgr.GetLogger(),
		client:       mgr.GetClient(),
		k8scs:        utils.K8sClient,
		aliyunClient: aliyunClient,
		swPool:       swPool,

		terwayConfigmapName:      terwayConfigmapName,
		terwayConfigmapNamespace: terwayConfigmapNamespace,

		clusterID: os.Getenv("ALIYUN_CLUSTER_ID"),
		vpcID:     os.Getenv("ALIYUN_VPC_ID"),
	}
	cm, err := utils.K8sClient.CoreV1().ConfigMaps(terwayConfigmapNamespace).Get(context.TODO(), terwayConfigmapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	podNetworkingList, err := utils.NetworkClient.NetworkV1beta1().PodNetworkings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	err = r.sync(cm, podNetworkingList.Items)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// sync sync with conf
func (c *ReconcileConfig) sync(cm *corev1.ConfigMap, networkings []v1beta1.PodNetworking) error {
	eniConf, err := c.parseTerwayConfig(cm)
	if err != nil {
		return err
	}
	for _, ids := range eniConf.VSwitches {
		for _, id := range ids {
			_, err = c.swPool.GetByID(id)
			if err != nil {
				c.logger.Error(err, fmt.Sprintf("error get vSwitch %s", id))
			}
		}
	}

	// sync from terway podNetworking
	for _, podNetworking := range networkings {
		for _, vsw := range podNetworking.Status.VSwitches {
			_, err = c.swPool.GetByID(vsw.ID)
			if err != nil {
				c.logger.Error(err, "error get vSwitch")
			}
		}
	}

	c.Lock()
	defer c.Unlock()
	c.terwayConfig = eniConf

	return nil
}

// Start run and never stop
func (c *ReconcileConfig) Start(ctx context.Context) error {
	wait.Until(func() {
		// sync from terway configmap
		cm := &corev1.ConfigMap{}
		err := c.client.Get(context.TODO(), k8stypes.NamespacedName{
			Namespace: c.terwayConfigmapNamespace,
			Name:      c.terwayConfigmapName,
		}, cm)
		if err != nil {
			c.logger.Error(err, "error get terway configmap")
			return
		}
		podNetworkingList := &v1beta1.PodNetworkingList{}
		err = c.client.List(context.TODO(), podNetworkingList)
		if err != nil {
			c.logger.Error(err, "error get podNetworking")
			return
		}

		err = c.sync(cm, podNetworkingList.Items)
		if err != nil {
			c.logger.Error(err, "error sync")
		}
	}, time.Minute, ctx.Done())

	return fmt.Errorf("sync loop end")
}

// NeedLeaderElection need election
func (c *ReconcileConfig) NeedLeaderElection() bool {
	return true
}

// GetTerwayConfig from local
func (c *ReconcileConfig) GetTerwayConfig() *config.Configure {
	c.RLock()
	defer c.RUnlock()

	return c.terwayConfig
}

func (c *ReconcileConfig) GetClusterID() string {
	return c.clusterID
}

func (c *ReconcileConfig) GetVPCID() string {
	return c.vpcID
}

func (c *ReconcileConfig) parseTerwayConfig(cm *corev1.ConfigMap) (*config.Configure, error) {
	eniConfStr, ok := cm.Data["eni_conf"]
	if !ok {
		return nil, fmt.Errorf("error get eni_conf field")
	}

	eniConf, err := config.MergeConfigAndUnmarshal(nil, []byte(eniConfStr))
	if err != nil {
		return nil, err
	}
	return eniConf, nil
}

// ReconcilePod implements Interface
var _ Interface = &ReconcileConfig{}

type Interface interface {
	GetTerwayConfig() *config.Configure
	GetClusterID() string
	GetVPCID() string
}
