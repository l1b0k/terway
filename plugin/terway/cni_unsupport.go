//go:build !linux && !windows
// +build !linux,!windows

package main

import (
	"context"
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"

	"github.com/AliyunContainerService/terway/plugin/driver/types"
	"github.com/AliyunContainerService/terway/rpc"
	terwayTypes "github.com/AliyunContainerService/terway/types"
)

func getCmdArgs(args *skel.CmdArgs) (*cniCmdArgs, error) {
	panic("not implement")
}

type cniCmdArgs struct {
}

func (args *cniCmdArgs) GetCNIConf() *types.CNIConf {
	panic("not implement")
}

func (args *cniCmdArgs) GetK8SConfig() *types.K8SArgs {
	panic("not implement")
}

func (args *cniCmdArgs) GetInputArgs() *skel.CmdArgs {
	panic("not implement")
}

func (args *cniCmdArgs) GetNetNSPath() string {
	panic("not implement")
}

func (args *cniCmdArgs) Close() error {
	panic("not implement")
}

func isNSPathNotExist(err error) bool {
	panic("not implement")
}

func doCmdAdd(ctx context.Context, client rpc.TerwayBackendClient, cmdArgs *cniCmdArgs) (*terwayTypes.IPNetSet, *terwayTypes.IPSet, error) {
	panic("not implement")
}

func doCmdDel(ctx context.Context, client rpc.TerwayBackendClient, cmdArgs *cniCmdArgs) error {
	panic("not implement")
}

func doCmdCheck(ctx context.Context, client rpc.TerwayBackendClient, cmdArgs *cniCmdArgs) error {
	panic("not implement")
}
func prepareVF(ctx context.Context, id int, mac string) (int32, error) {
	return 0, fmt.Errorf("not support")
}
