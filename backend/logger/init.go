package logger

import (
	"backend/settings"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(conf *settings.LogConfig) {
	level := new(zapcore.Level)
	err := level.UnmarshalText([]byte(conf.Level))
	if err != nil {
		panic(err)
	}
	confCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), //Encoder
		getWriteSync(conf), //WriteSyncer
		level)              //LogLevel

	// if 配置文件中Mode="dev"
	var devCore zapcore.Core
	if settings.Conf.Mode == "dev" {
		devCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), //Encode
			os.Stdout,          //WriteSyncer
			zapcore.DebugLevel) //LogLevel
	}
	core := zapcore.NewTee(confCore, devCore)
	logger := zap.New(core, zap.WithCaller(true))
	zap.ReplaceGlobals(logger)

	// 输出测试
	zap.L().Debug("Logger Init Success")
}

func getWriteSync(s *settings.LogConfig) zapcore.WriteSyncer {
	l := lumberjack.Logger{
		Filename:   s.FileName,
		MaxSize:    s.MaxSize,
		MaxAge:     s.MaxAge,
		MaxBackups: s.MaxBackups,
	}
	return zapcore.AddSync(&l)
}
