package redis

import (
	"fmt"
	"github.com/aijie/michat/config"
	"github.com/aijie/michat/datas/model"
	"github.com/aijie/michat/datas/repository"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"sync"
	"time"
)

//public
//var RedisCli *redis.Conn

var (
	once      sync.Once
)

const (
	appKey = "imApp"
	appExpire = time.Duration(24 * time.Hour)
)

func NewRedis() redis.Conn {
	var con redis.Conn
	var err error
	readTimeout := redis.DialReadTimeout(appExpire)
	once.Do(func() {
		con, err = redis.Dial(
			"tcp",
			config.LogicConf.RedisIp,
			readTimeout,
		)
		if err != nil {
			log.Fatal("Failed to Dial redis ", err)
		}
	})
	return con
}

type imRedis struct {
	con redis.Conn
}

func NewIMRedisResp(con redis.Conn) repository.IMRepository {
	return &imRedis{con: con}
}

func (r *imRedis)fmtSubject(id int64, key string) string {
	subject := fmt.Sprintf("%s.%s", key, strconv.Itoa(int(id)))
	return subject
}

func (r *imRedis) SetAppInfo(app model.App) error {
	field := r.fmtSubject(app.Id, appKey)
	_, err := r.con.Do("HMSET", redis.Args{}.Add(field).AddFlat(app)...)
	return err
}

func (r *imRedis) GetAppInfo(appId int64) (*model.App, error) {
	field := r.fmtSubject(appId, appKey)
	reply, err := redis.Values(r.con.Do("HGET", field))
	if err != nil {
		return nil, err
	}
	app := &model.App{}
	redis.Scan(reply, app)
	return app, nil
}

//func InitDb() {
//	addr := config.LogicConf.RedisIp
//	RedisCli = redis.NewClient(&redis.Options{
//		Addr: addr,
//		DB: 0,
//	})
//	_, err := RedisCli.Ping().Result()
//	if err != nil {
//		logger.Sugar.Error(err)
//		panic(err)
//		return
//	}
//}
