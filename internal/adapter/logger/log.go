package logger

import (
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Encoder     string `mapstructure:"encoder" default:"console"`
	Level       string `mapstructure:"level" default:"INFO"`
	CallerDepth int    `mapstructure:"caller_depth" default:"0"`
}

const (
	callerDepthAdjustment = 0
)

var (
	jsonEncoderConfig = zapcore.EncoderConfig{
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

	consoleEncoderConfig = zapcore.EncoderConfig{
		TimeKey:          "ts",
		MessageKey:       "msg",
		LevelKey:         "level",
		NameKey:          "logger",
		StacktraceKey:    "stacktrace",
		ConsoleSeparator: "\t",
		FunctionKey:      zapcore.OmitKey,
		EncodeTime:       zapcore.RFC3339TimeEncoder,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
	}
)

// provideEncoder returns a zapcore.Encoder based on the provided logging mode.
//
// Supported modes:
//   - "json": Returns a JSON encoder using jsonEncoderConfig.
//   - "console" or any other value: Returns a Console encoder using consoleEncoderConfig.
//
// If the mode is unrecognized or empty, it defaults to the Console encoder.
func provideEncoder(mode string) zapcore.Encoder {
	var encoder zapcore.Encoder

	switch mode {
	case "json":
		encoder = zapcore.NewJSONEncoder(jsonEncoderConfig)
	case "console":
	default:
		encoder = zapcore.NewConsoleEncoder(consoleEncoderConfig)

	}

	return encoder
}

func newLogger(cfg Config, level zapcore.Level) *zap.Logger {
	core := zapcore.NewCore(
		provideEncoder(cfg.Encoder),
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

func levelFromString(s string) (zapcore.Level, error) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(s))
	return level, err
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
