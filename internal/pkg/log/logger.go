package log

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
}

func New(level string, serviceName string) Logger {
	zerolog.TimeFieldFormat = "02-01-2006 15:04:05"
	zerolog.LevelWarnValue = "warning"
	zerolog.MessageFieldName = "description"

	isNotParsed := false
	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		isNotParsed = true
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	if isNotParsed {
		logger.Warn().Msgf("loglevel %s not parsed, defaulting to info", level)
	}

	return &logger
}
