// Code generated by protoc-gen-bm v0.1, DO NOT EDIT.
// source: api.proto

/*
Package api is a generated blademaster stub package.
This code was generated with kratos/tool/protobuf/protoc-gen-bm v0.1.

package 命名使用 {appid}.{version} 的方式, version 形如 v1, v2 ..

It is generated from these files:
	api.proto
*/
package api

import (
	"context"

	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
)
import google_protobuf1 "google.golang.org/protobuf/types/known/emptypb"

// to suppressed 'imported but not used warning'
var _ *bm.Context
var _ context.Context
var _ binding.StructValidator

var PathCarRecordPing = "/demo.service.v1.CarRecord/Ping"
var PathCarRecordSayHello = "/demo.service.v1.CarRecord/SayHello"
var PathCarRecordSayHelloURL = "/kratos-demo/say_hello"
var PathCarRecordGetCarRecord = "/api/custom-lifang/getcarrecord"

// CarRecordBMServer is the server API for CarRecord service.
type CarRecordBMServer interface {
	Ping(ctx context.Context, req *google_protobuf1.Empty) (resp *google_protobuf1.Empty, err error)

	SayHello(ctx context.Context, req *HelloReq) (resp *google_protobuf1.Empty, err error)

	SayHelloURL(ctx context.Context, req *HelloReq) (resp *HelloResp, err error)

	//  rpc UploadCarRecord(UploadCarRecordReq) returns (UploadCarRecordResp) {
	//    option (google.api.http) = {
	//      post: "/api/custom-lifang/uploadrecord"
	//    };
	//  };
	GetCarRecord(ctx context.Context, req *GetCarRecordReq) (resp *CarRecordInfoList, err error)
}

var CarRecordSvc CarRecordBMServer

func carRecordPing(c *bm.Context) {
	p := new(google_protobuf1.Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := CarRecordSvc.Ping(c, p)
	c.JSON(resp, err)
}

func carRecordSayHello(c *bm.Context) {
	p := new(HelloReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := CarRecordSvc.SayHello(c, p)
	c.JSON(resp, err)
}

func carRecordSayHelloURL(c *bm.Context) {
	p := new(HelloReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := CarRecordSvc.SayHelloURL(c, p)
	c.JSON(resp, err)
}

func carRecordGetCarRecord(c *bm.Context) {
	p := new(GetCarRecordReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := CarRecordSvc.GetCarRecord(c, p)
	c.JSON(resp, err)
}

// RegisterCarRecordBMServer Register the blademaster route
func RegisterCarRecordBMServer(e *bm.Engine, server CarRecordBMServer) {
	CarRecordSvc = server
	e.GET("/demo.service.v1.CarRecord/Ping", carRecordPing)
	e.GET("/demo.service.v1.CarRecord/SayHello", carRecordSayHello)
	e.GET("/kratos-demo/say_hello", carRecordSayHelloURL)
	e.GET("/api/custom-lifang/getcarrecord", carRecordGetCarRecord)
}
