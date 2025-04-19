package logger

import (
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level       string `mapstructure:"level" default:"INFO"`
	CallerDepth int    `mapstructure:"caller_depth" default:"0"`
}

var Module = fx.Options(
	// GetLogger provides an *otelzap.Logger instance. This logger is
	// configured with JSON encoding, writing to standard output, and
	// using the Info level for logging. It also integrates with
	// OpenTelemetry for tracing context.
	fx.Provide(GetLogger),

	// This anonymous function provides an *otelzap.SugaredLogger.
	// A SugaredLogger offers a more convenient interface for logging
	// without explicit formatting, using methods like Infof, Debugf, etc.
	// It's derived from the base *otelzap.Logger.
	fx.Provide(func(logger *otelzap.Logger) *otelzap.SugaredLogger {
		return logger.Sugar()
	}),

	// This anonymous function provides the underlying *zap.Logger.
	// This is the core Zap logger instance that provides structured logging
	// capabilities with fields. The *otelzap.Logger wraps this *zap.Logger
	// to add OpenTelemetry integration.
	fx.Provide(func(logger *otelzap.Logger) *zap.Logger {
		return logger.Logger
	}),

	// WithLogger configures Fx's internal event logging to use the
	// provided *otelzap.Logger. It creates a fxevent.ZapLogger,
	// sets its underlying Zap Logger, and then sets the log level for
	// Fx's events to DebugLevel, ensuring verbose output of Fx lifecycle events.
	fx.WithLogger(func(logger *otelzap.Logger) fxevent.Logger {
		l := &fxevent.ZapLogger{Logger: logger.Logger}
		l.UseLogLevel(zap.DebugLevel)
		return l
	}),
)

const callerDepthAdjustment = 0

func newLogger(cfg Config, level zapcore.Level) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		StacktraceKey:  "stacktrace",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		os.Stdout,
		level,
	)

	return zap.New(core, zap.AddCaller()).
		WithOptions(
			zap.WithCaller(true),
			zap.AddStacktrace(zap.ErrorLevel),
			zap.AddCallerSkip(cfg.CallerDepth),
		)
}

func GetLogger(cfg Config) (*otelzap.Logger, error) {
	level, err := levelFromString(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	logger := otelzap.New(newLogger(cfg, level),
		otelzap.WithErrorStatusLevel(zapcore.WarnLevel),
		otelzap.WithMinLevel(level),
		otelzap.WithCaller(true),
		otelzap.WithCallerDepth(callerDepthAdjustment),
		otelzap.WithMinLevel(level),
		otelzap.WithErrorStatusLevel(zapcore.WarnLevel),
	)

	_ = otelzap.ReplaceGlobals(logger)

	return logger, nil
}

func levelFromString(s string) (zapcore.Level, error) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(s))
	return level, err
}
