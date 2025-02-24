package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/AliyunContainerService/terway/pkg/aliyun/client"
	"github.com/AliyunContainerService/terway/pkg/aliyun/credential"
	"github.com/AliyunContainerService/terway/pkg/aliyun/instance"
)

var (
	accessKeyID     string
	accessKeySecret string
	credentialPath  string
	region          string
	mode            string
	ipStack         string
)

func init() {
	flag.StringVar(&accessKeyID, "access-key-id", "", "AlibabaCloud Access Key ID")
	flag.StringVar(&accessKeySecret, "access-key-secret", "", "AlibabaCloud Access Key Secret")
	flag.StringVar(&credentialPath, "credential-path", "", "AlibabaCloud credential path")
	flag.StringVar(&region, "region", "", "AlibabaCloud Access Key Secret")
	flag.StringVar(&mode, "mode", "terway-eniip", "max pod cal mode: eni-ip|eni")
	flag.StringVar(&ipStack, "ip-stack", "ipv4", "ip stack")
}

func main() {
	flag.Parse()
	log.SetOutput(io.Discard)
	regionID, err := instance.GetInstanceMeta().GetRegionID()
	if err != nil {
		panic(err)
	}
	instanceType, err := instance.GetInstanceMeta().GetInstanceType()
	if err != nil {
		panic(err)
	}
	providers := []credential.Interface{
		credential.NewAKPairProvider(accessKeyID, accessKeySecret),
		credential.NewEncryptedCredentialProvider(credentialPath),
		credential.NewMetadataProvider(),
	}

	c, err := credential.NewClientMgr(regionID, providers...)
	if err != nil {
		panic(err)
	}

	api, err := client.New(c, nil)
	if err != nil {
		panic(err)
	}

	if mode == "terway-eniip" {
		limit, err := client.LimitProviders["ecs"].GetLimit(api, instanceType)
		if err != nil {
			panic(err)
		}
		fmt.Println(limit.IPv4PerAdapter * (limit.Adapters - 1))
	} else if mode == "terway-eni" {
		limit, err := client.LimitProviders["ecs"].GetLimit(api, instanceType)
		if err != nil {
			panic(err)
		}
		fmt.Println(limit.Adapters - 1)
	}
}
