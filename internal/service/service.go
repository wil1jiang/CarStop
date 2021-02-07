package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
	pb "lifang/api"
	"lifang/internal/dao"
	"lifang/internal/model"
	"lifang/internal/mq"
	"net/http"
	"time"
)

var(
	g_projectUuid string
	g_carip string
	g_carport int
	g_rabbitmqip string
	g_rabbitmqport int
	g_mqexchange string
	g_mqtopic string
	g_mqDns string
	G_stopaddr string
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.CarRecordBMServer), new(*Service)))

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

func configInit()(err error){
	var(
		ct paladin.TOML
		comcfgs model.Config
		carcfgs model.CarInterface
		mqconfig model.Rabbitmq
	)
	if err = paladin.Get("cfgs.toml").Unmarshal(&ct); err !=nil{
		return err
	}
	if err = ct.Get("common").UnmarshalTOML(&comcfgs); err!=nil{
		log.Info("解析common字段配置文件失败")
		return err
	}
	g_projectUuid = comcfgs.ProjectID
	G_stopaddr = comcfgs.Stopaddr

	if err = ct.Get("carinterface").UnmarshalTOML(&carcfgs); err!=nil{
		log.Info("解析carinterface字段配置文件失败")
		return err
	}
	g_carip  = carcfgs.Ip
	g_carport = carcfgs.Port

	if err = ct.Get("rabbitmq").UnmarshalTOML(&mqconfig); err!=nil{
		log.Info("解析rabbitmq字段配置文件失败")
		return err
	}
	g_mqexchange = mqconfig.Exchange
	g_mqtopic = mqconfig.Topic
	g_mqDns = fmt.Sprintf("amqp://%s:%s@%s:%d/%s",mqconfig.User,mqconfig.Password,mqconfig.Ip,mqconfig.Port,mqconfig.Vhost)
	return nil
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	configInit()
	mq.SetupMQ(g_mqDns)
	return
}

// SayHello grpc demo func.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	fmt.Printf("hello %s", req.Name)
	return
}

// SayHelloURL bm demo func.
func (s *Service) SayHelloURL(ctx context.Context, req *pb.HelloReq) (reply *pb.HelloResp, err error) {
	reply = &pb.HelloResp{
		Content: "hello " + req.Name,
	}
	fmt.Printf("hello url %s", req.Name)
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}

func (s *Service) UploadCarRecord(ctx context.Context, req *pb.UploadCarRecordReq) (reply *pb.UploadCarRecordResp, err error) {
	//首先通过车牌得到车主的名字，车主的名字也需要放入数据库中
	user, err := s.getCarUser(req.CarCode)
	if err != nil {
		return &pb.UploadCarRecordResp{
			ResCode: 0,
			ResMsg:  "获取车主信息失败",
		}, err
		log.Errorv(context.Background(), log.KV("log", "getCarUser"), log.KV("error", err))
	}
	//将请求的数据保存到数据库中
	err = s.dao.SaveCarRecord(ctx, req, user)
	if err != nil {
		reply = &pb.UploadCarRecordResp{
			ResCode: 0,
			ResMsg:  "数据库保存失败",
		}
		log.Errorv(context.Background(), log.KV("log", "数据库保存失败"), log.KV("error", err))
		return reply, err
	}

	mqdata := model.Data{
		Caruser: user,
		Carcode: req.CarCode,
		Stopaddr: G_stopaddr,
		PassType :req.InOrOut,
		PassTime :req.PassTime,
		ImagePath: req.ImagePath,
	}

	mq.Publish(g_mqexchange,g_mqtopic,model.MQCarInfo{ProjectUuid: g_projectUuid,
		Topic: g_mqtopic,
		Data: mqdata,
	})

	reply = &pb.UploadCarRecordResp{
		ResCode: 0,
		ResMsg:  "成功",
	}
	return reply, nil
}

func (s *Service) GetCarRecord(ctx context.Context, req *pb.GetCarRecordReq) (reply *pb.CarRecordInfoList, err error) {
	//前端没有选择进出时间直接返回错误
	if len(req.InTime) == 0 || len(req.PassTime) == 0 {
		err = ecode.Error(-1, "参数错误")
		return
	}

	//参数中只有进出时间
	if len(req.CarCode) == 0 && len(req.ParkID) == 0 {
		reply, err = s.dao.GetRecordOnlyByTime(ctx, req)
		if err != nil {
			err = ecode.Error(-1, "只通过进出时间查询记录失败")
		}
		return reply, err
	}

	//参数有进出时间和车牌号
	if len(req.ParkID) == 0 && len(req.CarCode) != 0 {
		reply, err = s.dao.GetRecordbyCarCode(ctx, req)
		if err != nil {
			err = ecode.Error(-1, "有进出时间和车牌号查询记录失败")
		}
		return reply, err
	}

	//参数有进出时间和停车场ID
	if len(req.CarCode) == 0 && len(req.ParkID) != 0 {
		reply, err = s.dao.GetRecordbyparkID(ctx, req)
		if err != nil {
			err = ecode.Error(-2, "通过进出时间和车场ID查询失败")
		}
		return reply, err
	}

	//四个参数都有
	reply, err = s.dao.GetCarRecord(ctx, req)
	if err != nil {
		err = ecode.Error(-1, "四个参数都有查询记录失败")
	}
	return reply, err
}

func (s *Service) getCarUser(carcode string) (user string, err error) {
	httpClient := resty.New()
	host := fmt.Sprintf("http://%s:%d", g_carip, g_carport)
	httpClient.SetHostURL(host)
	httpClient.SetTimeout(3 * time.Second)
	httpReq := httpClient.R()
	httpReq.Method = http.MethodGet
	httpReq.URL = fmt.Sprintf("/api/statistics-basedata-server/project/%s/parking/vehicle/list",g_projectUuid)
	httpReq.SetQueryParams(map[string]string{
		"licenseNumber": carcode,
		"page":          "1",
		"limit":         "13",
	})

	response, err := httpReq.Send()
	if err != nil {
		log.Info("getCarUser failed!")
		return "", nil
	}
	result := struct {
		Data struct {
			Total int `json:"total"`
			List  []struct {
				Picture        string `json:"picture"`
				VehicleUuid    string `json:"vehicleUuid"`
				LicenseNumber  string `json:"licenseNumber"`
				LicenseType    string `json:"licenseType"`
				LicenseColor   string `json:"licenseColor"`
				VehicleType    string `json:"vehicleType"`
				VehicleColor   string `json:"vehicleColor"`
				VehiclePicture string `json:"vehiclePicture"`
				CreateTime     string `json:"createTime"`
				UpdateTime     string `json:"updateTime"`
				StaffUuid      string `json:"staffUuid"`
				StaffName      string `json:"staffName"`
			} `json:"list"`
		} `json:"data"`

		Msg     string `json:"msg"`
		Success bool   `json:"success"`
		ErrCode int    `json:"errCode"`
	}{}
	err = json.Unmarshal(response.Body(), &result)
	if err != nil || result.Data.Total == 0 {
		return "", err
	}
	return result.Data.List[0].StaffName, err
}
