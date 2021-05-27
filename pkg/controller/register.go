package register

import (
	"github.com/AliyunContainerService/terway/pkg/aliyun"
	"github.com/AliyunContainerService/terway/pkg/controller/dynamic"
	"github.com/AliyunContainerService/terway/pkg/controller/vswitch"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Creator func(mgr manager.Manager, aliyunClient *aliyun.OpenAPI, swPool *vswitch.SwitchPool, p dynamic.Interface) error

var Controllers = map[string]Creator{}

func Add(name string, creator Creator) {
	Controllers[name] = creator
}
