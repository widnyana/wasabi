package pg

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

// ZapLoggerAdapter is a GORM logger that uses Zap for logging.
type ZapLoggerAdapter struct {
	lgr *otelzap.Logger
	cfg gormlogger.Config
}

// NewZapLoggerAdapter creates a new ZapLoggerAdapter.
// It takes a Zap logger and a GORM logger configuration.
// It returns a GORM logger interface.
func NewZapLoggerAdapter(lgr *otelzap.Logger, cfg gormlogger.Config) gormlogger.Interface {
	return &ZapLoggerAdapter{
		lgr: lgr,
		cfg: cfg,
	}
}

// Error implements logger.Interface.
func (z *ZapLoggerAdapter) Error(_ context.Context, msg string, args ...interface{}) {
	if z.cfg.LogLevel >= gormlogger.Info {
		z.lgr.Sugar().Errorf(msg, args...)
	}
}

// Info implements logger.Interface.
func (z *ZapLoggerAdapter) Info(_ context.Context, msg string, args ...interface{}) {
	if z.cfg.LogLevel >= gormlogger.Info {
		z.lgr.Sugar().Infof(msg, args...)
	}
}

// Warn implements logger.Interface.
func (z *ZapLoggerAdapter) Warn(_ context.Context, msg string, args ...interface{}) {
	if z.cfg.LogLevel >= gormlogger.Info {
		z.lgr.Sugar().Warnf(msg, args...)
	}
}

// Trace implements logger.Interface.
func (z *ZapLoggerAdapter) Trace(
	_ context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if z.cfg.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.Error(err),
		zap.String("location", fileWithLineNum()),
		zap.String("elapsed", elapsed.String()),
		zap.String("sql", sql),
		zap.Int64("rows_affected", rows),
	}

	switch {
	case err != nil && z.cfg.LogLevel >= gormlogger.Error:
		z.lgr.Error("", fields...)
	case elapsed >= z.cfg.SlowThreshold && z.cfg.LogLevel >= gormlogger.Warn:
		z.lgr.Warn("", fields...)
	case z.cfg.LogLevel >= gormlogger.Info:
		z.lgr.Info("", fields...)
	}
}

// LogMode sets the log level for the logger.
// It returns a new logger with the specified log level.
func (z *ZapLoggerAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	z.cfg.LogLevel = level
	return z
}

func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !strings.Contains(file, "gorm.io") && !strings.HasSuffix(file, "_test.go") {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	return ""
}
