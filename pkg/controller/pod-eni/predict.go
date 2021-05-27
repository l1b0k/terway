package pod_eni

import (
	"strconv"

	"github.com/AliyunContainerService/terway/types"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type predicateForPodEvent struct {
	predicate.Funcs
}

func (p *predicateForPodEvent) Create(e event.CreateEvent) bool {
	return OKToProcess(e.Object)
}

func (p *predicateForPodEvent) Update(e event.UpdateEvent) bool {
	return OKToProcess(e.ObjectNew)
}

func (p *predicateForPodEvent) Delete(e event.DeleteEvent) bool {
	return OKToProcess(e.Object)
}

func (p *predicateForPodEvent) Generic(e event.GenericEvent) bool {
	return OKToProcess(e.Object)
}

// OKToProcess filter pod which is ready to process
func OKToProcess(obj interface{}) bool {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return false
	}
	if pod.Spec.NodeName == "" {
		return false
	}
	key, ok := pod.GetAnnotations()[types.PodENI]
	if !ok {
		return false
	}
	v, err := strconv.ParseBool(key)
	if err != nil {
		return false
	}
	return v
}
