package daemon

import (
	"context"

	"github.com/AliyunContainerService/terway/pkg/aliyun/client"
	"github.com/AliyunContainerService/terway/pkg/aliyun/instance"
	"github.com/AliyunContainerService/terway/pkg/k8s"
	"github.com/AliyunContainerService/terway/pkg/utils"
	"github.com/AliyunContainerService/terway/pkg/vswitch"
	"github.com/AliyunContainerService/terway/types"
	"github.com/AliyunContainerService/terway/types/daemon"
)

// getDynamicConfig returns (config, label, error) specified in node
// ("", "", nil) for no dynamic config for this node
func getDynamicConfig(ctx context.Context, k8s k8s.Kubernetes) (string, string, error) {
	label := k8s.GetNodeDynamicConfigLabel()
	if label == "" {
		return "", "", nil
	}

	cfg, err := k8s.GetDynamicConfigWithName(ctx, label)

	return cfg, label, err
}

func getENIConfig(cfg *daemon.Config) *types.ENIConfig {
	vswitchSelectionPolicy := vswitch.VSwitchSelectionPolicyRandom
	switch cfg.VSwitchSelectionPolicy {
	case "ordered":
		// keep the previous behave
		vswitchSelectionPolicy = vswitch.VSwitchSelectionPolicyMost
	}

	eniSelectionPolicy := types.EniSelectionPolicyMostIPs
	switch cfg.EniSelectionPolicy {
	case "least_ips":
		eniSelectionPolicy = types.EniSelectionPolicyLeastIPs
	}

	eniConfig := &types.ENIConfig{
		ZoneID:                 instance.GetInstanceMeta().ZoneID,
		VSwitchOptions:         nil,
		ENITags:                cfg.ENITags,
		SecurityGroupIDs:       cfg.GetSecurityGroups(),
		InstanceID:             instance.GetInstanceMeta().InstanceID,
		VSwitchSelectionPolicy: vswitchSelectionPolicy,
		EniSelectionPolicy:     eniSelectionPolicy,
		ResourceGroupID:        cfg.ResourceGroupID,
		EniTypeAttr:            0,
		TagFilter:              cfg.ENITagFilter,
	}

	if cfg.VSwitches != nil {
		zoneVswitchs, ok := cfg.VSwitches[eniConfig.ZoneID]
		if ok && len(zoneVswitchs) > 0 {
			eniConfig.VSwitchOptions = cfg.VSwitches[eniConfig.ZoneID]
		}
	}

	if len(eniConfig.VSwitchOptions) == 0 {
		eniConfig.VSwitchOptions = []string{instance.GetInstanceMeta().VSwitchID}
	}

	if cfg.EnableENITrunking {
		types.EnableFeature(&eniConfig.EniTypeAttr, types.FeatTrunk)
	}
	if cfg.EnableERDMA {
		types.EnableFeature(&eniConfig.EniTypeAttr, types.FeatERDMA)
	}

	return eniConfig
}

// the actual size for pool is minIdle and maxIdle
func getPoolConfig(cfg *daemon.Config, daemonMode string, limit *client.Limits) (*types.PoolConfig, error) {

	poolConfig := &types.PoolConfig{
		BatchSize: 10,
	}

	if cfg.ENITags == nil {
		cfg.ENITags = make(map[string]string)
	}
	cfg.ENITags[types.NetworkInterfaceTagCreatorKey] = types.NetworkInterfaceTagCreatorValue

	capacity := 0
	maxENI := 0
	maxMemberENI := 0

	switch daemonMode {
	case daemon.ModeENIMultiIP:
		maxENI = limit.Adapters
		maxENI = int(float64(maxENI)*cfg.EniCapRatio) + cfg.EniCapShift - 1

		// set max eni node can use
		if cfg.MaxENI > 0 && cfg.MaxENI < maxENI {
			maxENI = cfg.MaxENI
		}

		ipPerENI := limit.IPv4PerAdapter
		if utils.IsWindowsOS() {
			// NB(thxCode): don't assign the primary IP of one assistant eni.
			ipPerENI--
		}

		capacity = maxENI * ipPerENI
		if cfg.MaxPoolSize > capacity {
			poolConfig.MaxPoolSize = capacity
		} else {
			poolConfig.MaxPoolSize = cfg.MaxPoolSize
		}

		poolConfig.MinPoolSize = cfg.MinPoolSize

		if cfg.MinENI > 0 {
			poolConfig.MinPoolSize = cfg.MinENI * ipPerENI
		}
		if poolConfig.MinPoolSize > poolConfig.MaxPoolSize {
			poolConfig.MinPoolSize = poolConfig.MaxPoolSize
		}

		maxMemberENI = limit.MemberAdapterLimit

		poolConfig.MaxIPPerENI = ipPerENI

		if cfg.EnableERDMA {
			poolConfig.ERdmaCapacity = limit.ERDMARes() * limit.IPv4PerAdapter
		}
	}

	if cfg.IPAMType == types.IPAMTypeCRD {
		poolConfig.MaxPoolSize = 0
		poolConfig.MinPoolSize = 0
	}

	poolConfig.Capacity = capacity
	poolConfig.MaxENI = maxENI
	poolConfig.MaxMemberENI = maxMemberENI

	return poolConfig, nil
}
