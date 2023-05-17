// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: api/grpc/comagic.proto

package pb

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

// ComagicServiceClient is the client API for ComagicService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ComagicServiceClient interface {
	PushCallsToBQ(ctx context.Context, in *PushCallsToBQRequest, opts ...grpc.CallOption) (*PushCallsToBQResponse, error)
	PushOfflineMessagesToBQ(ctx context.Context, in *PushOfflineMessagesToBQRequest, opts ...grpc.CallOption) (*PushOfflineMessagesToBQResponse, error)
	GetCampaigns(ctx context.Context, in *GetCampaignsRequest, opts ...grpc.CallOption) (*GetCampaignsResponse, error)
	GetCampaignsConditions(ctx context.Context, in *GetCampaignsConditionsRequest, opts ...grpc.CallOption) (*GetCampaignsConditionsResponse, error)
}

type comagicServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewComagicServiceClient(cc grpc.ClientConnInterface) ComagicServiceClient {
	return &comagicServiceClient{cc}
}

func (c *comagicServiceClient) PushCallsToBQ(ctx context.Context, in *PushCallsToBQRequest, opts ...grpc.CallOption) (*PushCallsToBQResponse, error) {
	out := new(PushCallsToBQResponse)
	err := c.cc.Invoke(ctx, "/comagic.ComagicService/PushCallsToBQ", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *comagicServiceClient) PushOfflineMessagesToBQ(ctx context.Context, in *PushOfflineMessagesToBQRequest, opts ...grpc.CallOption) (*PushOfflineMessagesToBQResponse, error) {
	out := new(PushOfflineMessagesToBQResponse)
	err := c.cc.Invoke(ctx, "/comagic.ComagicService/PushOfflineMessagesToBQ", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *comagicServiceClient) GetCampaigns(ctx context.Context, in *GetCampaignsRequest, opts ...grpc.CallOption) (*GetCampaignsResponse, error) {
	out := new(GetCampaignsResponse)
	err := c.cc.Invoke(ctx, "/comagic.ComagicService/GetCampaigns", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *comagicServiceClient) GetCampaignsConditions(ctx context.Context, in *GetCampaignsConditionsRequest, opts ...grpc.CallOption) (*GetCampaignsConditionsResponse, error) {
	out := new(GetCampaignsConditionsResponse)
	err := c.cc.Invoke(ctx, "/comagic.ComagicService/GetCampaignsConditions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ComagicServiceServer is the server API for ComagicService service.
// All implementations must embed UnimplementedComagicServiceServer
// for forward compatibility
type ComagicServiceServer interface {
	PushCallsToBQ(context.Context, *PushCallsToBQRequest) (*PushCallsToBQResponse, error)
	PushOfflineMessagesToBQ(context.Context, *PushOfflineMessagesToBQRequest) (*PushOfflineMessagesToBQResponse, error)
	GetCampaigns(context.Context, *GetCampaignsRequest) (*GetCampaignsResponse, error)
	GetCampaignsConditions(context.Context, *GetCampaignsConditionsRequest) (*GetCampaignsConditionsResponse, error)
	mustEmbedUnimplementedComagicServiceServer()
}

// UnimplementedComagicServiceServer must be embedded to have forward compatible implementations.
type UnimplementedComagicServiceServer struct {
}

func (UnimplementedComagicServiceServer) PushCallsToBQ(context.Context, *PushCallsToBQRequest) (*PushCallsToBQResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushCallsToBQ not implemented")
}
func (UnimplementedComagicServiceServer) PushOfflineMessagesToBQ(context.Context, *PushOfflineMessagesToBQRequest) (*PushOfflineMessagesToBQResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushOfflineMessagesToBQ not implemented")
}
func (UnimplementedComagicServiceServer) GetCampaigns(context.Context, *GetCampaignsRequest) (*GetCampaignsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCampaigns not implemented")
}
func (UnimplementedComagicServiceServer) GetCampaignsConditions(context.Context, *GetCampaignsConditionsRequest) (*GetCampaignsConditionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCampaignsConditions not implemented")
}
func (UnimplementedComagicServiceServer) mustEmbedUnimplementedComagicServiceServer() {}

// UnsafeComagicServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ComagicServiceServer will
// result in compilation errors.
type UnsafeComagicServiceServer interface {
	mustEmbedUnimplementedComagicServiceServer()
}

func RegisterComagicServiceServer(s grpc.ServiceRegistrar, srv ComagicServiceServer) {
	s.RegisterService(&ComagicService_ServiceDesc, srv)
}

func _ComagicService_PushCallsToBQ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushCallsToBQRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComagicServiceServer).PushCallsToBQ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comagic.ComagicService/PushCallsToBQ",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComagicServiceServer).PushCallsToBQ(ctx, req.(*PushCallsToBQRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ComagicService_PushOfflineMessagesToBQ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushOfflineMessagesToBQRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComagicServiceServer).PushOfflineMessagesToBQ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comagic.ComagicService/PushOfflineMessagesToBQ",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComagicServiceServer).PushOfflineMessagesToBQ(ctx, req.(*PushOfflineMessagesToBQRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ComagicService_GetCampaigns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCampaignsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComagicServiceServer).GetCampaigns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comagic.ComagicService/GetCampaigns",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComagicServiceServer).GetCampaigns(ctx, req.(*GetCampaignsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ComagicService_GetCampaignsConditions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCampaignsConditionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComagicServiceServer).GetCampaignsConditions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comagic.ComagicService/GetCampaignsConditions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComagicServiceServer).GetCampaignsConditions(ctx, req.(*GetCampaignsConditionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ComagicService_ServiceDesc is the grpc.ServiceDesc for ComagicService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ComagicService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "comagic.ComagicService",
	HandlerType: (*ComagicServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushCallsToBQ",
			Handler:    _ComagicService_PushCallsToBQ_Handler,
		},
		{
			MethodName: "PushOfflineMessagesToBQ",
			Handler:    _ComagicService_PushOfflineMessagesToBQ_Handler,
		},
		{
			MethodName: "GetCampaigns",
			Handler:    _ComagicService_GetCampaigns_Handler,
		},
		{
			MethodName: "GetCampaignsConditions",
			Handler:    _ComagicService_GetCampaignsConditions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/grpc/comagic.proto",
}
