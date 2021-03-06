// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package prime_number

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

// PrimeNumberServiceClient is the client API for PrimeNumberService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrimeNumberServiceClient interface {
	IsPrimeNumber(ctx context.Context, in *IsPrimeNumberRequest, opts ...grpc.CallOption) (*IsPrimeNumberResponse, error)
}

type primeNumberServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPrimeNumberServiceClient(cc grpc.ClientConnInterface) PrimeNumberServiceClient {
	return &primeNumberServiceClient{cc}
}

func (c *primeNumberServiceClient) IsPrimeNumber(ctx context.Context, in *IsPrimeNumberRequest, opts ...grpc.CallOption) (*IsPrimeNumberResponse, error) {
	out := new(IsPrimeNumberResponse)
	err := c.cc.Invoke(ctx, "/prime_number.PrimeNumberService/IsPrimeNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrimeNumberServiceServer is the server API for PrimeNumberService service.
// All implementations must embed UnimplementedPrimeNumberServiceServer
// for forward compatibility
type PrimeNumberServiceServer interface {
	IsPrimeNumber(context.Context, *IsPrimeNumberRequest) (*IsPrimeNumberResponse, error)
	mustEmbedUnimplementedPrimeNumberServiceServer()
}

// UnimplementedPrimeNumberServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPrimeNumberServiceServer struct {
}

func (UnimplementedPrimeNumberServiceServer) IsPrimeNumber(context.Context, *IsPrimeNumberRequest) (*IsPrimeNumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsPrimeNumber not implemented")
}
func (UnimplementedPrimeNumberServiceServer) mustEmbedUnimplementedPrimeNumberServiceServer() {}

// UnsafePrimeNumberServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrimeNumberServiceServer will
// result in compilation errors.
type UnsafePrimeNumberServiceServer interface {
	mustEmbedUnimplementedPrimeNumberServiceServer()
}

func RegisterPrimeNumberServiceServer(s grpc.ServiceRegistrar, srv PrimeNumberServiceServer) {
	s.RegisterService(&PrimeNumberService_ServiceDesc, srv)
}

func _PrimeNumberService_IsPrimeNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsPrimeNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrimeNumberServiceServer).IsPrimeNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/prime_number.PrimeNumberService/IsPrimeNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrimeNumberServiceServer).IsPrimeNumber(ctx, req.(*IsPrimeNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PrimeNumberService_ServiceDesc is the grpc.ServiceDesc for PrimeNumberService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrimeNumberService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "prime_number.PrimeNumberService",
	HandlerType: (*PrimeNumberServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsPrimeNumber",
			Handler:    _PrimeNumberService_IsPrimeNumber_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "prime_number.proto",
}
