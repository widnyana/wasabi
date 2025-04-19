package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/widnyana/wasabi/internal/adapter/database/pg"
	"github.com/widnyana/wasabi/internal/adapter/http"
	"github.com/widnyana/wasabi/internal/adapter/logger"
	"github.com/widnyana/wasabi/internal/adapter/metrics"
	"github.com/widnyana/wasabi/internal/adapter/redis"
	"github.com/widnyana/wasabi/internal/adapter/tracing"
	"github.com/widnyana/wasabi/internal/constant"
)

var ErrParseConfigFailed = errors.New("failed to parse configuration from environment variables")

// AppConfig contains structure of Application Config
type AppConfig struct {
	Env      string         `envconfig:"env"`
	HTTP     http.Config    `envconfig:"http"`
	Redis    redis.Config   `envconfig:"redis"`
	Postgres pg.Config      `envconfig:"postgres"`
	Metrics  metrics.Config `envconfig:"metrics"`
	Tracing  tracing.Config `envconfig:"tracing"`
	Log      logger.Config  `envconfig:"log"`
}

// NewAppConfig Provide a configuration instance
// Usage:
//
//	fx.Options(
//		fx.Provide(NewAppConfig),
//	)
func NewAppConfig() *AppConfig {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Can't load the configuration. Error: %s", err.Error())
	}

	return cfg
}

// LoadConfig load configuration from environment variables
func LoadConfig() (*AppConfig, error) {
	var cfg AppConfig
	if err := envconfig.Process(constant.AppName, &cfg); err != nil {
		return &cfg, errors.Join(
			err,
			ErrParseConfigFailed,
		)
	}

	return &cfg, nil
}

// PrintBanner print application banner
func PrintBanner(c *AppConfig) {
	if strings.EqualFold(c.Env, "production") {
		return
	}

	fmt.Println(constant.AppBanner)
}
