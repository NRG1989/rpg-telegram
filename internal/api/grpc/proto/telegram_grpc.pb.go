// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: assets/telegram.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GoAutRegistrationUserTelegramClient is the client API for GoAutRegistrationUserTelegram service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GoAutRegistrationUserTelegramClient interface {
	SendCode(ctx context.Context, in *SendCodeRequest, opts ...grpc.CallOption) (*SendCodeResponse, error)
}

type goAutRegistrationUserTelegramClient struct {
	cc grpc.ClientConnInterface
}

func NewGoAutRegistrationUserTelegramClient(cc grpc.ClientConnInterface) GoAutRegistrationUserTelegramClient {
	return &goAutRegistrationUserTelegramClient{cc}
}

func (c *goAutRegistrationUserTelegramClient) SendCode(ctx context.Context, in *SendCodeRequest, opts ...grpc.CallOption) (*SendCodeResponse, error) {
	out := new(SendCodeResponse)
	err := c.cc.Invoke(ctx, "/telegram_server.GoAutRegistrationUserTelegram/SendCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GoAutRegistrationUserTelegramServer is the server API for GoAutRegistrationUserTelegram service.
// All implementations must embed UnimplementedGoAutRegistrationUserTelegramServer
// for forward compatibility
type GoAutRegistrationUserTelegramServer interface {
	SendCode(context.Context, *SendCodeRequest) (*SendCodeResponse, error)
	mustEmbedUnimplementedGoAutRegistrationUserTelegramServer()
}

// UnimplementedGoAutRegistrationUserTelegramServer must be embedded to have forward compatible implementations.
type UnimplementedGoAutRegistrationUserTelegramServer struct {
}

func (UnimplementedGoAutRegistrationUserTelegramServer) SendCode(context.Context, *SendCodeRequest) (*SendCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCode not implemented")
}
func (UnimplementedGoAutRegistrationUserTelegramServer) mustEmbedUnimplementedGoAutRegistrationUserTelegramServer() {
}

// UnsafeGoAutRegistrationUserTelegramServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GoAutRegistrationUserTelegramServer will
// result in compilation errors.
type UnsafeGoAutRegistrationUserTelegramServer interface {
	mustEmbedUnimplementedGoAutRegistrationUserTelegramServer()
}

func RegisterGoAutRegistrationUserTelegramServer(s grpc.ServiceRegistrar, srv GoAutRegistrationUserTelegramServer) {
	s.RegisterService(&GoAutRegistrationUserTelegram_ServiceDesc, srv)
}

func _GoAutRegistrationUserTelegram_SendCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoAutRegistrationUserTelegramServer).SendCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/telegram_server.GoAutRegistrationUserTelegram/SendCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoAutRegistrationUserTelegramServer).SendCode(ctx, req.(*SendCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GoAutRegistrationUserTelegram_ServiceDesc is the grpc.ServiceDesc for GoAutRegistrationUserTelegram service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GoAutRegistrationUserTelegram_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "telegram_server.GoAutRegistrationUserTelegram",
	HandlerType: (*GoAutRegistrationUserTelegramServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendCode",
			Handler:    _GoAutRegistrationUserTelegram_SendCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "assets/telegram.proto",
}
