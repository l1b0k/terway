package eni

import (
	"context"

	"github.com/AliyunContainerService/terway/pkg/aliyun"
	"github.com/AliyunContainerService/terway/pkg/apis/network.alibabacloud.com/v1beta1"
	register "github.com/AliyunContainerService/terway/pkg/controller"
	"github.com/AliyunContainerService/terway/pkg/controller/dynamic"
	"github.com/AliyunContainerService/terway/pkg/controller/vswitch"
	"github.com/AliyunContainerService/terway/pkg/utils"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const controllerName = "pod-networking"

func init() {
	register.Add(controllerName, func(mgr manager.Manager, aliyunClient *aliyun.OpenAPI, swPool *vswitch.SwitchPool, cfg dynamic.Interface) error {
		c, err := controller.New(controllerName, mgr, controller.Options{
			Reconciler:              NewReconcilePodNetworking(mgr, aliyunClient, swPool),
			MaxConcurrentReconciles: 1,
		})
		if err != nil {
			return err
		}

		return c.Watch(
			&source.Kind{
				Type: &v1beta1.PodNetworking{},
			},
			&handler.EnqueueRequestForObject{},
			&predicate.ResourceVersionChangedPredicate{},
			&predicateForPodNetworkingEvent{},
		)
	})
}

// ReconcilePodNetworking implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePodNetworking{}

// ReconcilePodNetworking reconciles a AutoRepair object
type ReconcilePodNetworking struct {
	cache        cache.Cache
	client       client.Client
	scheme       *runtime.Scheme
	aliyunClient *aliyun.OpenAPI
	swPool       *vswitch.SwitchPool

	//record event recorder
	record record.EventRecorder
}

// NewReconcilePodNetworking watch pod lifecycle events and sync to podENI resource
func NewReconcilePodNetworking(mgr manager.Manager, aliyunClient *aliyun.OpenAPI, swPool *vswitch.SwitchPool) *ReconcilePodNetworking {
	r := &ReconcilePodNetworking{
		client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		record:       mgr.GetEventRecorderFor("PodNetworking"),
		aliyunClient: aliyunClient,
		swPool:       swPool,
	}
	return r
}

// Reconcile podNetworking when user create or vSwitch fields changed
func (m *ReconcilePodNetworking) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Reconcile")

	old := &v1beta1.PodNetworking{}
	err := m.client.Get(context.TODO(), request.NamespacedName, old)
	if err != nil {
		if errors.IsNotFound(err) {
			l.Info("not found")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	var statusVSW []v1beta1.VSwitch

	for _, id := range old.Spec.VSwitchIDs {
		sw, err := m.swPool.GetByID(id)
		if err != nil {
			continue
		}
		statusVSW = append(statusVSW, v1beta1.VSwitch{
			ID:   sw.ID,
			Zone: sw.Zone,
		})
	}
	//TODO check really changed

	update := old.DeepCopy()
	update.Status.UpdateAt = metav1.Now()
	update.Status.VSwitches = statusVSW
	update.Status.Status = v1beta1.NetworkingStatusReady

	return reconcile.Result{}, m.updateStatus(ctx, update, old)
}

// NeedLeaderElection need election
func (m *ReconcilePodNetworking) NeedLeaderElection() bool {
	return true
}

func (m *ReconcilePodNetworking) updateStatus(ctx context.Context, update, old *v1beta1.PodNetworking) error {
	err := wait.ExponentialBackoff(utils.DefaultPatchBackoff, func() (done bool, err error) {
		innerErr := m.client.Status().Patch(ctx, update, client.MergeFrom(old))
		if innerErr != nil {
			if errors.IsNotFound(innerErr) {
				l := log.FromContext(ctx)
				l.Info("podNetworking is not found")
				return true, nil
			}
			return false, err
		}
		return true, nil
	})
	return err
}
