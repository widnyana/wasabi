package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	defaultMaxIdleConns    int           = 10
	defaultMaxOpenConns    int           = 50
	defaultConnMaxLifetime time.Duration = 2 * time.Hour / time.Millisecond
	defaultConnMaxIdleTime time.Duration = 1 * time.Hour / time.Millisecond
)

// Config represents the configuration for the PostgreSQL database.
type Config struct {
	DSN                  string ` envconfig:"dsn"`
	Debug                bool   `envconfig:"debug"`
	SlowThresholdMS      int    `envconfig:"slow_threshold_ms"`
	PreferSimpleProtocol bool   `envconfig:"prefer_simple_protocol"`

	ConnMaxIdleTimeMillis int `envconfig:"conn_max_idle_time_millis" default:"60000"`  // Default to 1 minute in milliseconds
	ConnMaxLifetimeMillis int `envconfig:"conn_max_lifetime_millis" default:"7200000"` // Default to 2 hours in milliseconds
	MaxIdleConns          int `envconfig:"max_idle_conns" default:"10"`
	MaxOpenConns          int `envconfig:"max_open_conns" default:"50"`
}

// NewGorm initializes a new GORM database connection for PostgreSQL.
// It takes the application's PostgreSQL configuration and a logger instance.
// It configures the GORM connection with specified settings, including logging level
// based on the debug flag, prepared statements, disabled default transactions,
// and a custom logger adapter that integrates with Zap and OpenTelemetry.
// It also retrieves the underlying `sql.DB` instance to configure the connection pool
// using the provided configuration and performs an initial connection check.
// Returns a GORM database instance or an error if the connection fails.
func NewGorm(config Config, logger *otelzap.Logger) (*gorm.DB, error) {
	ctx, span := otel.Tracer("postgres").Start(context.TODO(), "new-gorm")
	defer span.End()

	logger.Ctx(ctx).Info("initializing gorm postgresql connection")

	level := gormlogger.Error
	if config.Debug {
		level = gormlogger.Info
	}

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  config.DSN,
			PreferSimpleProtocol: config.PreferSimpleProtocol,
		}),
		&gorm.Config{
			AllowGlobalUpdate:      false,
			PrepareStmt:            true, // https://gorm.io/docs/performance.html#SQL-Builder-with-PreparedStmt
			SkipDefaultTransaction: true, // https://gorm.io/docs/performance.html#Disable-Default-Transaction
			FullSaveAssociations:   false,
			Logger: NewZapLoggerAdapter(logger, gormlogger.Config{
				LogLevel:                  level,
				SlowThreshold:             time.Duration(config.SlowThresholdMS) * time.Millisecond,
				ParameterizedQueries:      true, // Don't include params in the SQL log,
				IgnoreRecordNotFoundError: true, // Ignore ErrRecordNotFound error for logger
			}),
		},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewSQLDB extracts the underlying *sql.DB instance from the provided GORM database connection.
// This allows direct access to the standard Go SQL driver's functionality, such as
// connection pool management and lower-level database operations if needed.
// It returns the *sql.DB instance or an error if it cannot be retrieved.
func NewSQLDB(orm *gorm.DB) (*sql.DB, error) { return orm.DB() }
