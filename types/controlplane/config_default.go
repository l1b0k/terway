//go:build default_build

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

package controlplane

import (
	"github.com/AliyunContainerService/terway/pkg/backoff"
	"github.com/AliyunContainerService/terway/types/secret"
)

type Config struct {
	// controller config
	LeaseLockName      string `json:"leaseLockName" validate:"required" mod:"default=terway-controller-lock"`
	LeaseLockNamespace string `json:"leaseLockNamespace" validate:"required" mod:"default=kube-system"`
	LeaseDuration      string `json:"leaseDuration"`
	RenewDeadline      string `json:"renewDeadline"`
	RetryPeriod        string `json:"retryPeriod"`

	ControllerNamespace string `json:"controllerNamespace" validate:"required" mod:"default=kube-system"`
	ControllerName      string `json:"controllerName" validate:"required" mod:"default=terway-controlplane"`

	HealthzBindAddress string `json:"healthzBindAddress" validate:"required,tcp_addr" mod:"default=0.0.0.0:80"`
	MetricsBindAddress string `json:"metricsBindAddress" validate:"required" mod:"default=0"`
	ClusterDomain      string `json:"clusterDomain" validate:"required" mod:"default=cluster.local"`
	DisableWebhook     bool   `json:"disableWebhook"`
	WebhookPort        int    `json:"webhookPort" validate:"gt=0,lte=65535" mod:"default=4443"`
	CertDir            string `json:"certDir" validate:"required" mod:"default=/var/run/webhook-cert"`
	LeaderElection     bool   `json:"leaderElection"`
	EnableTrace        bool   `json:"enableTrace"`

	PodMaxConcurrent    int `json:"podMaxConcurrent" validate:"gt=0,lte=10000" mod:"default=10"`
	PodENIMaxConcurrent int `json:"podENIMaxConcurrent" validate:"gt=0,lte=10000" mod:"default=10"`
	NodeController
	MultiIPController
	ENIController

	Controllers []string `json:"controllers"`

	// cluster info for controlplane
	RegionID  string `json:"regionID" validate:"required"`
	ClusterID string `json:"clusterID" validate:"required"`
	VPCID     string `json:"vpcID" validate:"required"`

	EnableTrunk        *bool  `json:"enableTrunk,omitempty"`
	EnableDevicePlugin bool   `json:"enableDevicePlugin"`
	IPStack            string `json:"ipStack,omitempty" validate:"oneof=ipv4 ipv6 dual" mod:"default=ipv4"`

	EnableWebhookInjectResource *bool `json:"enableWebhookInjectResource,omitempty"`

	KubeClientQPS   float32 `json:"kubeClientQPS" validate:"gt=0,lte=10000" mod:"default=20"`
	KubeClientBurst int     `json:"kubeClientBurst" validate:"gt=0,lte=10000" mod:"default=30"`

	VSwitchPoolSize int    `json:"vSwitchPoolSize" validate:"gt=0" mod:"default=1000"`
	VSwitchCacheTTL string `json:"vSwitchCacheTTL" mod:"default=20m0s"`

	CustomStatefulWorkloadKinds []string `json:"customStatefulWorkloadKinds"`

	BackoffOverride map[string]backoff.ExtendedBackoff `json:"backoffOverride,omitempty"`
	IPAMType        string                             `json:"ipamType"`
	CentralizedIPAM bool                               `json:"centralizedIPAM,omitempty"`

	RateLimit map[string]int `json:"rateLimit"`

	Credential
}

type Credential struct {
	AccessKey      secret.Secret `json:"accessKey" validate:"required_with=AccessSecret"`
	AccessSecret   secret.Secret `json:"accessSecret" validate:"required_with=AccessKey"`
	CredentialPath string        `json:"credentialPath"`
	OtelEndpoint   string        `json:"otelEndpoint"`
	OtelToken      secret.Secret `json:"otelToken"`
}

type MultiIPController struct {
	MultiIPPodMaxConcurrent       int    `json:"multiIPPodMaxConcurrent" validate:"gt=0,lte=20000" mod:"default=500"`
	MultiIPNodeMaxConcurrent      int    `json:"multiIPNodeMaxConcurrent" validate:"gt=0,lte=20000" mod:"default=500"`
	MultiIPNodeSyncPeriod         string `json:"multiIPNodeSyncPeriod" mod:"default=12h"`
	MultiIPGCPeriod               string `json:"multiIPGCPeriod" mod:"default=2m"`
	MultiIPMinSyncPeriodOnFailure string `json:"multiIPMinSyncPeriodOnFailure" mod:"default=1s"`
	MultiIPMaxSyncPeriodOnFailure string `json:"multiIPMaxSyncPeriodOnFailure" mod:"default=300s"`
}

type NodeController struct {
	NodeMaxConcurrent  int               `json:"nodeMaxConcurrent" validate:"gt=0,lte=10000" mod:"default=10"`
	NodeLabelWhiteList map[string]string `json:"nodeLabelWhiteList"`
}

type ENIController struct {
	ENIMaxConcurrent int `json:"eniMaxConcurrent" validate:"gt=0,lte=10000" mod:"default=300"`
}
