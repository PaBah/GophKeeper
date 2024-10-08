// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/gophkeeper/v1/service.proto

package proto

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
	GophKeeperService_SignUp_FullMethodName             = "/proto.gophkeeper.v1.GophKeeperService/SignUp"
	GophKeeperService_SignIn_FullMethodName             = "/proto.gophkeeper.v1.GophKeeperService/SignIn"
	GophKeeperService_CreateCredentials_FullMethodName  = "/proto.gophkeeper.v1.GophKeeperService/CreateCredentials"
	GophKeeperService_GetCredentials_FullMethodName     = "/proto.gophkeeper.v1.GophKeeperService/GetCredentials"
	GophKeeperService_UpdateCredentials_FullMethodName  = "/proto.gophkeeper.v1.GophKeeperService/UpdateCredentials"
	GophKeeperService_DeleteCredentials_FullMethodName  = "/proto.gophkeeper.v1.GophKeeperService/DeleteCredentials"
	GophKeeperService_CreateCard_FullMethodName         = "/proto.gophkeeper.v1.GophKeeperService/CreateCard"
	GophKeeperService_GetCards_FullMethodName           = "/proto.gophkeeper.v1.GophKeeperService/GetCards"
	GophKeeperService_UpdateCard_FullMethodName         = "/proto.gophkeeper.v1.GophKeeperService/UpdateCard"
	GophKeeperService_DeleteCard_FullMethodName         = "/proto.gophkeeper.v1.GophKeeperService/DeleteCard"
	GophKeeperService_GetFiles_FullMethodName           = "/proto.gophkeeper.v1.GophKeeperService/GetFiles"
	GophKeeperService_DeleteFile_FullMethodName         = "/proto.gophkeeper.v1.GophKeeperService/DeleteFile"
	GophKeeperService_SubscribeToChanges_FullMethodName = "/proto.gophkeeper.v1.GophKeeperService/SubscribeToChanges"
	GophKeeperService_UploadFile_FullMethodName         = "/proto.gophkeeper.v1.GophKeeperService/UploadFile"
	GophKeeperService_DownloadFile_FullMethodName       = "/proto.gophkeeper.v1.GophKeeperService/DownloadFile"
)

// GophKeeperServiceClient is the client API for GophKeeperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// GophKeeperService service provides ability to store date securely
type GophKeeperServiceClient interface {
	SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*SignUpResponse, error)
	SignIn(ctx context.Context, in *SignInRequest, opts ...grpc.CallOption) (*SignInResponse, error)
	CreateCredentials(ctx context.Context, in *CreateCredentialsRequest, opts ...grpc.CallOption) (*CreateCredentialsResponse, error)
	GetCredentials(ctx context.Context, in *GetCredentialsRequest, opts ...grpc.CallOption) (*GetCredentialsResponse, error)
	UpdateCredentials(ctx context.Context, in *UpdateCredentialsRequest, opts ...grpc.CallOption) (*UpdateCredentialsResponse, error)
	DeleteCredentials(ctx context.Context, in *DeleteCredentialsRequest, opts ...grpc.CallOption) (*DeleteCredentialsResponse, error)
	CreateCard(ctx context.Context, in *CreateCardRequest, opts ...grpc.CallOption) (*CreateCardResponse, error)
	GetCards(ctx context.Context, in *GetCardsRequest, opts ...grpc.CallOption) (*GetCardsResponse, error)
	UpdateCard(ctx context.Context, in *UpdateCardRequest, opts ...grpc.CallOption) (*UpdateCardResponse, error)
	DeleteCard(ctx context.Context, in *DeleteCardRequest, opts ...grpc.CallOption) (*DeleteCardResponse, error)
	GetFiles(ctx context.Context, in *GetFilesRequest, opts ...grpc.CallOption) (*GetFilesResponse, error)
	DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error)
	SubscribeToChanges(ctx context.Context, in *SubscribeToChangesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeToChangesResponse], error)
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[UploadFileRequest, UploadFileResponse], error)
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadFileResponse], error)
}

type gophKeeperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperServiceClient(cc grpc.ClientConnInterface) GophKeeperServiceClient {
	return &gophKeeperServiceClient{cc}
}

