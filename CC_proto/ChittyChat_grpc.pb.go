// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package CCBuffers

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

// ChittyChatClient is the client API for ChittyChat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittyChatClient interface {
	Join(ctx context.Context, in *ParticipantInfo, opts ...grpc.CallOption) (*ParticipantId, error)
	Recieve(ctx context.Context, in *ParticipantId, opts ...grpc.CallOption) (ChittyChat_RecieveClient, error)
	Leave(ctx context.Context, in *ParticipantId, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Publish(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type chittyChatClient struct {
	cc grpc.ClientConnInterface
}

func NewChittyChatClient(cc grpc.ClientConnInterface) ChittyChatClient {
	return &chittyChatClient{cc}
}

func (c *chittyChatClient) Join(ctx context.Context, in *ParticipantInfo, opts ...grpc.CallOption) (*ParticipantId, error) {
	out := new(ParticipantId)
	err := c.cc.Invoke(ctx, "/CCBuffers.ChittyChat/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chittyChatClient) Recieve(ctx context.Context, in *ParticipantId, opts ...grpc.CallOption) (ChittyChat_RecieveClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChittyChat_ServiceDesc.Streams[0], "/CCBuffers.ChittyChat/Recieve", opts...)
	if err != nil {
		return nil, err
	}
	x := &chittyChatRecieveClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChittyChat_RecieveClient interface {
	Recv() (*MsgRecieved, error)
	grpc.ClientStream
}

type chittyChatRecieveClient struct {
	grpc.ClientStream
}

func (x *chittyChatRecieveClient) Recv() (*MsgRecieved, error) {
	m := new(MsgRecieved)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chittyChatClient) Leave(ctx context.Context, in *ParticipantId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/CCBuffers.ChittyChat/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chittyChatClient) Publish(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/CCBuffers.ChittyChat/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChittyChatServer is the server API for ChittyChat service.
// All implementations must embed UnimplementedChittyChatServer
// for forward compatibility
type ChittyChatServer interface {
	Join(context.Context, *ParticipantInfo) (*ParticipantId, error)
	Recieve(*ParticipantId, ChittyChat_RecieveServer) error
	Leave(context.Context, *ParticipantId) (*emptypb.Empty, error)
	Publish(context.Context, *Msg) (*emptypb.Empty, error)
	mustEmbedUnimplementedChittyChatServer()
}

// UnimplementedChittyChatServer must be embedded to have forward compatible implementations.
type UnimplementedChittyChatServer struct {
}

func (UnimplementedChittyChatServer) Join(context.Context, *ParticipantInfo) (*ParticipantId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedChittyChatServer) Recieve(*ParticipantId, ChittyChat_RecieveServer) error {
	return status.Errorf(codes.Unimplemented, "method Recieve not implemented")
}
func (UnimplementedChittyChatServer) Leave(context.Context, *ParticipantId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedChittyChatServer) Publish(context.Context, *Msg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedChittyChatServer) mustEmbedUnimplementedChittyChatServer() {}

// UnsafeChittyChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittyChatServer will
// result in compilation errors.
type UnsafeChittyChatServer interface {
	mustEmbedUnimplementedChittyChatServer()
}

func RegisterChittyChatServer(s grpc.ServiceRegistrar, srv ChittyChatServer) {
	s.RegisterService(&ChittyChat_ServiceDesc, srv)
}

func _ChittyChat_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParticipantInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CCBuffers.ChittyChat/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServer).Join(ctx, req.(*ParticipantInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChittyChat_Recieve_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ParticipantId)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChittyChatServer).Recieve(m, &chittyChatRecieveServer{stream})
}

type ChittyChat_RecieveServer interface {
	Send(*MsgRecieved) error
	grpc.ServerStream
}

type chittyChatRecieveServer struct {
	grpc.ServerStream
}

func (x *chittyChatRecieveServer) Send(m *MsgRecieved) error {
	return x.ServerStream.SendMsg(m)
}

func _ChittyChat_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParticipantId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CCBuffers.ChittyChat/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServer).Leave(ctx, req.(*ParticipantId))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChittyChat_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/CCBuffers.ChittyChat/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServer).Publish(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

// ChittyChat_ServiceDesc is the grpc.ServiceDesc for ChittyChat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChittyChat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CCBuffers.ChittyChat",
	HandlerType: (*ChittyChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _ChittyChat_Join_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _ChittyChat_Leave_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _ChittyChat_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Recieve",
			Handler:       _ChittyChat_Recieve_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "CC_proto/ChittyChat.proto",
}
