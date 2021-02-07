package dao

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/pkg/cache/memcache"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	xtime "github.com/go-kratos/kratos/pkg/time"
	"lifang/internal/model"

	"github.com/google/wire"
	carapi "lifang/api"
)

var Provider = wire.NewSet(New, NewDB, NewRedis, NewMC)

type (
	//go:generate kratos tool genbts
	// Dao dao interface
	Dao interface {
		Close()
		Ping(ctx context.Context) (err error)
		// bts: -nullcache=&model.Article{ID:-1} -check_null_code=$!=nil&&$.ID==-1
		Article(c context.Context, id int64) (*model.Article, error)
		// 存储车辆进出记录
		SaveCarRecord(c context.Context, req *carapi.UploadCarRecordReq, carId string) (err error)
		// 查询车辆进出记录
		GetCarRecord(c context.Context, req *carapi.GetCarRecordReq) (reply *carapi.CarRecordInfoList, err error)
		// 只通过时间来查询车辆进出记录
		GetRecordOnlyByTime(c context.Context, req *carapi.GetCarRecordReq) (reply *carapi.CarRecordInfoList, err error)
		//通过进出时间和车牌号来查询记录
		GetRecordbyCarCode(c context.Context, req *carapi.GetCarRecordReq) (reply *carapi.CarRecordInfoList, err error)
		//参数有进出时间和停车场ID
		GetRecordbyparkID(c context.Context, req *carapi.GetCarRecordReq) (reply *carapi.CarRecordInfoList, err error)
	}
)

// dao dao.
type dao struct {
	db         *sql.DB
	redis      *redis.Redis
	mc         *memcache.Memcache
	cache      *fanout.Fanout
	demoExpire int32
}

// New new a dao and return.
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d Dao, cf func(), err error) {
	return newDao(r, mc, db)
}

func newDao(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d *dao, cf func(), err error) {
	var cfg struct {
		DemoExpire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &dao{
		db:         db,
		redis:      r,
		mc:         mc,
		cache:      fanout.New("cache"),
		demoExpire: int32(time.Duration(cfg.DemoExpire) / time.Second),
	}
	cf = d.Close
	return
}

// Close close the resource.
func (d *dao) Close() {
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return nil
}
