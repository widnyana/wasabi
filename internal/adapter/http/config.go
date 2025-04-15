package http

// Config represents the configuration for the HTTP server.
type (
	Config struct {
		Port int `mapstructure:"PORT"`
	}
)
