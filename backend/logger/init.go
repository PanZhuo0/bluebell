package logger

import (
	"backend/settings"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(conf *settings.LogConfig) {
	var cores zapcore.Core
	var confCore zapcore.Core //the core which setted by configuration file

	// Get LogLevel from config file
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(conf.Level)); err != nil {
		panic(err)
	}
	confCore = zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), //Encoder
		getWriteSync(conf), //WriteSyncer
		level)              //LogLevel
	cores = confCore

	// if Mode="dev"
	var devCore zapcore.Core
	if settings.Conf.Mode == "dev" {
		devCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			os.Stdout,
			zapcore.DebugLevel)
		cores = zapcore.NewTee(confCore, devCore)
	}
	logger := zap.New(cores, zap.WithCaller(true))
	zap.ReplaceGlobals(logger)

	// just a test
	zap.L().Info("Logger Init Success")
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
