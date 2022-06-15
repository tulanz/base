package providers

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLoggerProvider(vip *viper.Viper) *zap.Logger {

	environment := vip.GetString("app.environment")
	// var log *zap.Logger
	cores := []zapcore.Core{}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "logger_name",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
	}

	writer := zapcore.Lock(os.Stdout)
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writer, zap.DebugLevel)
	cores = append(cores, core)

	if environment != "" {
		//file
		logdir := "./logs"
		_ = os.MkdirAll(logdir, 0744)
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./logs/api.index.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     14,
			Compress:   false,
		})

		level := zap.AtomicLevel{}
		if err := level.UnmarshalText([]byte("info")); err != nil {
			level = zap.NewAtomicLevel() //默认用info级别
		}
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, level)
		cores = append(cores, core)
	}
	combinedCore := zapcore.NewTee(cores...)
	log := zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(0))
	zap.ReplaceGlobals(log)

	return log
}
