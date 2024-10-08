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

package pod

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/AliyunContainerService/terway/pkg/utils"
	"github.com/AliyunContainerService/terway/types"
)

type predicateForPodEvent struct {
	predicate.Funcs
}

func (p *predicateForPodEvent) Create(e event.CreateEvent) bool {
	return processPod(e.Object)
}

func (p *predicateForPodEvent) Update(e event.UpdateEvent) bool {
	return processPod(e.ObjectNew)
}

func (p *predicateForPodEvent) Delete(e event.DeleteEvent) bool {
	return processPod(e.Object)
}

func (p *predicateForPodEvent) Generic(e event.GenericEvent) bool {
	return processPod(e.Object)
}

func processPod(o client.Object) bool {
	pod, ok := o.(*corev1.Pod)
	if !ok {
		return false
	}

	if pod.Spec.NodeName == "" ||
		pod.Spec.HostNetwork {
		return false
	}

	if types.IgnoredByTerway(pod.Labels) {
		return false
	}
	return true
}

func processNode(node *corev1.Node) bool {
	if types.IgnoredByTerway(node.Labels) ||
		utils.ISVKNode(node) {
		return false
	}

	return true
}
