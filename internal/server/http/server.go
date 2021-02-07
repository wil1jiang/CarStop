package http

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/render"
	"lifang/internal/service"
	"net/http"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	pb "lifang/api"
	"lifang/internal/model"
)

var svc *service.Service

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine, err error) {
	var (
		cfg bm.ServerConfig
		ct  paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	svc = s
	engine = bm.DefaultServer(&cfg)
	pb.RegisterCarRecordBMServer(engine, s)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/lifang")
	{
		g.GET("/start", howToStart)
		e.POST("/api/custom-lifang/uploadrecord", UploadCarRecord)
	}
}


func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.
func howToStart(c *bm.Context) {
	k := &model.Kratos{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(k, nil)
}

func UploadCarRecord(c *bm.Context) {
	p := new(pb.UploadCarRecordReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	var bytes []byte
	recordResp, _ := svc.UploadCarRecord(context.Background(), p)
	bytes, _ = json.Marshal(recordResp)
	c.Render(http.StatusOK, render.Data{
		ContentType: "application/json",
		Data:        [][]byte{bytes},
	})
}

/*func carRecordGetCarRecord(c *bm.Context) {
	p := new(pb.GetCarRecordReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	var bytes []byte
	recordResp, _ := svc.GetCarRecord(context.Background(),p)
	bytes, _ = json.Marshal(recordResp)
	c.Render(http.StatusOK,render.Data{
		ContentType: "application/json",
		Data: [][]byte{bytes},
	})
	/*resp, err := CarRecordSvc.GetCarRecord(c, p)
	c.JSON(resp, err)
}*/