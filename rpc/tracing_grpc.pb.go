// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: tracing.proto

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TerwayTracing_GetResourceTypes_FullMethodName   = "/rpc.TerwayTracing/GetResourceTypes"
	TerwayTracing_GetResources_FullMethodName       = "/rpc.TerwayTracing/GetResources"
	TerwayTracing_GetResourceConfig_FullMethodName  = "/rpc.TerwayTracing/GetResourceConfig"
	TerwayTracing_GetResourceTrace_FullMethodName   = "/rpc.TerwayTracing/GetResourceTrace"
	TerwayTracing_ResourceExecute_FullMethodName    = "/rpc.TerwayTracing/ResourceExecute"
	TerwayTracing_GetResourceMapping_FullMethodName = "/rpc.TerwayTracing/GetResourceMapping"
)

// TerwayTracingClient is the client API for TerwayTracing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TerwayTracingClient interface {
	GetResourceTypes(ctx context.Context, in *Placeholder, opts ...grpc.CallOption) (*ResourcesTypesReply, error)
	GetResources(ctx context.Context, in *ResourceTypeRequest, opts ...grpc.CallOption) (*ResourcesNamesReply, error)
	GetResourceConfig(ctx context.Context, in *ResourceTypeNameRequest, opts ...grpc.CallOption) (*ResourceConfigReply, error)
	GetResourceTrace(ctx context.Context, in *ResourceTypeNameRequest, opts ...grpc.CallOption) (*ResourceTraceReply, error)
	ResourceExecute(ctx context.Context, in *ResourceExecuteRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ResourceExecuteReply], error)
	GetResourceMapping(ctx context.Context, in *Placeholder, opts ...grpc.CallOption) (*ResourceMappingReply, error)
}

type terwayTracingClient struct {
	cc grpc.ClientConnInterface
}

func NewTerwayTracingClient(cc grpc.ClientConnInterface) TerwayTracingClient {
	return &terwayTracingClient{cc}
}

func (c *terwayTracingClient) GetResourceTypes(ctx context.Context, in *Placeholder, opts ...grpc.CallOption) (*ResourcesTypesReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResourcesTypesReply)
	err := c.cc.Invoke(ctx, TerwayTracing_GetResourceTypes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terwayTracingClient) GetResources(ctx context.Context, in *ResourceTypeRequest, opts ...grpc.CallOption) (*ResourcesNamesReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResourcesNamesReply)
	err := c.cc.Invoke(ctx, TerwayTracing_GetResources_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terwayTracingClient) GetResourceConfig(ctx context.Context, in *ResourceTypeNameRequest, opts ...grpc.CallOption) (*ResourceConfigReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResourceConfigReply)
	err := c.cc.Invoke(ctx, TerwayTracing_GetResourceConfig_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terwayTracingClient) GetResourceTrace(ctx context.Context, in *ResourceTypeNameRequest, opts ...grpc.CallOption) (*ResourceTraceReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResourceTraceReply)
	err := c.cc.Invoke(ctx, TerwayTracing_GetResourceTrace_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terwayTracingClient) ResourceExecute(ctx context.Context, in *ResourceExecuteRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ResourceExecuteReply], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &TerwayTracing_ServiceDesc.Streams[0], TerwayTracing_ResourceExecute_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ResourceExecuteRequest, ResourceExecuteReply]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TerwayTracing_ResourceExecuteClient = grpc.ServerStreamingClient[ResourceExecuteReply]

