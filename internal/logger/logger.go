package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Log *zap.Logger

func Init() {
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/usr/local/kamailio/logs/kamailio_manage.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.FunctionKey = "FUNC"
	// Using your specified format
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2026-01-02 15:04:05")
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer),
		zap.InfoLevel)

	Log = zap.New(core, zap.AddCaller())
}
