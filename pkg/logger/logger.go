package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill"
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

type Logger struct {
	zerolog.Logger
}

func (l Logger) WarpedGormLogger() GormLogger {
	return GormLogger{l.Logger}
}

func (l Logger) WrapedWatermillLogger() WatermillLogger {
	return WatermillLogger{l.Logger}
}

var _ logger.Interface = (*GormLogger)(nil)

type GormLogger struct {
	logger zerolog.Logger
}

func (logger GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return logger
}

func (logger GormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	logger.logger.Info().Msgf(format, args...)
}

func (logger GormLogger) Warn(tx context.Context, format string, args ...interface{}) {
	logger.logger.Warn().Msgf(format, args...)
}

func (logger GormLogger) Error(tx context.Context, format string, args ...interface{}) {
	logger.logger.Error().Msgf(format, args...)
}

func (logger GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, row := fc()
	logger.logger.Trace().Dur("elapsed", elapsed).Str("sql", sql).Int64("row", row).Err(err).Msg("trace sql")
}

var _ watermill.LoggerAdapter = (*WatermillLogger)(nil)

type WatermillLogger struct {
	logger zerolog.Logger
}

func (logger WatermillLogger) Error(msg string, err error, fields watermill.LogFields) {
	logger.logger.Error().Err(err).Fields(fields).Msg(msg)
}

func (logger WatermillLogger) Info(msg string, fields watermill.LogFields) {
	logger.logger.Info().Fields(fields).Msg(msg)
}

func (logger WatermillLogger) Debug(msg string, fields watermill.LogFields) {
	logger.logger.Debug().Fields(fields).Msg(msg)
}

func (logger WatermillLogger) Trace(msg string, fields watermill.LogFields) {
	logger.logger.Trace().Fields(fields).Msg(msg)
}

func (logger WatermillLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return logger
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
