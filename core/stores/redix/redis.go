package redix

import (
	"context"
	"errors"
	"time"

	"github.com/go-zelus/zelus/core/logx"

	"github.com/go-redis/redis/v8"
	"github.com/go-zelus/zelus/core/config"
	"github.com/go-zelus/zelus/core/util"
)

type Config struct {
	Name     string `mapstructure:"name"`
	Conn     string `mapstructure:"conn"`
	Password string `mapstructure:"password"`
	Timeout  int    `mapstructure:"timeout"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

var instance = make(map[string]*Redis)

func Init() {
	register()
}

// 注册redis
func register() {
	slice := make([]*Config, 0)
	err := config.UnmarshalKey("redis", &slice)
	if err != nil {
		logx.Fatalf("解析配置失败 key:[%s], %v", "redis", err)
	}
	if len(slice) == 0 {
		logx.Fatalf("未发现redis配置信息")
	}
	for _, conf := range slice {
		c := new(Redis)
		c.ctx = context.Background()
		if conf.Timeout == 0 {
			conf.Timeout = 60
		}
		if conf.PoolSize == 0 {
			conf.PoolSize = 15
		}
		c.client = redis.NewClient(&redis.Options{
			Addr:     conf.Conn,
			Password: conf.Password,
			DB:       conf.DB,
			//连接池容量及闲置连接数量
			PoolSize:     conf.PoolSize,   // 连接池数量
			MinIdleConns: 10,              //好比最小连接数
			DialTimeout:  5 * time.Second, //连接建立超时时间
			ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
			WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
			PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

			//闲置连接检查包括IdleTimeout，MaxConnAge
			IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
			MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

			//命令执行失败时的重试策略
			MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
			MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
			MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
			IdleTimeout:     time.Second * time.Duration(conf.Timeout),
		})
		err = c.client.Ping(c.ctx).Err()
		if err != nil {
			logx.Fatalf("redis 连接失败: %v", err)
		}
		instance[conf.Name] = c
	}
	logx.Info("redis 注册成功...")
}

// DB 根据名称获取数据库连接
func DB(name string) *Redis {
	v, ok := instance[name]
	if !ok {
		panic("get " + name + " redis db failed")
	}
	return v
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Context() context.Context {
	return r.ctx
}

// Set 将字符串值 value 关联到 key
func (r *Redis) Set(key string, value interface{}, ops ...*Operation) *Result {
	exp := Operations(ops).Find(EXPIRE).Result(time.Second * 0).(time.Duration)
	nx := Operations(ops).Find(NX).Result(nil)

	if nx != nil {
		return NewResult(r.client.SetNX(r.ctx, key, value, exp).Result())
	}
	xx := Operations(ops).Find(XX).Result(nil)
	if xx != nil {
		return NewResult(r.client.SetXX(r.ctx, key, value, exp).Result())
	}
	return NewResult(r.client.Set(r.ctx, key, value, exp).Result())
}

// Get 返回与键 key 相关联的字符串值
func (r *Redis) Get(key string) *Result {
	return NewResult(r.client.Get(r.ctx, key).Result())
}

// TTL 返回剩余有效期
func (r *Redis) TTL(key string) *Result {
	return NewResult(r.client.TTL(r.ctx, key).Result())
}

// Keys 返回key列表
func (r *Redis) Keys(arg string) *Result {
	return NewResult(r.client.Keys(r.ctx, arg).Result())
}

// MGet 返回包含了所有给定键的值的列表
func (r *Redis) MGet(key string) *MResult {
	return NewMResult(r.client.MGet(r.ctx, key).Result())
}

// Del 删除key
func (r *Redis) Del(key string) *Result {
	return NewResult(r.client.Del(r.ctx, key).Result())
}

// Inrc 数字加一
func (r *Redis) Inrc(key string) *Result {
	return NewResult(r.client.Incr(r.ctx, key).Result())
}

// Lock 分布式锁
func (r *Redis) Lock(key string, acquire, timeout time.Duration) (string, error) {
	code := util.NewUUID()
	endTime := time.Now().Add(acquire).UnixNano()
	for time.Now().UnixNano() <= endTime {
		if success, err := r.client.SetNX(r.ctx, key, code, timeout).Result(); err != nil {
			return "", err
		} else if success {
			return code, nil
		} else if r.client.TTL(r.ctx, key).Val() == -1 {
			r.client.Expire(r.ctx, key, timeout)
		}
		time.Sleep(time.Millisecond)
	}
	return "", errors.New("lock timeout")
}

// UnLock 释放分布式锁
func (r *Redis) UnLock(key, code string) bool {
	txf := func(tx *redis.Tx) error {
		if v, err := tx.Get(r.ctx, key).Result(); err != nil && err != redis.Nil {
			return err
		} else if v == code {
			_, err = tx.Pipelined(r.ctx, func(pipe redis.Pipeliner) error {
				pipe.Del(r.ctx, key)
				return nil
			})
			return err
		}
		return nil
	}

	for {
		if err := r.client.Watch(r.ctx, txf, key); err == nil {
			return true
		} else if err == redis.TxFailedErr {
			logx.Infof("watch key is modified,retry to release lock. err: %s\n", err.Error())
			logx.Infof("key: %s,code: %s\n", key, code)
		} else {
			return false
		}
	}
}
