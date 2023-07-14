// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: upload.proto

package upload

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

// Client API for UploadService service

type UploadService interface {
	// 获取上传入口地址
	UploadEntry(ctx context.Context, in *ReqEntry, opts ...client.CallOption) (*RespEntry, error)
}

type uploadService struct {
	c    client.Client
	name string
}

func NewUploadService(name string, c client.Client) UploadService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.service.upload"
	}
	return &uploadService{
		c:    c,
		name: name,
	}
}

func (c *uploadService) UploadEntry(ctx context.Context, in *ReqEntry, opts ...client.CallOption) (*RespEntry, error) {
	req := c.c.NewRequest(c.name, "UploadService.UploadEntry", in)
	out := new(RespEntry)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UploadService service

type UploadServiceHandler interface {
	// 获取上传入口地址
	UploadEntry(context.Context, *ReqEntry, *RespEntry) error
}

func RegisterUploadServiceHandler(s server.Server, hdlr UploadServiceHandler, opts ...server.HandlerOption) error {
	type uploadService interface {
		UploadEntry(ctx context.Context, in *ReqEntry, out *RespEntry) error
	}
	type UploadService struct {
		uploadService
	}
	h := &uploadServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UploadService{h}, opts...))
}

type uploadServiceHandler struct {
	UploadServiceHandler
}

func (h *uploadServiceHandler) UploadEntry(ctx context.Context, in *ReqEntry, out *RespEntry) error {
	return h.UploadServiceHandler.UploadEntry(ctx, in, out)
}