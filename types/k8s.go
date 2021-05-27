package types

// AnnotationPrefix is the annotation prefix
const AnnotationPrefix = "k8s.aliyun.com/"

const (
	// TrunkOn is the key for eni
	TrunkOn = AnnotationPrefix + "trunk-on"

	// PodENI whether pod is using eni (trunking mode)
	PodENI = AnnotationPrefix + "pod-eni"
	VPCID  = AnnotationPrefix + "vpc-id"
)
