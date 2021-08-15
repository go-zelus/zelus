package mongox

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/go-zelus/zelus/core/config"
	"github.com/go-zelus/zelus/core/logx"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ErrorDataExisted = "E11000"                        // 数据已存在
	ErrorDataIsEmpty = "mongo: no documents in result" // 没有查询到数据
)

type Config struct {
	Database   string `mapstructure:"db"`
	Conn       string `mapstructure:"conn"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	AuthSource string `mapstructure:"auth_source"`
	Timeout    int    `mapstructure:"timeout"`
	MaxPool    int    `mapstructure:"max_pool"`
}

var instance *mongo.Client
var conf = Config{}
var once sync.Once

// Init 初始化mongo数据库
func Init() {
	err := config.UnmarshalKey("mongo", &conf)
	if err != nil {
		logx.Fatalf("未发现mongo配置信息, %v", err)
	}
	if conf.Conn == "" {
		logx.Fatalf("mongo配置信息错误, %v", err)
	}
	New(conf)
}

// NewID 生成ObjectID
func NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// New 注册mongodb
func New(conf Config) *mongo.Client {
	wg := sync.WaitGroup{}
	var err error
	wg.Add(1)
	once.Do(func() {
		defer func() {
			wg.Done()
		}()
		if conf.Timeout <= 0 {
			conf.Timeout = 5
		}
		if conf.MaxPool <= 0 {
			conf.MaxPool = 10
		}
		clientOptions := options.Client().ApplyURI(conf.Conn).SetConnectTimeout(time.Duration(conf.Timeout) * time.Second).SetMaxPoolSize(uint64(conf.MaxPool))
		clientOptions.SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    conf.Database,
			Username:      conf.User,
			Password:      conf.Password,
		})
		instance, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			logx.Fatalf("创建连接mongo失败, %v", err)
		}
		err = instance.Ping(context.TODO(), nil)
		if err != nil {
			logx.Fatalf("连接mongo失败, %v", err)
		}
		logx.Info("mongo 注册成功...")
	})
	wg.Wait()
	return instance
}

// DB 获取数据库操作对象
func DB(name string) *mongo.Database {
	return instance.Database(name)
}

// Collection 获取集合操作对象
func Collection(key string) *mongo.Collection {
	return instance.Database(conf.Database).Collection(key)
}

// IndexExisted 检查索引是否存在
func IndexExisted(coll *mongo.Collection, index string) bool {
	curs, err := coll.Indexes().List(context.TODO(), nil)
	if err != nil {
		return true
	}
	for curs.Next(context.TODO()) {
		curr, _ := strconv.Unquote(curs.Current.Lookup("name").String())
		if index == curr {
			return true
		}
	}
	if curs.Err() != nil {
		return true
	}
	err = curs.Close(context.TODO())
	if err != nil {
		return true
	}
	return false
}

// Close 关闭Mongodb连接
func Close() {
	if instance != nil {
		err := instance.Disconnect(context.TODO())
		if err != nil {
			logx.Fatalf("关闭Mongodb连接失败 %v", err)
		}
	}
}
