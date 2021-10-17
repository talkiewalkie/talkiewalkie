// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// UtilsClient is the client API for Utils service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UtilsClient interface {
	HealthCheck(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type utilsClient struct {
	cc grpc.ClientConnInterface
}

func NewUtilsClient(cc grpc.ClientConnInterface) UtilsClient {
	return &utilsClient{cc}
}

func (c *utilsClient) HealthCheck(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/app.Utils/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UtilsServer is the server API for Utils service.
// All implementations should embed UnimplementedUtilsServer
// for forward compatibility
type UtilsServer interface {
	HealthCheck(context.Context, *Empty) (*Empty, error)
}

// UnimplementedUtilsServer should be embedded to have forward compatible implementations.
type UnimplementedUtilsServer struct {
}

func (UnimplementedUtilsServer) HealthCheck(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}

// UnsafeUtilsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UtilsServer will
// result in compilation errors.
type UnsafeUtilsServer interface {
	mustEmbedUnimplementedUtilsServer()
}

func RegisterUtilsServer(s grpc.ServiceRegistrar, srv UtilsServer) {
	s.RegisterService(&Utils_ServiceDesc, srv)
}

func _Utils_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UtilsServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.Utils/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UtilsServer).HealthCheck(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Utils_ServiceDesc is the grpc.ServiceDesc for Utils service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Utils_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "app.Utils",
	HandlerType: (*UtilsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _Utils_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "app.proto",
}

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	SyncContacts(ctx context.Context, in *SyncContactsInput, opts ...grpc.CallOption) (*SyncContactsOutput, error)
	Onboarding(ctx context.Context, in *OnboardingInput, opts ...grpc.CallOption) (*MeUser, error)
	Me(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MeUser, error)
	Get(ctx context.Context, in *UserGetInput, opts ...grpc.CallOption) (*User, error)
	List(ctx context.Context, in *UserListInput, opts ...grpc.CallOption) (UserService_ListClient, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) SyncContacts(ctx context.Context, in *SyncContactsInput, opts ...grpc.CallOption) (*SyncContactsOutput, error) {
	out := new(SyncContactsOutput)
	err := c.cc.Invoke(ctx, "/app.UserService/SyncContacts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Onboarding(ctx context.Context, in *OnboardingInput, opts ...grpc.CallOption) (*MeUser, error) {
	out := new(MeUser)
	err := c.cc.Invoke(ctx, "/app.UserService/Onboarding", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Me(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MeUser, error) {
	out := new(MeUser)
	err := c.cc.Invoke(ctx, "/app.UserService/Me", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Get(ctx context.Context, in *UserGetInput, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/app.UserService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) List(ctx context.Context, in *UserListInput, opts ...grpc.CallOption) (UserService_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserService_ServiceDesc.Streams[0], "/app.UserService/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserService_ListClient interface {
	Recv() (*User, error)
	grpc.ClientStream
}

type userServiceListClient struct {
	grpc.ClientStream
}

func (x *userServiceListClient) Recv() (*User, error) {
	m := new(User)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations should embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	SyncContacts(context.Context, *SyncContactsInput) (*SyncContactsOutput, error)
	Onboarding(context.Context, *OnboardingInput) (*MeUser, error)
	Me(context.Context, *Empty) (*MeUser, error)
	Get(context.Context, *UserGetInput) (*User, error)
	List(*UserListInput, UserService_ListServer) error
}

// UnimplementedUserServiceServer should be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) SyncContacts(context.Context, *SyncContactsInput) (*SyncContactsOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncContacts not implemented")
}
func (UnimplementedUserServiceServer) Onboarding(context.Context, *OnboardingInput) (*MeUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Onboarding not implemented")
}
func (UnimplementedUserServiceServer) Me(context.Context, *Empty) (*MeUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Me not implemented")
}
func (UnimplementedUserServiceServer) Get(context.Context, *UserGetInput) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedUserServiceServer) List(*UserListInput, UserService_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_SyncContacts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncContactsInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).SyncContacts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.UserService/SyncContacts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).SyncContacts(ctx, req.(*SyncContactsInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Onboarding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnboardingInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Onboarding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.UserService/Onboarding",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Onboarding(ctx, req.(*OnboardingInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Me_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Me(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.UserService/Me",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Me(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserGetInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.UserService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Get(ctx, req.(*UserGetInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserListInput)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServiceServer).List(m, &userServiceListServer{stream})
}

type UserService_ListServer interface {
	Send(*User) error
	grpc.ServerStream
}

type userServiceListServer struct {
	grpc.ServerStream
}

func (x *userServiceListServer) Send(m *User) error {
	return x.ServerStream.SendMsg(m)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "app.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SyncContacts",
			Handler:    _UserService_SyncContacts_Handler,
		},
		{
			MethodName: "Onboarding",
			Handler:    _UserService_Onboarding_Handler,
		},
		{
			MethodName: "Me",
			Handler:    _UserService_Me_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _UserService_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _UserService_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "app.proto",
}

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	Send(ctx context.Context, in *MessageSendInput, opts ...grpc.CallOption) (*Empty, error)
	Incoming(ctx context.Context, in *Empty, opts ...grpc.CallOption) (MessageService_IncomingClient, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) Send(ctx context.Context, in *MessageSendInput, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/app.MessageService/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) Incoming(ctx context.Context, in *Empty, opts ...grpc.CallOption) (MessageService_IncomingClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageService_ServiceDesc.Streams[0], "/app.MessageService/Incoming", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageServiceIncomingClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessageService_IncomingClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type messageServiceIncomingClient struct {
	grpc.ClientStream
}

func (x *messageServiceIncomingClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations should embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	Send(context.Context, *MessageSendInput) (*Empty, error)
	Incoming(*Empty, MessageService_IncomingServer) error
}

// UnimplementedMessageServiceServer should be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) Send(context.Context, *MessageSendInput) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedMessageServiceServer) Incoming(*Empty, MessageService_IncomingServer) error {
	return status.Errorf(codes.Unimplemented, "method Incoming not implemented")
}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageSendInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.MessageService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Send(ctx, req.(*MessageSendInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_Incoming_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessageServiceServer).Incoming(m, &messageServiceIncomingServer{stream})
}

type MessageService_IncomingServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type messageServiceIncomingServer struct {
	grpc.ServerStream
}

func (x *messageServiceIncomingServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "app.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _MessageService_Send_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Incoming",
			Handler:       _MessageService_Incoming_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "app.proto",
}

// ConversationServiceClient is the client API for ConversationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConversationServiceClient interface {
	Get(ctx context.Context, in *ConversationGetInput, opts ...grpc.CallOption) (*Conversation, error)
	// TODO: Use ConversationService as output, delayed for demo
	List(ctx context.Context, in *ConversationListInput, opts ...grpc.CallOption) (ConversationService_ListClient, error)
}

type conversationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConversationServiceClient(cc grpc.ClientConnInterface) ConversationServiceClient {
	return &conversationServiceClient{cc}
}

func (c *conversationServiceClient) Get(ctx context.Context, in *ConversationGetInput, opts ...grpc.CallOption) (*Conversation, error) {
	out := new(Conversation)
	err := c.cc.Invoke(ctx, "/app.ConversationService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *conversationServiceClient) List(ctx context.Context, in *ConversationListInput, opts ...grpc.CallOption) (ConversationService_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConversationService_ServiceDesc.Streams[0], "/app.ConversationService/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &conversationServiceListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConversationService_ListClient interface {
	Recv() (*Conversation, error)
	grpc.ClientStream
}

type conversationServiceListClient struct {
	grpc.ClientStream
}

func (x *conversationServiceListClient) Recv() (*Conversation, error) {
	m := new(Conversation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConversationServiceServer is the server API for ConversationService service.
// All implementations should embed UnimplementedConversationServiceServer
// for forward compatibility
type ConversationServiceServer interface {
	Get(context.Context, *ConversationGetInput) (*Conversation, error)
	// TODO: Use ConversationService as output, delayed for demo
	List(*ConversationListInput, ConversationService_ListServer) error
}

// UnimplementedConversationServiceServer should be embedded to have forward compatible implementations.
type UnimplementedConversationServiceServer struct {
}

func (UnimplementedConversationServiceServer) Get(context.Context, *ConversationGetInput) (*Conversation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedConversationServiceServer) List(*ConversationListInput, ConversationService_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}

// UnsafeConversationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConversationServiceServer will
// result in compilation errors.
type UnsafeConversationServiceServer interface {
	mustEmbedUnimplementedConversationServiceServer()
}

func RegisterConversationServiceServer(s grpc.ServiceRegistrar, srv ConversationServiceServer) {
	s.RegisterService(&ConversationService_ServiceDesc, srv)
}

func _ConversationService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConversationGetInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConversationServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.ConversationService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConversationServiceServer).Get(ctx, req.(*ConversationGetInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConversationService_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConversationListInput)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConversationServiceServer).List(m, &conversationServiceListServer{stream})
}

type ConversationService_ListServer interface {
	Send(*Conversation) error
	grpc.ServerStream
}

type conversationServiceListServer struct {
	grpc.ServerStream
}

func (x *conversationServiceListServer) Send(m *Conversation) error {
	return x.ServerStream.SendMsg(m)
}

// ConversationService_ServiceDesc is the grpc.ServiceDesc for ConversationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConversationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "app.ConversationService",
	HandlerType: (*ConversationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _ConversationService_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _ConversationService_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "app.proto",
}
