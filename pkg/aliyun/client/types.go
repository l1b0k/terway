package client

import (
	"errors"
	"reflect"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/go-logr/logr"
)

var ErrInvalidArgs = errors.New("invalid args")

// log fields
const (
	LogFieldAPI              = "api"
	LogFieldRequestID        = "requestID"
	LogFieldInstanceID       = "instanceID"
	LogFieldSecondaryIPCount = "secondaryIPCount"
	LogFieldENIID            = "eni"
	LogFieldIPs              = "ips"
	LogFieldEIPID            = "eip"
	LogFieldPrivateIP        = "privateIP"
	LogFieldVSwitchID        = "vSwitchID"
	LogFieldSgID             = "securityGroupID"
	LogFieldResourceGroupID  = "resourceGroupID"
)

const (
	eniDescription    = "interface create by terway"
	maxSinglePageSize = 500
)

// status for eni
const (
	ENIStatusInUse     string = "InUse"
	ENIStatusAvailable string = "Available"
	ENIStatusAttaching string = "Attaching"
	ENIStatusDetaching string = "Detaching"
	ENIStatusDeleting  string = "Deleting"
)

const (
	ENITypePrimary   string = "Primary"
	ENITypeSecondary string = "Secondary"
	ENITypeTrunk     string = "Trunk"
	ENITypeMember    string = "Member"
)

const (
	ENITrafficModeRDMA     string = "HighPerformance"
	ENITrafficModeStandard string = "Standard"
)

const EIPInstanceTypeNetworkInterface = "NetworkInterface"

// NetworkInterface openAPI result for ecs.CreateNetworkInterfaceResponse and ecs.NetworkInterfaceSet
type NetworkInterface struct {
	Status             string             `json:"status,omitempty"`
	MacAddress         string             `json:"mac_address,omitempty"`
	NetworkInterfaceID string             `json:"network_interface_id,omitempty"`
	VSwitchID          string             `json:"v_switch_id,omitempty"`
	PrivateIPAddress   string             `json:"private_ip_address,omitempty"`
	PrivateIPSets      []ecs.PrivateIpSet `json:"private_ip_sets"`
	ZoneID             string             `json:"zone_id,omitempty"`
	SecurityGroupIDs   []string           `json:"security_group_ids,omitempty"`
	ResourceGroupID    string             `json:"resource_group_id,omitempty"`
	IPv6Set            []ecs.Ipv6Set      `json:"ipv6_set,omitempty"`
	Tags               []ecs.Tag          `json:"tags,omitempty"`

	// fields for DescribeNetworkInterface
	Type                        string `json:"type,omitempty"`
	InstanceID                  string `json:"instance_id,omitempty"`
	TrunkNetworkInterfaceID     string `json:"trunk_network_interface_id,omitempty"`
	NetworkInterfaceTrafficMode string `json:"network_interface_traffic_mode"`
	DeviceIndex                 int    `json:"device_index,omitempty"`
	CreationTime                string `json:"creation_time,omitempty"`
}

func FromCreateResp(in *ecs.CreateNetworkInterfaceResponse) *NetworkInterface {
	return &NetworkInterface{
		Status:             in.Status,
		MacAddress:         in.MacAddress,
		NetworkInterfaceID: in.NetworkInterfaceId,
		VSwitchID:          in.VSwitchId,
		PrivateIPAddress:   in.PrivateIpAddress,
		PrivateIPSets:      in.PrivateIpSets.PrivateIpSet,
		ZoneID:             in.ZoneId,
		SecurityGroupIDs:   in.SecurityGroupIds.SecurityGroupId,
		IPv6Set:            in.Ipv6Sets.Ipv6Set,
		Tags:               in.Tags.Tag,
		Type:               in.Type,
		ResourceGroupID:    in.ResourceGroupId,
	}
}

func FromDescribeResp(in *ecs.NetworkInterfaceSet) *NetworkInterface {
	ins := in.InstanceId
	if in.InstanceId == "" {
		ins = in.Attachment.InstanceId
	}

	return &NetworkInterface{
		Status:                      in.Status,
		MacAddress:                  in.MacAddress,
		NetworkInterfaceID:          in.NetworkInterfaceId,
		InstanceID:                  ins,
		VSwitchID:                   in.VSwitchId,
		PrivateIPAddress:            in.PrivateIpAddress,
		ZoneID:                      in.ZoneId,
		SecurityGroupIDs:            in.SecurityGroupIds.SecurityGroupId,
		IPv6Set:                     in.Ipv6Sets.Ipv6Set,
		PrivateIPSets:               in.PrivateIpSets.PrivateIpSet,
		Tags:                        in.Tags.Tag,
		TrunkNetworkInterfaceID:     in.Attachment.TrunkNetworkInterfaceId,
		NetworkInterfaceTrafficMode: in.NetworkInterfaceTrafficMode,
		DeviceIndex:                 in.Attachment.DeviceIndex,
		Type:                        in.Type,
		CreationTime:                in.CreationTime,
	}
}

// LogFields function enhances the provided logger with key-value pairs extracted from the fields of the given object.
//
// Parameters:
// l     - The original logr.Logger instance to be augmented with object field information.
// obj   - An arbitrary object whose fields will be inspected for logging. Must be of a struct type.
//
// Return Value:
// Returns an updated logr.Logger instance that includes key-value pairs for non-empty, non-zero fields of the input object.
// The original logger `l` is modified in place, and the returned logger is a reference to the same instance.
func LogFields(l logr.Logger, obj any) logr.Logger {
	r := l
	t := reflect.TypeOf(obj)

	realObj := obj
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		objValue := reflect.ValueOf(obj).Elem()
		realObj = objValue.Interface()
	}

	if t.Kind() == reflect.Struct {
		r = r.WithValues(LogFieldAPI, strings.TrimSuffix(t.Name(), "Request"))
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tagValue := field.Tag.Get("name")
		if tagValue == "" {
			continue
		}

		fieldValue := reflect.ValueOf(realObj).FieldByName(field.Name)
		if !fieldValue.IsValid() || fieldValue.IsZero() {
			continue
		}

		r = r.WithValues(field.Name, fieldValue.Interface())
	}
	return r
}
