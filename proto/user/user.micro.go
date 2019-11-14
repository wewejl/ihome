// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/user/user.proto

package go_micro_srv_user

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for User service

type UserService interface {
	MicroUser(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	UpdataUser(ctx context.Context, in *UpdataUserRequest, opts ...client.CallOption) (*UpdataUserResponse, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.user"
	}
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) MicroUser(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "User.MicroUser", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) UpdataUser(ctx context.Context, in *UpdataUserRequest, opts ...client.CallOption) (*UpdataUserResponse, error) {
	req := c.c.NewRequest(c.name, "User.UpdataUser", in)
	out := new(UpdataUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for User service

type UserHandler interface {
	MicroUser(context.Context, *Request, *Response) error
	UpdataUser(context.Context, *UpdataUserRequest, *UpdataUserResponse) error
}

func RegisterUserHandler(s server.Server, hdlr UserHandler, opts ...server.HandlerOption) error {
	type user interface {
		MicroUser(ctx context.Context, in *Request, out *Response) error
		UpdataUser(ctx context.Context, in *UpdataUserRequest, out *UpdataUserResponse) error
	}
	type User struct {
		user
	}
	h := &userHandler{hdlr}
	return s.Handle(s.NewHandler(&User{h}, opts...))
}

type userHandler struct {
	UserHandler
}

func (h *userHandler) MicroUser(ctx context.Context, in *Request, out *Response) error {
	return h.UserHandler.MicroUser(ctx, in, out)
}

func (h *userHandler) UpdataUser(ctx context.Context, in *UpdataUserRequest, out *UpdataUserResponse) error {
	return h.UserHandler.UpdataUser(ctx, in, out)
}
