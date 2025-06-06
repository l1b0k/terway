// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	daemon "github.com/AliyunContainerService/terway/types/daemon"

	mock "github.com/stretchr/testify/mock"

	netip "net/netip"
)

// Factory is an autogenerated mock type for the Factory type
type Factory struct {
	mock.Mock
}

// AssignNIPv4 provides a mock function with given fields: eniID, count, mac
func (_m *Factory) AssignNIPv4(eniID string, count int, mac string) ([]netip.Addr, error) {
	ret := _m.Called(eniID, count, mac)

	if len(ret) == 0 {
		panic("no return value specified for AssignNIPv4")
	}

	var r0 []netip.Addr
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int, string) ([]netip.Addr, error)); ok {
		return rf(eniID, count, mac)
	}
	if rf, ok := ret.Get(0).(func(string, int, string) []netip.Addr); ok {
		r0 = rf(eniID, count, mac)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]netip.Addr)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int, string) error); ok {
		r1 = rf(eniID, count, mac)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AssignNIPv6 provides a mock function with given fields: eniID, count, mac
func (_m *Factory) AssignNIPv6(eniID string, count int, mac string) ([]netip.Addr, error) {
	ret := _m.Called(eniID, count, mac)

	if len(ret) == 0 {
		panic("no return value specified for AssignNIPv6")
	}

	var r0 []netip.Addr
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int, string) ([]netip.Addr, error)); ok {
		return rf(eniID, count, mac)
	}
	if rf, ok := ret.Get(0).(func(string, int, string) []netip.Addr); ok {
		r0 = rf(eniID, count, mac)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]netip.Addr)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int, string) error); ok {
		r1 = rf(eniID, count, mac)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateNetworkInterface provides a mock function with given fields: ipv4, ipv6, eniType
func (_m *Factory) CreateNetworkInterface(ipv4 int, ipv6 int, eniType string) (*daemon.ENI, []netip.Addr, []netip.Addr, error) {
	ret := _m.Called(ipv4, ipv6, eniType)

	if len(ret) == 0 {
		panic("no return value specified for CreateNetworkInterface")
	}

	var r0 *daemon.ENI
	var r1 []netip.Addr
	var r2 []netip.Addr
	var r3 error
	if rf, ok := ret.Get(0).(func(int, int, string) (*daemon.ENI, []netip.Addr, []netip.Addr, error)); ok {
		return rf(ipv4, ipv6, eniType)
	}
	if rf, ok := ret.Get(0).(func(int, int, string) *daemon.ENI); ok {
		r0 = rf(ipv4, ipv6, eniType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*daemon.ENI)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, string) []netip.Addr); ok {
		r1 = rf(ipv4, ipv6, eniType)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]netip.Addr)
		}
	}

	if rf, ok := ret.Get(2).(func(int, int, string) []netip.Addr); ok {
		r2 = rf(ipv4, ipv6, eniType)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).([]netip.Addr)
		}
	}

	if rf, ok := ret.Get(3).(func(int, int, string) error); ok {
		r3 = rf(ipv4, ipv6, eniType)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// DeleteNetworkInterface provides a mock function with given fields: eniID
func (_m *Factory) DeleteNetworkInterface(eniID string) error {
	ret := _m.Called(eniID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteNetworkInterface")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(eniID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAttachedNetworkInterface provides a mock function with given fields: preferTrunkID
func (_m *Factory) GetAttachedNetworkInterface(preferTrunkID string) ([]*daemon.ENI, error) {
	ret := _m.Called(preferTrunkID)

	if len(ret) == 0 {
		panic("no return value specified for GetAttachedNetworkInterface")
	}

	var r0 []*daemon.ENI
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*daemon.ENI, error)); ok {
		return rf(preferTrunkID)
	}
	if rf, ok := ret.Get(0).(func(string) []*daemon.ENI); ok {
		r0 = rf(preferTrunkID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*daemon.ENI)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(preferTrunkID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadNetworkInterface provides a mock function with given fields: mac
func (_m *Factory) LoadNetworkInterface(mac string) ([]netip.Addr, []netip.Addr, error) {
	ret := _m.Called(mac)

	if len(ret) == 0 {
		panic("no return value specified for LoadNetworkInterface")
	}

	var r0 []netip.Addr
	var r1 []netip.Addr
	var r2 error
	if rf, ok := ret.Get(0).(func(string) ([]netip.Addr, []netip.Addr, error)); ok {
		return rf(mac)
	}
	if rf, ok := ret.Get(0).(func(string) []netip.Addr); ok {
		r0 = rf(mac)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]netip.Addr)
		}
	}

	if rf, ok := ret.Get(1).(func(string) []netip.Addr); ok {
		r1 = rf(mac)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]netip.Addr)
		}
	}

	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(mac)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UnAssignNIPv4 provides a mock function with given fields: eniID, ips, mac
func (_m *Factory) UnAssignNIPv4(eniID string, ips []netip.Addr, mac string) error {
	ret := _m.Called(eniID, ips, mac)

	if len(ret) == 0 {
		panic("no return value specified for UnAssignNIPv4")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []netip.Addr, string) error); ok {
		r0 = rf(eniID, ips, mac)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnAssignNIPv6 provides a mock function with given fields: eniID, ips, mac
func (_m *Factory) UnAssignNIPv6(eniID string, ips []netip.Addr, mac string) error {
	ret := _m.Called(eniID, ips, mac)

	if len(ret) == 0 {
		panic("no return value specified for UnAssignNIPv6")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []netip.Addr, string) error); ok {
		r0 = rf(eniID, ips, mac)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFactory creates a new instance of Factory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *Factory {
	mock := &Factory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
