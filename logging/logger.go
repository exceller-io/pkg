package logging

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey struct{}

var fallbackLogger *zap.SugaredLogger

func init() {
	if logger, err := zap.NewProduction(); err != nil {
		fallbackLogger = zap.NewNop().Sugar()
	} else {
		fallbackLogger = logger.Named("fallback").Sugar()
	}
}

//WithLogger add the current logger to context and returns the context
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	loggerKey := &contextKey{}

	return context.WithValue(ctx, loggerKey, logger)
}

//LoggerFromContext returns the logger from the context
func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	loggerKey := &contextKey{}

	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return fallbackLogger
}

//NewLogger returns a new instance of logger
func NewLogger(level string, opts ...zap.Option) (*zap.SugaredLogger, zap.AtomicLevel) {
	loggingCfg := zap.NewProductionConfig()
	if len(level) > 0 {
		if level, err := levelFromString(level); err == nil {
			loggingCfg.Level = zap.NewAtomicLevelAt(*level)
		}
	}
	logger, err := loggingCfg.Build(opts...)
	if err != nil {
		panic(err)
	}
	return logger.Sugar(), loggingCfg.Level
}

//NewLoggerFromConfig returns a new instance of logger from Config
func NewLoggerFromConfig(configYaml string, levelOverride string, opts ...zap.Option) (*zap.SugaredLogger, zap.AtomicLevel, error) {
	if len(configYaml) == 0 {
		return nil, zap.AtomicLevel{}, errors.New("No logging configuration passed")
	}
	var loggingCfg zap.Config
	if err := yaml.Unmarshal([]byte(configYaml), &loggingCfg); err != nil {
		return nil, zap.AtomicLevel{}, err
	}
	if len(levelOverride) > 0 {
		if level, err := levelFromString(levelOverride); err == nil {
			loggingCfg.Level = zap.NewAtomicLevelAt(*level)
		}
	}
	logger, err := loggingCfg.Build(opts...)
	if err != nil {
		return nil, zap.AtomicLevel{}, err
	}
	logger.Info("Successfully created the logger.", zap.String("", configYaml))
	logger.Sugar().Infof("Logging level set to %v", loggingCfg.Level)
	return logger.Sugar(), loggingCfg.Level, nil
}

func levelFromString(level string) (*zapcore.Level, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid logging level: %v", level)
	}
	return &zapLevel, nil
}