func (c *terwayTracingClient) GetResourceMapping(ctx context.Context, in *Placeholder, opts ...grpc.CallOption) (*ResourceMappingReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResourceMappingReply)
	err := c.cc.Invoke(ctx, TerwayTracing_GetResourceMapping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TerwayTracingServer is the server API for TerwayTracing service.
// All implementations must embed UnimplementedTerwayTracingServer
// for forward compatibility.
type TerwayTracingServer interface {
	GetResourceTypes(context.Context, *Placeholder) (*ResourcesTypesReply, error)
	GetResources(context.Context, *ResourceTypeRequest) (*ResourcesNamesReply, error)
	GetResourceConfig(context.Context, *ResourceTypeNameRequest) (*ResourceConfigReply, error)
	GetResourceTrace(context.Context, *ResourceTypeNameRequest) (*ResourceTraceReply, error)
	ResourceExecute(*ResourceExecuteRequest, grpc.ServerStreamingServer[ResourceExecuteReply]) error
	GetResourceMapping(context.Context, *Placeholder) (*ResourceMappingReply, error)
	mustEmbedUnimplementedTerwayTracingServer()
}

// UnimplementedTerwayTracingServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTerwayTracingServer struct{}

func (UnimplementedTerwayTracingServer) GetResourceTypes(context.Context, *Placeholder) (*ResourcesTypesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResourceTypes not implemented")
}
func (UnimplementedTerwayTracingServer) GetResources(context.Context, *ResourceTypeRequest) (*ResourcesNamesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResources not implemented")
}
func (UnimplementedTerwayTracingServer) GetResourceConfig(context.Context, *ResourceTypeNameRequest) (*ResourceConfigReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResourceConfig not implemented")
}
func (UnimplementedTerwayTracingServer) GetResourceTrace(context.Context, *ResourceTypeNameRequest) (*ResourceTraceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResourceTrace not implemented")
}
func (UnimplementedTerwayTracingServer) ResourceExecute(*ResourceExecuteRequest, grpc.ServerStreamingServer[ResourceExecuteReply]) error {
	return status.Errorf(codes.Unimplemented, "method ResourceExecute not implemented")
}
func (UnimplementedTerwayTracingServer) GetResourceMapping(context.Context, *Placeholder) (*ResourceMappingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResourceMapping not implemented")
}
func (UnimplementedTerwayTracingServer) mustEmbedUnimplementedTerwayTracingServer() {}
func (UnimplementedTerwayTracingServer) testEmbeddedByValue()                       {}

// UnsafeTerwayTracingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TerwayTracingServer will
// result in compilation errors.
type UnsafeTerwayTracingServer interface {
	mustEmbedUnimplementedTerwayTracingServer()
}

func RegisterTerwayTracingServer(s grpc.ServiceRegistrar, srv TerwayTracingServer) {
	// If the following call pancis, it indicates UnimplementedTerwayTracingServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TerwayTracing_ServiceDesc, srv)
}

func _TerwayTracing_GetResourceTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Placeholder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerwayTracingServer).GetResourceTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TerwayTracing_GetResourceTypes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerwayTracingServer).GetResourceTypes(ctx, req.(*Placeholder))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerwayTracing_GetResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerwayTracingServer).GetResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TerwayTracing_GetResources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerwayTracingServer).GetResources(ctx, req.(*ResourceTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerwayTracing_GetResourceConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceTypeNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerwayTracingServer).GetResourceConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TerwayTracing_GetResourceConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerwayTracingServer).GetResourceConfig(ctx, req.(*ResourceTypeNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerwayTracing_GetResourceTrace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceTypeNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerwayTracingServer).GetResourceTrace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TerwayTracing_GetResourceTrace_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerwayTracingServer).GetResourceTrace(ctx, req.(*ResourceTypeNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerwayTracing_ResourceExecute_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ResourceExecuteRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TerwayTracingServer).ResourceExecute(m, &grpc.GenericServerStream[ResourceExecuteRequest, ResourceExecuteReply]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type TerwayTracing_ResourceExecuteServer = grpc.ServerStreamingServer[ResourceExecuteReply]

func _TerwayTracing_GetResourceMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Placeholder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerwayTracingServer).GetResourceMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TerwayTracing_GetResourceMapping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerwayTracingServer).GetResourceMapping(ctx, req.(*Placeholder))
	}
	return interceptor(ctx, in, info, handler)
}

// TerwayTracing_ServiceDesc is the grpc.ServiceDesc for TerwayTracing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TerwayTracing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.TerwayTracing",
	HandlerType: (*TerwayTracingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetResourceTypes",
			Handler:    _TerwayTracing_GetResourceTypes_Handler,
		},
		{
			MethodName: "GetResources",
			Handler:    _TerwayTracing_GetResources_Handler,
		},
		{
			MethodName: "GetResourceConfig",
			Handler:    _TerwayTracing_GetResourceConfig_Handler,
		},
		{
			MethodName: "GetResourceTrace",
			Handler:    _TerwayTracing_GetResourceTrace_Handler,
		},
		{
			MethodName: "GetResourceMapping",
			Handler:    _TerwayTracing_GetResourceMapping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ResourceExecute",
			Handler:       _TerwayTracing_ResourceExecute_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "tracing.proto",
}
