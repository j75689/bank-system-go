package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type LogFormat string

func (logFormat LogFormat) String() string {
	return string(logFormat)
}

const (
	JSONFormat    LogFormat = "json"
	ConsoleFormat LogFormat = "console"
)

var _ logger.Interface = (*Logger)(nil)

type Logger struct {
	zerolog.Logger
}

func (logger Logger) LogMode(logger.LogLevel) logger.Interface {
	return logger
}

func (logger Logger) Info(ctx context.Context, format string, args ...interface{}) {
	logger.Logger.Info().Msgf(format, args...)
}

func (logger Logger) Warn(tx context.Context, format string, args ...interface{}) {
	logger.Logger.Warn().Msgf(format, args...)
}

func (logger Logger) Error(tx context.Context, format string, args ...interface{}) {
	logger.Logger.Error().Msgf(format, args...)
}

func (logger Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, row := fc()
	logger.Logger.Trace().Dur("elapsed", elapsed).Str("sql", sql).Int64("row", row).Err(err).Msg("trace sql")
}

// NewLogger returns a Logger
func NewLogger(logLevel string, logFormat LogFormat) (Logger, error) {
	level := zerolog.InfoLevel
	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		return Logger{zerolog.Logger{}}, errors.Wrap(err, "init logger failed")
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	switch logFormat {
	case JSONFormat:
		return Logger{zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()}, nil
	case ConsoleFormat:
		return Logger{zerolog.New(os.Stdout).With().Caller().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})}, nil
	}

	return Logger{zerolog.Logger{}}, errors.Errorf("not support log format [%s]", logFormat)
}
