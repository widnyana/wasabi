package redis

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	grds "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Config is the redis configuration.
type (
	Config struct {
		Addr     string `mapstructure:"addr"`
		DB       int    `mapstructure:"db"`
		Password string `mapstructure:"password"`
		Enable   bool   `mapstructure:"enable" default:"false"`
	}

	// HealthChecker is the redis health checker.
	HealthChecker struct {
		client *grds.Client
	}
)

// Module is the fx module for redis.
func Module(cfg Config) fx.Option {
	if !cfg.Enable {
		return fx.Options()
	}

	return fx.Module("redis",
		fx.Provide(
			func() Config { return cfg },
			NewClient,
			NewHealthChecker,
		),
		fx.Invoke(
			HookRedis,
			EnableTracing,
			EnableMetrics,
		),
	)
}

// NewClient creates a new redis client.
func NewClient(cfg Config) (*grds.Client, error) {
	if !cfg.Enable {
		return nil, ErrRedisDisabled
	}

	return grds.NewClient(&grds.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	}), nil
}

// HookRedis hooks the redis client to the fx lifecycle.
func HookRedis(lifecycle fx.Lifecycle, redis *grds.Client) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return redis.Ping(ctx).Err()
		},
		OnStop: func(_ context.Context) error {
			return redis.Close()
		},
	})
}

// EnableTracing enables tracing for the redis client.
func EnableTracing(redis *grds.Client) error {
	return redisotel.InstrumentTracing(redis)
}

// EnableMetrics enables metrics for the redis client.
func EnableMetrics(redis *grds.Client) error {
	return redisotel.InstrumentMetrics(redis)
}

// NewHealthChecker creates a new redis health checker.
func NewHealthChecker(client *grds.Client) HealthChecker {
	return HealthChecker{client}
}

// CheckHealth checks the health of the redis client.
func (healthChecker HealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx).Err()
}
