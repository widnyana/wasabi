package http

// Config represents the configuration for the HTTP server.
type (
	Config struct {
		Port int    `mapstructure:"PORT"`
		Host string `mapstructure:"HOST" default:"127.0.0.1"`
	}
)
