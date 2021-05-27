package common

import (
	"fmt"

	"github.com/AliyunContainerService/terway/pkg/apis/network.alibabacloud.com/v1beta1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// MatchOnePodNetworking will range all podNetworking and try to found a match
func MatchOnePodNetworking(pod *corev1.Pod, networkings []v1beta1.PodNetworking) (*v1beta1.PodNetworking, error) {
	l := labels.Set(pod.Labels)
	for _, podNetworking := range networkings {
		if podNetworking.Status.Status != v1beta1.NetworkingStatusReady {
			continue
		}
		matchOne := false
		if podNetworking.Spec.Selector.PodSelector != nil {
			ok, err := PodMatchSelector(podNetworking.Spec.Selector.PodSelector, l)
			if err != nil {
				return nil, fmt.Errorf("error match pod selector, %w", err)
			}
			if !ok {
				continue
			}
			matchOne = true
		}
		if podNetworking.Spec.Selector.NamespaceSelector != nil {
			ok, err := PodMatchSelector(podNetworking.Spec.Selector.NamespaceSelector, l)
			if err != nil {
				return nil, fmt.Errorf("error match namespace selector, %w", err)
			}
			if !ok {
				continue
			}
			matchOne = true
		}
		if matchOne {
			return &podNetworking, nil
		}
	}
	return nil, nil
}

// PodMatchSelector pod is selected by selector
func PodMatchSelector(labelSelector *metav1.LabelSelector, l labels.Set) (bool, error) {
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return false, err
	}
	return selector.Matches(l), nil
}
