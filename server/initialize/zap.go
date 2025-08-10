package initialize

import (
	"os"
	"path/filepath"
	"server/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitZap() {
	zapConfig := global.Config.Zap
	logDir := filepath.Dir(zapConfig.Filename)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("failed to create log directory" + err.Error())
	}

	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   zapConfig.Filename,
		MaxSize:    zapConfig.MaxSize,
		MaxAge:     zapConfig.MaxAge,
		MaxBackups: zapConfig.MaxBackups,
	})
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(zapConfig.Level)); err != nil {
		*level = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer),
		level,
	)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	global.ZapLog = logger
}
