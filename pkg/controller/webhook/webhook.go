package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AliyunContainerService/terway/deviceplugin"
	"github.com/AliyunContainerService/terway/pkg/apis/network.alibabacloud.com/v1beta1"
	"github.com/AliyunContainerService/terway/pkg/controller/common"
	"github.com/AliyunContainerService/terway/pkg/utils"

	"gomodules.xyz/jsonpatch/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var log = ctrl.Log.WithName("webhook")

func MutatingHook(client client.Client) *webhook.Admission {
	return &webhook.Admission{
		Handler: admission.HandlerFunc(func(ctx context.Context, req webhook.AdmissionRequest) webhook.AdmissionResponse {
			if req.Kind.Kind != "Pod" {
				return webhook.Allowed("none pod")
			}
			if req.Namespace == "kube-system" {
				return webhook.Allowed("namespace is kube-system")
			}

			original := req.Object.Raw

			pod := &corev1.Pod{}
			err := json.Unmarshal(original, pod)
			if err != nil {
				return admission.Errored(http.StatusInternalServerError, fmt.Errorf("failed decoding pod: %s, %w", string(original), err))
			}
			log.Info("pod", fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))

			if pod.Spec.HostNetwork {
				return webhook.Allowed("host network")
			}

			// 1. check pod with podNetworking config and get one
			podNetworkings := &v1beta1.PodNetworkingList{}
			err = client.List(ctx, podNetworkings)
			if err != nil {
				return webhook.Errored(1, fmt.Errorf("error list podNetworking, %w", err))
			}

			podNetworking, err := common.MatchOnePodNetworking(pod, podNetworkings.Items)
			if err != nil {
				log.Error(err, "error match podNetworking")
				return webhook.Errored(1, err)
			}
			if podNetworking == nil {
				log.V(4).Info("no selector is matched or CRD is not ready")
				return webhook.Allowed("not match")
			}

			if len(pod.Spec.Containers) == 0 {
				return webhook.Allowed("pod do not have containers")
			}
			// we only patch one container for res request
			if pod.Spec.Containers[0].Resources.Requests == nil {
				pod.Spec.Containers[0].Resources.Requests = make(corev1.ResourceList)
			}
			if pod.Spec.Containers[0].Resources.Limits == nil {
				pod.Spec.Containers[0].Resources.Limits = make(corev1.ResourceList)
			}
			pod.Spec.Containers[0].Resources.Requests[deviceplugin.MemberENIResName] = resource.MustParse("1")
			pod.Spec.Containers[0].Resources.Limits[deviceplugin.MemberENIResName] = resource.MustParse("1")

			// 2. get pod previous zone

			previousZone, err := getPreviousZone(client, pod)
			if err != nil {
				msg := fmt.Sprintf("error get previous podENI conf, %s", err)
				log.Error(err, msg)
				return webhook.Errored(1, fmt.Errorf(msg))
			}
			var zones []string
			if previousZone != "" {
				zones = []string{previousZone}
			} else {
				// 3. if no previous conf found, we will add zone limit by vSwitches
				for _, vsw := range podNetworking.Status.VSwitches {
					zones = append(zones, vsw.Zone)
				}
			}
			if len(zones) > 0 {
				pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = append(pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms,
					corev1.NodeSelectorTerm{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      corev1.LabelZoneFailureDomainStable,
								Operator: corev1.NodeSelectorOpIn,
								Values:   zones,
							},
						},
					})
			}
			podPatched, err := json.Marshal(pod)
			if err != nil {
				log.Error(err, "error marshal pod")
				return webhook.Errored(1, err)
			}
			patches, err := jsonpatch.CreatePatch(original, podPatched)
			if err != nil {
				log.Error(err, "error create patch")
				return webhook.Errored(1, err)
			}
			log.Info("patch pod for trunking")
			return webhook.Patched("ok", patches...)
		}),
	}
}

func getPreviousZone(client client.Client, pod *corev1.Pod) (string, error) {
	if !utils.IsStsPod(pod) {
		return "", nil
	}

	podENI := &v1beta1.PodENI{}
	err := client.Get(context.TODO(), k8stypes.NamespacedName{
		Namespace: pod.Namespace,
		Name:      pod.Name,
	}, podENI)
	if err != nil {
		if errors.IsNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return podENI.Spec.Allocation.ENI.Zone, nil
}
