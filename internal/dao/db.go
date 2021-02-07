package dao

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/prometheus/common/log"
	carapi "lifang/api"
	"lifang/internal/model"
	"time"
)

var Db *sql.DB
var CarStop string = "治安总队停车场"

type GetCarRecordResp struct {
	carCode string
	inTime string
	passTime string
	parkID string
	inOrOut string
	GUID string
	channelID string
	channelName string
	imagePath string
	caruser string
}

func CarRecordTablrAdd(ctx context.Context) (err error){

	_, err = Db.Exec(ctx,"create table if not exists CarStopRecord(carCode varchar(1024),inTime datetime,passTime datetime,parkID varchar(1024),inOrOut varchar(1024),GUID varchar(1024),channelID varchar(1024),channelName varchar(1024),imagePath varchar(1024),caruser varchar(1024)) DEFAULT CHARACTER SET = utf8;")
	if err != nil{
		log.Info("CarRecordTablrAdd failed")
		return err
	}
	return
}

func NewDB() (db *sql.DB, cf func(), err error) {
	var (
		cfg sql.Config
		ct  paladin.TOML
	)
	if err = paladin.Get("db.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = sql.NewMySQL(&cfg)
	Db = db
	err = CarRecordTablrAdd(context.Background())
	if err !=nil{
		log.Info("创建表失败...")
	}
	cf = func() { db.Close() }
	return
}

func (d *dao) RawArticle(ctx context.Context, id int64) (art *model.Article, err error) {
	// get data from db
	return
}

func (d *dao) SaveCarRecord(c context.Context, req *carapi.UploadCarRecordReq, caruser string) (err error) {
	marshal, err := json.Marshal(req)
	if err != nil {
		return  err
	}

	err = CarRecordInsert(c,req,caruser)
	if err != nil {
		return  err
	}
	log.Info("marshal:",string(marshal))
	return
}

func (d *dao) GetCarRecord(ctx context.Context, req *carapi.GetCarRecordReq) (reply *carapi.CarRecordInfoList, err error) {
	var resp GetCarRecordResp
	reply = &carapi.CarRecordInfoList{}
	var CalNum int64
	start := (req.Page-1)*(req.Limit)
	end := (req.Page-1)*(req.Limit) + req.Limit
	inTime := time.Unix(0,0)
	passTime := time.Unix(0,0)

	rows,err := Db.Query(ctx,"SELECT * FROM CarStopRecord WHERE passTime >= ? AND passTime <= ? AND carCode = ? AND parkID = ? ",req.InTime,req.PassTime,req.CarCode,req.ParkID)
	if rows == nil || err != nil{
		return nil, err
	}
	for rows.Next(){
		if CalNum>=start && CalNum<end {
			err = rows.Scan(&resp.carCode,&inTime,&passTime,&resp.parkID,&resp.inOrOut,&resp.GUID,&resp.channelID,&resp.channelName,&resp.imagePath,&resp.caruser)
			if err != nil {
				log.Info("GetRecordOnlyByTime scan failed")
				return nil, err
			}
			info := carapi.CarRecordInfo{}
			info.CarCode = resp.carCode
			info.InTime = inTime.Format("2006-01-02 15:04:05")
			info.PassTime = passTime.Format("2006-01-02 15:04:05")
			info.ParkID = resp.parkID
			info.InOrOut = resp.inOrOut
			info.GUID = resp.GUID
			info.ChannelID = resp.channelID
			info.ChannelName = resp.channelName
			info.ImagePath = resp.imagePath
			info.Caruser = resp.caruser
			info.Stopaddr = CarStop
			reply.List = append(reply.List, &info)
			log.Info(resp.carCode,resp.inTime,resp.passTime,resp.parkID,resp.inOrOut,resp.GUID,resp.channelID,resp.channelName,resp.imagePath)
		}
		CalNum++
	}
	reply.Total = CalNum
	return reply, err
}

func (d *dao) GetRecordOnlyByTime(ctx context.Context, req *carapi.GetCarRecordReq) (reply *carapi.CarRecordInfoList, err error){
	/*start_time := req.InTime
	end_time := req.PassTime*/
	resp :=GetCarRecordResp{}
	reply =  &carapi.CarRecordInfoList{}
	var CalNum int64
	start := (req.Page-1)*(req.Limit)
	end := (req.Page-1)*(req.Limit) + req.Limit
	inTime := time.Unix(0,0)
	passTime := time.Unix(0,0)

	rows,err := Db.Query(ctx,"SELECT * FROM CarStopRecord WHERE passTime >= ? AND passTime <= ? ", req.InTime, req.PassTime)
	if rows == nil || err != nil{
		return nil, err
	}
	for rows.Next(){
		if CalNum>=start && CalNum<end {
			err = rows.Scan(&resp.carCode, &inTime, &passTime, &resp.parkID, &resp.inOrOut, &resp.GUID, &resp.channelID, &resp.channelName, &resp.imagePath, &resp.caruser)
			if err != nil {
				log.Info("GetRecordOnlyByTime scan failed")
				return nil, err
			}
			info := carapi.CarRecordInfo{}
			info.CarCode = resp.carCode
			info.InTime = inTime.Format("2006-01-02 15:04:05")
			info.PassTime = passTime.Format("2006-01-02 15:04:05")
			info.ParkID = resp.parkID
			info.InOrOut = resp.inOrOut
			info.GUID = resp.GUID
			info.ChannelID = resp.channelID
			info.ChannelName = resp.channelName
			info.ImagePath = resp.imagePath
			info.Caruser = resp.caruser
			info.Stopaddr = CarStop
			reply.List = append(reply.List, &info)

			log.Info(resp.carCode, resp.inTime, resp.passTime, resp.parkID, resp.inOrOut, resp.GUID, resp.channelID, resp.channelName, resp.imagePath)
		}
		CalNum++
	}
	reply.Total = CalNum
	return reply, err
}

func (d *dao)GetRecordbyCarCode(ctx context.Context, req *carapi.GetCarRecordReq) (reply *carapi.CarRecordInfoList, err error){
	var resp GetCarRecordResp
	reply = &carapi.CarRecordInfoList{}
	var CalNum int64
	start := (req.Page-1)*(req.Limit)
	end := (req.Page-1)*(req.Limit) + req.Limit
	inTime := time.Unix(0,0)
	passTime := time.Unix(0,0)

	rows,err := Db.Query(ctx,"SELECT * FROM CarStopRecord WHERE passTime >= ? AND passTime <= ? AND carCode LIKE ?",req.InTime,req.PassTime,req.CarCode + "%")
	if rows == nil || err != nil{
		return nil, err
	}
	for rows.Next() {
		if CalNum >= start && CalNum < end {
			err = rows.Scan(&resp.carCode, &inTime, &passTime, &resp.parkID, &resp.inOrOut, &resp.GUID, &resp.channelID, &resp.channelName, &resp.imagePath, &resp.caruser)
			if err != nil {
				log.Info("GetRecordOnlyByTime scan failed")
				return nil, err
			}
			info := carapi.CarRecordInfo{}
			info.CarCode = resp.carCode
			info.InTime = inTime.Format("2006-01-02 15:04:05")
			info.PassTime = passTime.Format("2006-01-02 15:04:05")
			info.ParkID = resp.parkID
			info.InOrOut = resp.inOrOut
			info.GUID = resp.GUID
			info.ChannelID = resp.channelID
			info.ChannelName = resp.channelName
			info.ImagePath = resp.imagePath
			info.Caruser = resp.caruser
			info.Stopaddr = CarStop
			reply.List = append(reply.List, &info)
			log.Info(resp.carCode, resp.inTime, resp.passTime, resp.parkID, resp.inOrOut, resp.GUID, resp.channelID, resp.channelName, resp.imagePath)
		}
		CalNum++
	}
	reply.Total = CalNum
	return reply,err
}

func (d *dao)GetRecordbyparkID(ctx context.Context, req *carapi.GetCarRecordReq) (reply *carapi.CarRecordInfoList, err error){
	var resp GetCarRecordResp
	reply = &carapi.CarRecordInfoList{}
	var CalNum int64
	inTime := time.Unix(0,0)
	passTime := time.Unix(0,0)
	start := (req.Page-1)*(req.Limit)
	end := (req.Page-1)*(req.Limit) + req.Limit

	rows,err := Db.Query(ctx,"SELECT * FROM CarStopRecord WHERE passTime >= ? AND passTime <= ? AND parkID = ? ",req.InTime,req.PassTime,req.ParkID)
	if rows == nil || err != nil{
		return nil, err
	}
	for rows.Next() {
		if CalNum >= start && CalNum < end {
			//err = rows.Scan(&resp.carCode, &resp.inTime, &resp.passTime, &resp.parkID, &resp.inOrOut, &resp.GUID, &resp.channelID, &resp.channelName, &resp.imagePath, &resp.caruser)
			err = rows.Scan(&resp.carCode, &inTime, &passTime, &resp.parkID, &resp.inOrOut, &resp.GUID, &resp.channelID, &resp.channelName, &resp.imagePath, &resp.caruser)

			if err != nil {
				log.Info("GetRecordOnlyByTime scan failed")
				return nil, err
			}
			info := carapi.CarRecordInfo{}
			info.CarCode = resp.carCode
			info.InTime = inTime.Format("2006-01-02 15:04:05")
			info.PassTime = passTime.Format("2006-01-02 15:04:05")
			info.ParkID = resp.parkID
			info.InOrOut = resp.inOrOut
			info.GUID = resp.GUID
			info.ChannelID = resp.channelID
			info.ChannelName = resp.channelName
			info.ImagePath = resp.imagePath
			info.Caruser = resp.caruser
			info.Stopaddr = CarStop
			reply.List = append(reply.List, &info)
			log.Info(resp.carCode, resp.inTime, resp.passTime, resp.parkID, resp.inOrOut, resp.GUID, resp.channelID, resp.channelName, resp.imagePath)
		}
		CalNum++
	}
	reply.Total = CalNum
	return reply,err
}

func CarRecordInsert(ctx context.Context, req *carapi.UploadCarRecordReq,caruser string) (err error) {
	stmt, err:= Db.Prepare("INSERT INTO CarStopRecord(carCode,inTime,passTime,parkID,inOrOut,GUID,channelID,channelName,imagePath,caruser) values(?,STR_TO_DATE(?,'%Y-%m-%d %H:%i:%s'),STR_TO_DATE(?,'%Y-%m-%d %H:%i:%s'),?,?,?,?,?,+ ?,?)")
	if err != nil{
		return err
	}
	//_, err = stmt.Exec(ctx,req.CarCode,req.InTime,req.PassTime,req.ParkID,req.InOrOut,req.GUID,req.ChannelID,req.ChannelName,"http://192.168.2.205:9988"+req.ImagePath,caruser)
	_, err = stmt.Exec(ctx,req.CarCode,req.InTime,req.PassTime,req.ParkID,req.InOrOut,req.GUID,req.ChannelID,req.ChannelName,req.ImagePath,caruser)
	if err != nil{
		return err
	}
	return nil
}

