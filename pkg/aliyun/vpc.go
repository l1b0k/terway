package aliyun

import (
	"fmt"
	"time"

	apiErr "github.com/AliyunContainerService/terway/pkg/aliyun/errors"
	"github.com/AliyunContainerService/terway/pkg/metric"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

// DescribeVSwitchByID get vsw by id
func (a *OpenAPI) DescribeVSwitchByID(vSwitch string) (*vpc.VSwitch, error) {
	req := vpc.CreateDescribeVSwitchesRequest()
	req.VSwitchId = vSwitch
	start := time.Now()
	resp, err := a.clientSet.VPC().DescribeVSwitches(req)
	metric.OpenAPILatency.WithLabelValues("DescribeVSwitches", fmt.Sprint(err != nil)).Observe(metric.MsSince(start))
	if err != nil {
		return nil, err
	}
	// For systems without RAM policy for VPC API permission, result is:
	// vsw is an empty slice, err is nil.
	// For systems which have RAM policy for VPC API permission,
	// (1) if vswitch indeed exists, result is:
	// vsw is a slice with a single element, err is nil.
	// (2) if vswitch doesn't exist, result is:
	// vsw is an empty slice, err is not nil.
	log.Debugf("result for DescribeVSwitches: vsw slice = %+v, err = %v", resp.VSwitches.VSwitch, err)
	if len(resp.VSwitches.VSwitch) == 0 {
		return nil, apiErr.ErrNotFound
	}
	if len(resp.VSwitches.VSwitch) > 0 {
		return &resp.VSwitches.VSwitch[0], nil
	}
	return nil, err
}
