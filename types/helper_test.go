package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AliyunContainerService/terway/rpc"
)

func TestBuildIPNet_EmptyInputs_ReturnsEmptyIPNetSet(t *testing.T) {
	ipNetSet, err := BuildIPNet(nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.Nil(t, ipNetSet.IPv4)
	assert.Nil(t, ipNetSet.IPv6)
}

func TestBuildIPNet_PartiallyEmptyInputs_ReturnsEmptyIPNetSet(t *testing.T) {
	ip := &rpc.IPSet{IPv4: "192.168.1.1", IPv6: "2001:db8::1"}
	ipNetSet, err := BuildIPNet(ip, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.Nil(t, ipNetSet.IPv4)
	assert.Nil(t, ipNetSet.IPv6)

	subnet := &rpc.IPSet{IPv4: "192.168.1.0/24", IPv6: "2001:db8::/64"}
	ipNetSet, err = BuildIPNet(nil, subnet)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.Nil(t, ipNetSet.IPv4)
	assert.Nil(t, ipNetSet.IPv6)
}

func TestBuildIPNet_ValidInputs_ReturnsCorrectIPNetSet(t *testing.T) {
	ip := &rpc.IPSet{IPv4: "192.168.1.1", IPv6: "2001:db8::1"}
	subnet := &rpc.IPSet{IPv4: "192.168.1.0/24", IPv6: "2001:db8::/64"}
	ipNetSet, err := BuildIPNet(ip, subnet)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.NotNil(t, ipNetSet.IPv4)
	assert.NotNil(t, ipNetSet.IPv6)
	assert.Equal(t, "192.168.1.1/24", ipNetSet.IPv4.String())
	assert.Equal(t, "2001:db8::1/64", ipNetSet.IPv6.String())
}

func TestBuildIPNet_InvalidIP_ReturnsError(t *testing.T) {
	ip := &rpc.IPSet{IPv4: "invalid", IPv6: "2001:db8::1"}
	subnet := &rpc.IPSet{IPv4: "192.168.1.0/24", IPv6: "2001:db8::/64"}
	ipNetSet, err := BuildIPNet(ip, subnet)
	assert.Error(t, err)
	assert.Nil(t, ipNetSet)
}

func TestBuildIPNet_InvalidSubnet_ReturnsError(t *testing.T) {
	ip := &rpc.IPSet{IPv4: "192.168.1.1", IPv6: "2001:db8::1"}
	subnet := &rpc.IPSet{IPv4: "invalid", IPv6: "2001:db8::/64"}
	ipNetSet, err := BuildIPNet(ip, subnet)
	assert.Error(t, err)
	assert.Nil(t, ipNetSet)
}

func TestBuildIPNet_OnlyIPv4_ReturnsCorrectIPNetSet(t *testing.T) {
	ip := &rpc.IPSet{IPv4: "192.168.1.1"}
	subnet := &rpc.IPSet{IPv4: "192.168.1.0/24"}
	ipNetSet, err := BuildIPNet(ip, subnet)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.NotNil(t, ipNetSet.IPv4)
	assert.Nil(t, ipNetSet.IPv6)
	assert.Equal(t, "192.168.1.1/24", ipNetSet.IPv4.String())
}

func TestBuildIPNet_OnlyIPv6_ReturnsCorrectIPNetSet(t *testing.T) {
	ip := &rpc.IPSet{IPv6: "2001:db8::1"}
	subnet := &rpc.IPSet{IPv6: "2001:db8::/64"}
	ipNetSet, err := BuildIPNet(ip, subnet)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.Nil(t, ipNetSet.IPv4)
	assert.NotNil(t, ipNetSet.IPv6)
	assert.Equal(t, "2001:db8::1/64", ipNetSet.IPv6.String())
}

func TestToIPNetSetReturnsErrorWhenIPIsNil(t *testing.T) {
	ipNetSet, err := ToIPNetSet(nil)
	assert.Error(t, err)
	assert.Nil(t, ipNetSet)
}

func TestToIPNetSetReturnsCorrectIPNetSetWhenValidIPv4(t *testing.T) {
	ip := &rpc.IPSet{IPv4: "192.168.1.0/24"}
	ipNetSet, err := ToIPNetSet(ip)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.NotNil(t, ipNetSet.IPv4)
	assert.Nil(t, ipNetSet.IPv6)
	assert.Equal(t, "192.168.1.0/24", ipNetSet.IPv4.String())
}

func TestToIPNetSetReturnsCorrectIPNetSetWhenValidIPv6(t *testing.T) {
	ip := &rpc.IPSet{IPv6: "2001:db8::/64"}
	ipNetSet, err := ToIPNetSet(ip)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.Nil(t, ipNetSet.IPv4)
	assert.NotNil(t, ipNetSet.IPv6)
	assert.Equal(t, "2001:db8::/64", ipNetSet.IPv6.String())
}

func TestToIPNetSetReturnsCorrectIPNetSetWhenValidIPv4AndIPv6(t *testing.T) {
	ip := &rpc.IPSet{IPv4: "192.168.1.0/24", IPv6: "2001:db8::/64"}
	ipNetSet, err := ToIPNetSet(ip)
	assert.NoError(t, err)
	assert.NotNil(t, ipNetSet)
	assert.NotNil(t, ipNetSet.IPv4)
	assert.NotNil(t, ipNetSet.IPv6)
	assert.Equal(t, "192.168.1.0/24", ipNetSet.IPv4.String())
	assert.Equal(t, "2001:db8::/64", ipNetSet.IPv6.String())
}

func TestToIPNetSetReturnsErrorWhenInvalidIPv4(t *testing.T) {
	ip := &rpc.IPSet{IPv4: "invalid"}
	ipNetSet, err := ToIPNetSet(ip)
	assert.Error(t, err)
	assert.Nil(t, ipNetSet)
}

func TestToIPNetSetReturnsErrorWhenInvalidIPv6(t *testing.T) {
	ip := &rpc.IPSet{IPv6: "invalid"}
	ipNetSet, err := ToIPNetSet(ip)
	assert.Error(t, err)
	assert.Nil(t, ipNetSet)
}
