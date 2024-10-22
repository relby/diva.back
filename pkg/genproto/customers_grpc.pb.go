// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: customers.proto

package genproto

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

// CustomersServiceClient is the client API for CustomersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CustomersServiceClient interface {
	ListCustomers(ctx context.Context, in *GetCustomersRequest, opts ...grpc.CallOption) (*GetCustomersResponse, error)
	UpdateCustomer(ctx context.Context, in *UpdateCustomerRequest, opts ...grpc.CallOption) (*UpdateCustomerResponse, error)
	AddCustomer(ctx context.Context, in *AddCustomerRequest, opts ...grpc.CallOption) (*AddCustomerResponse, error)
	DeleteCustomer(ctx context.Context, in *DeleteCustomerRequest, opts ...grpc.CallOption) (*DeleteCustomerResponse, error)
	ExportCustomersToExcel(ctx context.Context, in *ExportCustomersToExcelRequest, opts ...grpc.CallOption) (*ExportCustomersToExcelResponse, error)
}

type customersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCustomersServiceClient(cc grpc.ClientConnInterface) CustomersServiceClient {
	return &customersServiceClient{cc}
}

func (c *customersServiceClient) ListCustomers(ctx context.Context, in *GetCustomersRequest, opts ...grpc.CallOption) (*GetCustomersResponse, error) {
	out := new(GetCustomersResponse)
	err := c.cc.Invoke(ctx, "/customers.CustomersService/ListCustomers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) UpdateCustomer(ctx context.Context, in *UpdateCustomerRequest, opts ...grpc.CallOption) (*UpdateCustomerResponse, error) {
	out := new(UpdateCustomerResponse)
	err := c.cc.Invoke(ctx, "/customers.CustomersService/UpdateCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) AddCustomer(ctx context.Context, in *AddCustomerRequest, opts ...grpc.CallOption) (*AddCustomerResponse, error) {
	out := new(AddCustomerResponse)
	err := c.cc.Invoke(ctx, "/customers.CustomersService/AddCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) DeleteCustomer(ctx context.Context, in *DeleteCustomerRequest, opts ...grpc.CallOption) (*DeleteCustomerResponse, error) {
	out := new(DeleteCustomerResponse)
	err := c.cc.Invoke(ctx, "/customers.CustomersService/DeleteCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customersServiceClient) ExportCustomersToExcel(ctx context.Context, in *ExportCustomersToExcelRequest, opts ...grpc.CallOption) (*ExportCustomersToExcelResponse, error) {
	out := new(ExportCustomersToExcelResponse)
	err := c.cc.Invoke(ctx, "/customers.CustomersService/ExportCustomersToExcel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CustomersServiceServer is the server API for CustomersService service.
// All implementations must embed UnimplementedCustomersServiceServer
// for forward compatibility
type CustomersServiceServer interface {
	ListCustomers(context.Context, *GetCustomersRequest) (*GetCustomersResponse, error)
	UpdateCustomer(context.Context, *UpdateCustomerRequest) (*UpdateCustomerResponse, error)
	AddCustomer(context.Context, *AddCustomerRequest) (*AddCustomerResponse, error)
	DeleteCustomer(context.Context, *DeleteCustomerRequest) (*DeleteCustomerResponse, error)
	ExportCustomersToExcel(context.Context, *ExportCustomersToExcelRequest) (*ExportCustomersToExcelResponse, error)
	mustEmbedUnimplementedCustomersServiceServer()
}

// UnimplementedCustomersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCustomersServiceServer struct {
}

func (UnimplementedCustomersServiceServer) ListCustomers(context.Context, *GetCustomersRequest) (*GetCustomersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCustomers not implemented")
}
func (UnimplementedCustomersServiceServer) UpdateCustomer(context.Context, *UpdateCustomerRequest) (*UpdateCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCustomer not implemented")
}
func (UnimplementedCustomersServiceServer) AddCustomer(context.Context, *AddCustomerRequest) (*AddCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCustomer not implemented")
}
func (UnimplementedCustomersServiceServer) DeleteCustomer(context.Context, *DeleteCustomerRequest) (*DeleteCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCustomer not implemented")
}
func (UnimplementedCustomersServiceServer) ExportCustomersToExcel(context.Context, *ExportCustomersToExcelRequest) (*ExportCustomersToExcelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportCustomersToExcel not implemented")
}
func (UnimplementedCustomersServiceServer) mustEmbedUnimplementedCustomersServiceServer() {}

// UnsafeCustomersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CustomersServiceServer will
// result in compilation errors.
type UnsafeCustomersServiceServer interface {
	mustEmbedUnimplementedCustomersServiceServer()
}

func RegisterCustomersServiceServer(s grpc.ServiceRegistrar, srv CustomersServiceServer) {
	s.RegisterService(&CustomersService_ServiceDesc, srv)
}

func _CustomersService_ListCustomers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCustomersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).ListCustomers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/customers.CustomersService/ListCustomers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).ListCustomers(ctx, req.(*GetCustomersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_UpdateCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).UpdateCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/customers.CustomersService/UpdateCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).UpdateCustomer(ctx, req.(*UpdateCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_AddCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).AddCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/customers.CustomersService/AddCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).AddCustomer(ctx, req.(*AddCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_DeleteCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).DeleteCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/customers.CustomersService/DeleteCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).DeleteCustomer(ctx, req.(*DeleteCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomersService_ExportCustomersToExcel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportCustomersToExcelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomersServiceServer).ExportCustomersToExcel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/customers.CustomersService/ExportCustomersToExcel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomersServiceServer).ExportCustomersToExcel(ctx, req.(*ExportCustomersToExcelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CustomersService_ServiceDesc is the grpc.ServiceDesc for CustomersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CustomersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "customers.CustomersService",
	HandlerType: (*CustomersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListCustomers",
			Handler:    _CustomersService_ListCustomers_Handler,
		},
		{
			MethodName: "UpdateCustomer",
			Handler:    _CustomersService_UpdateCustomer_Handler,
		},
		{
			MethodName: "AddCustomer",
			Handler:    _CustomersService_AddCustomer_Handler,
		},
		{
			MethodName: "DeleteCustomer",
			Handler:    _CustomersService_DeleteCustomer_Handler,
		},
		{
			MethodName: "ExportCustomersToExcel",
			Handler:    _CustomersService_ExportCustomersToExcel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "customers.proto",
}
