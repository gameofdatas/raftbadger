// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.18.1
// source: badger.proto

package proto

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BadgerService_SaveEntries_FullMethodName   = "/badgerpb.BadgerService/SaveEntries"
	BadgerService_GetEntries_FullMethodName    = "/badgerpb.BadgerService/GetEntries"
	BadgerService_DeleteEntries_FullMethodName = "/badgerpb.BadgerService/DeleteEntries"
	BadgerService_UpdateEntries_FullMethodName = "/badgerpb.BadgerService/UpdateEntries"
)

// BadgerServiceClient is the client API for BadgerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// BadgerService defines RPC methods for distributed data management.
type BadgerServiceClient interface {
	// SaveEntries saves a secret in the system.
	SaveEntries(ctx context.Context, in *SaveEntriesRequest, opts ...grpc.CallOption) (*SaveEntriesResponse, error)
	// GetEntries retrieves a secret by ID.
	GetEntries(ctx context.Context, in *GetEntriesRequest, opts ...grpc.CallOption) (*GetEntriesResponse, error)
	// DeleteEntries removes a secret by ID.
	DeleteEntries(ctx context.Context, in *DeleteEntriesRequest, opts ...grpc.CallOption) (*DeleteEntriesResponse, error)
	// UpdateEntries updates a secret by ID.
	UpdateEntries(ctx context.Context, in *UpdateEntriesRequest, opts ...grpc.CallOption) (*UpdateEntriesResponse, error)
}

type badgerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBadgerServiceClient(cc grpc.ClientConnInterface) BadgerServiceClient {
	return &badgerServiceClient{cc}
}

func (c *badgerServiceClient) SaveEntries(ctx context.Context, in *SaveEntriesRequest, opts ...grpc.CallOption) (*SaveEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SaveEntriesResponse)
	err := c.cc.Invoke(ctx, BadgerService_SaveEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *badgerServiceClient) GetEntries(ctx context.Context, in *GetEntriesRequest, opts ...grpc.CallOption) (*GetEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEntriesResponse)
	err := c.cc.Invoke(ctx, BadgerService_GetEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *badgerServiceClient) DeleteEntries(ctx context.Context, in *DeleteEntriesRequest, opts ...grpc.CallOption) (*DeleteEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteEntriesResponse)
	err := c.cc.Invoke(ctx, BadgerService_DeleteEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *badgerServiceClient) UpdateEntries(ctx context.Context, in *UpdateEntriesRequest, opts ...grpc.CallOption) (*UpdateEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateEntriesResponse)
	err := c.cc.Invoke(ctx, BadgerService_UpdateEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BadgerServiceServer is the server API for BadgerService service.
// All implementations must embed UnimplementedBadgerServiceServer
// for forward compatibility.
//
// BadgerService defines RPC methods for distributed data management.
type BadgerServiceServer interface {
	// SaveEntries saves a secret in the system.
	SaveEntries(context.Context, *SaveEntriesRequest) (*SaveEntriesResponse, error)
	// GetEntries retrieves a secret by ID.
	GetEntries(context.Context, *GetEntriesRequest) (*GetEntriesResponse, error)
	// DeleteEntries removes a secret by ID.
	DeleteEntries(context.Context, *DeleteEntriesRequest) (*DeleteEntriesResponse, error)
	// UpdateEntries updates a secret by ID.
	UpdateEntries(context.Context, *UpdateEntriesRequest) (*UpdateEntriesResponse, error)
	mustEmbedUnimplementedBadgerServiceServer()
}

// UnimplementedBadgerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBadgerServiceServer struct{}

func (UnimplementedBadgerServiceServer) SaveEntries(context.Context, *SaveEntriesRequest) (*SaveEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveEntries not implemented")
}
func (UnimplementedBadgerServiceServer) GetEntries(context.Context, *GetEntriesRequest) (*GetEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEntries not implemented")
}
func (UnimplementedBadgerServiceServer) DeleteEntries(context.Context, *DeleteEntriesRequest) (*DeleteEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEntries not implemented")
}
func (UnimplementedBadgerServiceServer) UpdateEntries(context.Context, *UpdateEntriesRequest) (*UpdateEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEntries not implemented")
}
func (UnimplementedBadgerServiceServer) mustEmbedUnimplementedBadgerServiceServer() {}
func (UnimplementedBadgerServiceServer) testEmbeddedByValue()                       {}

// UnsafeBadgerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BadgerServiceServer will
// result in compilation errors.
type UnsafeBadgerServiceServer interface {
	mustEmbedUnimplementedBadgerServiceServer()
}

func RegisterBadgerServiceServer(s grpc.ServiceRegistrar, srv BadgerServiceServer) {
	// If the following call pancis, it indicates UnimplementedBadgerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BadgerService_ServiceDesc, srv)
}

func _BadgerService_SaveEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BadgerServiceServer).SaveEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BadgerService_SaveEntries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BadgerServiceServer).SaveEntries(ctx, req.(*SaveEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BadgerService_GetEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BadgerServiceServer).GetEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BadgerService_GetEntries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BadgerServiceServer).GetEntries(ctx, req.(*GetEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BadgerService_DeleteEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BadgerServiceServer).DeleteEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BadgerService_DeleteEntries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BadgerServiceServer).DeleteEntries(ctx, req.(*DeleteEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BadgerService_UpdateEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BadgerServiceServer).UpdateEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BadgerService_UpdateEntries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BadgerServiceServer).UpdateEntries(ctx, req.(*UpdateEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BadgerService_ServiceDesc is the grpc.ServiceDesc for BadgerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BadgerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "badgerpb.BadgerService",
	HandlerType: (*BadgerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveEntries",
			Handler:    _BadgerService_SaveEntries_Handler,
		},
		{
			MethodName: "GetEntries",
			Handler:    _BadgerService_GetEntries_Handler,
		},
		{
			MethodName: "DeleteEntries",
			Handler:    _BadgerService_DeleteEntries_Handler,
		},
		{
			MethodName: "UpdateEntries",
			Handler:    _BadgerService_UpdateEntries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "badger.proto",
}
