package gslog

import (
	"github.com/metagogs/gogs/global"
	"go.uber.org/zap"
)

var devLog = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Development:      true,
	Encoding:         "console",
	EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
}

var prodLog = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
	Development:      false,
	Sampling:         nil,
	Encoding:         "json",
	EncoderConfig:    zap.NewProductionEncoderConfig(),
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
}

var disableLog = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.FatalLevel),
	Development:      false,
	Sampling:         nil,
	Encoding:         "json",
	EncoderConfig:    zap.NewProductionEncoderConfig(),
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
}

func NewLog(name string) *zap.Logger {
	if global.GOGS_DISABLE_LOG {
		// only show the fatal log
		logger, _ := disableLog.Build()
		return logger.Named("gogs_" + name)
	}

	if global.GoGSDebug {
		logger, _ := devLog.Build()
		return logger.Named("gogs_" + name)
	}

	logger, _ := prodLog.Build()

	return logger.Named("gogs_" + name)
}
