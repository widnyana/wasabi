package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

// EnableTracing enables OpenTelemetry tracing for GORM database operations.
func EnableTracing(ctx context.Context, db *gorm.DB, logger *otelzap.Logger) error {
	_, span := otel.Tracer("postgres").Start(ctx, "enable-tracing")
	defer span.End()

	logger.Debug("enabling gorm tracing plugin")
	return db.Use(tracing.NewPlugin())
}

// HookConnection sets up lifecycle hooks for the database connection.
func HookConnection(lifecycle fx.Lifecycle, sqlDB *sql.DB) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return sqlDB.PingContext(ctx) },
		OnStop:  func(_ context.Context) error { return sqlDB.Close() },
	})
}

// configureConnPool configures the database connection pool settings during application startup.
func configureConnPool(lc fx.Lifecycle, sqlDB *sql.DB, cfg Config, logger *otelzap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			configurePool(ctx, sqlDB, logger, cfg)
			return nil
		},
	})
}

// configurePool configures the connection pool settings for the provided SQL database.
// It sets the maximum number of idle connections, the maximum lifetime of a connection,
// the maximum idle time of a connection, and the maximum number of open connections.
// If the corresponding configuration values are zero, default values are used and logged.
func configurePool(
	ctx context.Context,
	sqlDB *sql.DB,
	logger *otelzap.Logger,
	config Config,
) {
	_, span := otel.Tracer("postgres").Start(ctx, "configure-pool")
	defer span.End()

	// Apply all connection pool settings in a single function using a helper
	applyPoolSetting(sqlDB, config, logger)
}

// applyPoolSetting applies all database connection pool settings
func applyPoolSetting(sqlDB *sql.DB, config Config, logger *otelzap.Logger) {
	logger.Info("applying database connection pool settings")

	// Set MaxIdleConns
	maxIdleConns := config.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = defaultMaxIdleConns
		logger.Debug("using default value for max_idle_conns", zap.Int("default", defaultMaxIdleConns))
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// Set ConnMaxLifetime
	maxLifetime := time.Duration(config.ConnMaxLifetimeMillis) * time.Millisecond
	if config.ConnMaxLifetimeMillis == 0 {
		maxLifetime = defaultConnMaxLifetime
		logger.Debug("using default value for conn_max_lifetime", zap.Duration("default", defaultConnMaxLifetime))
	}
	sqlDB.SetConnMaxLifetime(maxLifetime)

	// Set ConnMaxIdleTime
	maxIdleTime := time.Duration(config.ConnMaxIdleTimeMillis) * time.Millisecond
	if config.ConnMaxIdleTimeMillis == 0 {
		maxIdleTime = defaultConnMaxIdleTime
		logger.Debug("using default value for conn_max_idle_time", zap.Duration("default", defaultConnMaxIdleTime))
	}
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	// Set MaxOpenConns
	maxOpenConns := config.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = defaultMaxOpenConns
		logger.Debug("using default value for max_open_conns", zap.Int("default", defaultMaxOpenConns))
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
}
