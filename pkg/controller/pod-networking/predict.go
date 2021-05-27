package eni

import (
	"github.com/AliyunContainerService/terway/pkg/apis/network.alibabacloud.com/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type predicateForPodNetworkingEvent struct {
	predicate.Funcs
}

func (p *predicateForPodNetworkingEvent) Create(e event.CreateEvent) bool {
	return OKToProcess(e.Object)
}

func (p *predicateForPodNetworkingEvent) Update(e event.UpdateEvent) bool {
	return OKToProcess(e.ObjectNew)
}

func (p *predicateForPodNetworkingEvent) Delete(e event.DeleteEvent) bool {
	return OKToProcess(e.Object)
}

func (p *predicateForPodNetworkingEvent) Generic(e event.GenericEvent) bool {
	return OKToProcess(e.Object)
}

// OKToProcess filter pod which is ready to process
func OKToProcess(obj interface{}) bool {
	_, ok := obj.(*v1beta1.PodNetworking)
	if !ok {
		return false
	}
	return true
}
