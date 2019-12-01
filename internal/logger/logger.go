package logger

import (
	"ImgCrop/internal/structs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func GetLogger(vp structs.Config) *zap.Logger {

	var level zapcore.Level
	switch vp.Logger.Level {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "ERROR":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.DebugLevel
	}
	outputs := make([]string, 0)
	outputs = append(outputs, "stderr")

	if vp.Logger.FileName == "" {
		dt := time.Now()
		vp.Logger.FileName = vp.Logger.Name + "_" + dt.Format("01_02_2006") + ".log"
	}
	outputs = append(outputs, vp.Logger.LogsPath+vp.Logger.FileName)

	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: outputs,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}
