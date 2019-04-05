//Stolen from Knative
package logging

import (
	"context"

	"go.uber.org/zap"
)

var (
	fallbackLogger *zap.SugaredLogger
)

type loggerKey struct{}

func init() {
	if logger, err := zap.NewProduction(); err != nil {
		fallbackLogger = zap.NewNop().Sugar()
	} else {
		fallbackLogger = logger.Named("fallback").Sugar()
	}
}

//FromContext returns a logger from context
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey{}).(*zap.SugaredLogger); ok {
		return logger
	}
	return fallbackLogger
}
