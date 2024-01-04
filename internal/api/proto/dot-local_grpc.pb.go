// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: proto/dot-local.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DotLocalClient is the client API for DotLocal service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DotLocalClient interface {
	CreateMapping(ctx context.Context, in *CreateMappingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type dotLocalClient struct {
	cc grpc.ClientConnInterface
}

func NewDotLocalClient(cc grpc.ClientConnInterface) DotLocalClient {
	return &dotLocalClient{cc}
}

func (c *dotLocalClient) CreateMapping(ctx context.Context, in *CreateMappingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/DotLocal/CreateMapping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DotLocalServer is the server API for DotLocal service.
// All implementations must embed UnimplementedDotLocalServer
// for forward compatibility
type DotLocalServer interface {
	CreateMapping(context.Context, *CreateMappingRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedDotLocalServer()
}

// UnimplementedDotLocalServer must be embedded to have forward compatible implementations.
type UnimplementedDotLocalServer struct {
}

func (UnimplementedDotLocalServer) CreateMapping(context.Context, *CreateMappingRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMapping not implemented")
}
func (UnimplementedDotLocalServer) mustEmbedUnimplementedDotLocalServer() {}

// UnsafeDotLocalServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DotLocalServer will
// result in compilation errors.
type UnsafeDotLocalServer interface {
	mustEmbedUnimplementedDotLocalServer()
}

func RegisterDotLocalServer(s grpc.ServiceRegistrar, srv DotLocalServer) {
	s.RegisterService(&DotLocal_ServiceDesc, srv)
}

func _DotLocal_CreateMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DotLocalServer).CreateMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DotLocal/CreateMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DotLocalServer).CreateMapping(ctx, req.(*CreateMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DotLocal_ServiceDesc is the grpc.ServiceDesc for DotLocal service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DotLocal_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "DotLocal",
	HandlerType: (*DotLocalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMapping",
			Handler:    _DotLocal_CreateMapping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/dot-local.proto",
}
