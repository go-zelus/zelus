package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-zelus/zelus/core/logx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-zelus/zelus/core/config"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoRows   = sql.ErrNoRows
	ErrConnDone = sql.ErrConnDone
	ErrTxDone   = sql.ErrTxDone
)

// Config mysql 配置信息
type Config struct {
	Name     string `mapstructure:"name"`     // 连接名称
	Host     string `mapstructure:"host"`     // 主机地址
	Port     int    `mapstructure:"port"`     // 连接端口号
	DB       string `mapstructure:"db"`       // 数据库名称
	User     string `mapstructure:"user"`     // 用户名
	Password string `mapstructure:"password"` // 密码
	MaxIDLE  int    `mapstructure:"max_idle"` // 最大空闲数
	MaxOpen  int    `mapstructure:"max_open"` // 最大连接数
	Timeout  int    `mapstructure:"timeout"`  // 连接超时
}

var instance map[string]*sqlx.DB = make(map[string]*sqlx.DB)

// Init 初始化mysql
func Init() {
	slice := make([]*Config, 0)
	err := config.UnmarshalKey("mysql", &slice)
	if err != nil {
		logx.Fatalf("解析配置失败 key:[%s], %v", "mysql", err)
	}
	if len(slice) == 0 {
		logx.Fatalf("未发现mysql配置信息, %v", err)
	}
	// 注册数据库
	for i, c := range slice {
		if i == 0 && c.Name == "" {
			c.Name = "default"
		}
		register(c.Name, "mysql", dataSource(c), c.MaxIDLE, c.MaxOpen)
	}
}

func dataSource(conf *Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&charset=utf8mb4&loc=%s&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DB, conf.Timeout, time.Local.String())
}

// register 注册mysql数据库
func register(aliasName, driverName, dataSource string, params ...int) {
	if aliasName == "" {
		logx.Fatalf("数据库别名为空%s", aliasName)
		return
	}
	db, err := sqlx.Open(driverName, dataSource)
	if err != nil {
		logx.Fatal(err)
	}
	for i, v := range params {
		switch i {
		case 0:
			db.SetMaxIdleConns(v)
		case 1:
			db.SetMaxOpenConns(v)
		}
	}
	instance[aliasName] = db
	err = db.Ping()
	if err != nil {
		logx.Fatal(err)
	}
	logx.Info("数据库 " + aliasName + " 注册成功...")
}

// ShowDataBase 展示数据库信息
func ShowDataBase() {
	logx.Infof("数据库: %v \n", instance)
}

// DB 根据名称获取数据库连接
func DB(key ...string) *sqlx.DB {
	if len(key) == 0 {
		return instance["default"]
	}
	return instance[key[0]]
}

// NewID 生成字符串唯一编号
func NewID() string {
	return primitive.NewObjectID().Hex()
}

// FieldNames 获取结构对应的数据库属性
func FieldNames(in interface{}) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("only accepts structs; got %T", v))
	}

	tp := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := tp.Field(i)
		if tag := fi.Tag.Get("db"); tag != "" {
			out = append(out, tag)
		}
	}
	return out
}

// FieldRows 组合各种条件字符串
func FieldRows(filedNames []string, sep string, removes ...string) string {
	if len(removes) == 0 {
		return strings.Join(filedNames, sep)
	}
	return strings.Join(remove(filedNames, removes...), sep)
}

// remove 删除数组中的制定字符串
func remove(strings []string, strs ...string) []string {
	out := append([]string(nil), strings...)
	for _, str := range strs {
		var n int
		for _, v := range out {
			if v != str {
				out[n] = v
				n++
			}
		}
		out = out[:n]
	}
	return out
}
