package logx

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-zelus/zelus/core/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once
var core zapcore.Core
var sugar *zap.SugaredLogger
var instance *Logger

// Config 日志配置信息
type Config struct {
	FileName   string `mapstructure:"FileName"`
	MaxSize    int    `mapstructure:"MaxSize"`
	MaxAge     int    `mapstructure:"MaxAge"`
	Level      string `mapstructure:"Level"`
	SizeIncise bool   `mapstructure:"SizeIncise"`
}

type Logger struct {
	sugar *zap.SugaredLogger
}

func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		FunctionKey:    zapcore.OmitKey,
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}

//  Init 初始化日志
func Init() {
	conf := Config{}
	config.UnmarshalKey("log", &conf)
	New(&conf)
}

// New 创建日志对象
func New(conf *Config) *Logger {
	if instance == nil {
		instance = initLog(conf)
	}
	return instance
}

func initLog(conf *Config) *Logger {
	once.Do(func() {
		if conf.MaxSize == 0 {
			conf.MaxSize = 100
		}
		if strings.TrimSpace(conf.FileName) == "" {
			conf.FileName = "logs/app.log"
		} else {
			idx := strings.LastIndex(conf.FileName, "/")
			if idx+1 <= utf8.RuneCountInString(conf.FileName) {
				suf := conf.FileName[idx+1:]
				if strings.TrimSpace(suf) == "" {
					conf.FileName += "log"
				}
			}
		}
		if strings.TrimSpace(conf.Level) == "" {
			conf.Level = "debug"
		}
		if conf.SizeIncise {
			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   conf.FileName,
				MaxSize:    conf.MaxSize,
				MaxBackups: 5,
				MaxAge:     conf.MaxAge,
				Compress:   false,
			})
			core = zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfig()),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w),
				getLevel(conf.Level),
			)
			sugar = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

		} else {
			logLevel := getLevel(conf.Level)
			infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl < zapcore.WarnLevel && lvl >= logLevel
			})
			warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.WarnLevel && lvl <= logLevel
			})
			infoWriter := getWriter(conf.FileName)
			warnWriter := getWriter(conf.FileName + ".error")
			encoder := zapcore.NewConsoleEncoder(encoderConfig())
			core = zapcore.NewTee(
				zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
				zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
				zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), getLevel(conf.Level)),
			)
			sugar = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel)).Sugar()
		}
		instance = &Logger{sugar: sugar}
	})
	return instance
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.Print(err)
	}
	return hook
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Panic(args ...interface{}) {
	sugar.Panic(args...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Warnf(template, args)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

// Sync 程序退出前执行，避免部分内容没有写入磁盘
func Sync() {
	sugar.Sync()
}