func (c *gophKeeperServiceClient) SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*SignUpResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SignUpResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_SignUp_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) SignIn(ctx context.Context, in *SignInRequest, opts ...grpc.CallOption) (*SignInResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SignInResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_SignIn_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) CreateCredentials(ctx context.Context, in *CreateCredentialsRequest, opts ...grpc.CallOption) (*CreateCredentialsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCredentialsResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_CreateCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) GetCredentials(ctx context.Context, in *GetCredentialsRequest, opts ...grpc.CallOption) (*GetCredentialsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCredentialsResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_GetCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) UpdateCredentials(ctx context.Context, in *UpdateCredentialsRequest, opts ...grpc.CallOption) (*UpdateCredentialsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateCredentialsResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_UpdateCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) DeleteCredentials(ctx context.Context, in *DeleteCredentialsRequest, opts ...grpc.CallOption) (*DeleteCredentialsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteCredentialsResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_DeleteCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) CreateCard(ctx context.Context, in *CreateCardRequest, opts ...grpc.CallOption) (*CreateCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCardResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_CreateCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) GetCards(ctx context.Context, in *GetCardsRequest, opts ...grpc.CallOption) (*GetCardsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCardsResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_GetCards_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) UpdateCard(ctx context.Context, in *UpdateCardRequest, opts ...grpc.CallOption) (*UpdateCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateCardResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_UpdateCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) DeleteCard(ctx context.Context, in *DeleteCardRequest, opts ...grpc.CallOption) (*DeleteCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteCardResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_DeleteCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) GetFiles(ctx context.Context, in *GetFilesRequest, opts ...grpc.CallOption) (*GetFilesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFilesResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_GetFiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteFileResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_DeleteFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) SubscribeToChanges(ctx context.Context, in *SubscribeToChangesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeToChangesResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &GophKeeperService_ServiceDesc.Streams[0], GophKeeperService_SubscribeToChanges_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeToChangesRequest, SubscribeToChangesResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type GophKeeperService_SubscribeToChangesClient = grpc.ServerStreamingClient[SubscribeToChangesResponse]

func (c *gophKeeperServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[UploadFileRequest, UploadFileResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &GophKeeperService_ServiceDesc.Streams[1], GophKeeperService_UploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadFileRequest, UploadFileResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type GophKeeperService_UploadFileClient = grpc.BidiStreamingClient[UploadFileRequest, UploadFileResponse]

func (c *gophKeeperServiceClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadFileResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &GophKeeperService_ServiceDesc.Streams[2], GophKeeperService_DownloadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadFileRequest, DownloadFileResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type GophKeeperService_DownloadFileClient = grpc.ServerStreamingClient[DownloadFileResponse]

// GophKeeperServiceServer is the server API for GophKeeperService service.
// All implementations must embed UnimplementedGophKeeperServiceServer
// for forward compatibility.
//
// GophKeeperService service provides ability to store date securely
type GophKeeperServiceServer interface {
	SignUp(context.Context, *SignUpRequest) (*SignUpResponse, error)
	SignIn(context.Context, *SignInRequest) (*SignInResponse, error)
	CreateCredentials(context.Context, *CreateCredentialsRequest) (*CreateCredentialsResponse, error)
	GetCredentials(context.Context, *GetCredentialsRequest) (*GetCredentialsResponse, error)
	UpdateCredentials(context.Context, *UpdateCredentialsRequest) (*UpdateCredentialsResponse, error)
	DeleteCredentials(context.Context, *DeleteCredentialsRequest) (*DeleteCredentialsResponse, error)
	CreateCard(context.Context, *CreateCardRequest) (*CreateCardResponse, error)
	GetCards(context.Context, *GetCardsRequest) (*GetCardsResponse, error)
	UpdateCard(context.Context, *UpdateCardRequest) (*UpdateCardResponse, error)
	DeleteCard(context.Context, *DeleteCardRequest) (*DeleteCardResponse, error)
	GetFiles(context.Context, *GetFilesRequest) (*GetFilesResponse, error)
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error)
	SubscribeToChanges(*SubscribeToChangesRequest, grpc.ServerStreamingServer[SubscribeToChangesResponse]) error
	UploadFile(grpc.BidiStreamingServer[UploadFileRequest, UploadFileResponse]) error
	DownloadFile(*DownloadFileRequest, grpc.ServerStreamingServer[DownloadFileResponse]) error
	mustEmbedUnimplementedGophKeeperServiceServer()
}

// UnimplementedGophKeeperServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGophKeeperServiceServer struct{}

func (UnimplementedGophKeeperServiceServer) SignUp(context.Context, *SignUpRequest) (*SignUpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedGophKeeperServiceServer) SignIn(context.Context, *SignInRequest) (*SignInResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignIn not implemented")
}
func (UnimplementedGophKeeperServiceServer) CreateCredentials(context.Context, *CreateCredentialsRequest) (*CreateCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCredentials not implemented")
}
func (UnimplementedGophKeeperServiceServer) GetCredentials(context.Context, *GetCredentialsRequest) (*GetCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCredentials not implemented")
}
func (UnimplementedGophKeeperServiceServer) UpdateCredentials(context.Context, *UpdateCredentialsRequest) (*UpdateCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCredentials not implemented")
}
func (UnimplementedGophKeeperServiceServer) DeleteCredentials(context.Context, *DeleteCredentialsRequest) (*DeleteCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCredentials not implemented")
}
func (UnimplementedGophKeeperServiceServer) CreateCard(context.Context, *CreateCardRequest) (*CreateCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCard not implemented")
}
func (UnimplementedGophKeeperServiceServer) GetCards(context.Context, *GetCardsRequest) (*GetCardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCards not implemented")
}
func (UnimplementedGophKeeperServiceServer) UpdateCard(context.Context, *UpdateCardRequest) (*UpdateCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCard not implemented")
}
func (UnimplementedGophKeeperServiceServer) DeleteCard(context.Context, *DeleteCardRequest) (*DeleteCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCard not implemented")
}
func (UnimplementedGophKeeperServiceServer) GetFiles(context.Context, *GetFilesRequest) (*GetFilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFiles not implemented")
}
func (UnimplementedGophKeeperServiceServer) DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedGophKeeperServiceServer) SubscribeToChanges(*SubscribeToChangesRequest, grpc.ServerStreamingServer[SubscribeToChangesResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToChanges not implemented")
}
func (UnimplementedGophKeeperServiceServer) UploadFile(grpc.BidiStreamingServer[UploadFileRequest, UploadFileResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedGophKeeperServiceServer) DownloadFile(*DownloadFileRequest, grpc.ServerStreamingServer[DownloadFileResponse]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedGophKeeperServiceServer) mustEmbedUnimplementedGophKeeperServiceServer() {}
func (UnimplementedGophKeeperServiceServer) testEmbeddedByValue()                           {}

// UnsafeGophKeeperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServiceServer will
// result in compilation errors.
type UnsafeGophKeeperServiceServer interface {
	mustEmbedUnimplementedGophKeeperServiceServer()
}

func RegisterGophKeeperServiceServer(s grpc.ServiceRegistrar, srv GophKeeperServiceServer) {
	// If the following call pancis, it indicates UnimplementedGophKeeperServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GophKeeperService_ServiceDesc, srv)
}

func _GophKeeperService_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_SignUp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).SignUp(ctx, req.(*SignUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_SignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).SignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_SignIn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).SignIn(ctx, req.(*SignInRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_CreateCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).CreateCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_CreateCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).CreateCredentials(ctx, req.(*CreateCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_GetCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).GetCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_GetCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).GetCredentials(ctx, req.(*GetCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_UpdateCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).UpdateCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_UpdateCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).UpdateCredentials(ctx, req.(*UpdateCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_DeleteCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).DeleteCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_DeleteCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).DeleteCredentials(ctx, req.(*DeleteCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_CreateCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).CreateCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_CreateCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).CreateCard(ctx, req.(*CreateCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_GetCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).GetCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_GetCards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).GetCards(ctx, req.(*GetCardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_UpdateCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).UpdateCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_UpdateCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).UpdateCard(ctx, req.(*UpdateCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_DeleteCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).DeleteCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_DeleteCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).DeleteCard(ctx, req.(*DeleteCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_GetFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).GetFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_GetFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).GetFiles(ctx, req.(*GetFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_DeleteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).DeleteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_DeleteFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).DeleteFile(ctx, req.(*DeleteFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_SubscribeToChanges_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeToChangesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GophKeeperServiceServer).SubscribeToChanges(m, &grpc.GenericServerStream[SubscribeToChangesRequest, SubscribeToChangesResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type GophKeeperService_SubscribeToChangesServer = grpc.ServerStreamingServer[SubscribeToChangesResponse]

func _GophKeeperService_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GophKeeperServiceServer).UploadFile(&grpc.GenericServerStream[UploadFileRequest, UploadFileResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type GophKeeperService_UploadFileServer = grpc.BidiStreamingServer[UploadFileRequest, UploadFileResponse]

func _GophKeeperService_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GophKeeperServiceServer).DownloadFile(m, &grpc.GenericServerStream[DownloadFileRequest, DownloadFileResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type GophKeeperService_DownloadFileServer = grpc.ServerStreamingServer[DownloadFileResponse]

// GophKeeperService_ServiceDesc is the grpc.ServiceDesc for GophKeeperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.gophkeeper.v1.GophKeeperService",
	HandlerType: (*GophKeeperServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignUp",
			Handler:    _GophKeeperService_SignUp_Handler,
		},
		{
			MethodName: "SignIn",
			Handler:    _GophKeeperService_SignIn_Handler,
		},
		{
			MethodName: "CreateCredentials",
			Handler:    _GophKeeperService_CreateCredentials_Handler,
		},
		{
			MethodName: "GetCredentials",
			Handler:    _GophKeeperService_GetCredentials_Handler,
		},
		{
			MethodName: "UpdateCredentials",
			Handler:    _GophKeeperService_UpdateCredentials_Handler,
		},
		{
			MethodName: "DeleteCredentials",
			Handler:    _GophKeeperService_DeleteCredentials_Handler,
		},
		{
			MethodName: "CreateCard",
			Handler:    _GophKeeperService_CreateCard_Handler,
		},
		{
			MethodName: "GetCards",
			Handler:    _GophKeeperService_GetCards_Handler,
		},
		{
			MethodName: "UpdateCard",
			Handler:    _GophKeeperService_UpdateCard_Handler,
		},
		{
			MethodName: "DeleteCard",
			Handler:    _GophKeeperService_DeleteCard_Handler,
		},
		{
			MethodName: "GetFiles",
			Handler:    _GophKeeperService_GetFiles_Handler,
		},
		{
			MethodName: "DeleteFile",
			Handler:    _GophKeeperService_DeleteFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToChanges",
			Handler:       _GophKeeperService_SubscribeToChanges_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadFile",
			Handler:       _GophKeeperService_UploadFile_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFile",
			Handler:       _GophKeeperService_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/gophkeeper/v1/service.proto",
}
