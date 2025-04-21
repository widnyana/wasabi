package tracing

import "time"

type Config struct {
	Enable     bool          `mapstructure:"enable"`
	Addr       string        `mapstructure:"addr"`
	Secure     bool          `mapstructure:"secure"`
	Timeout    time.Duration `mapstructure:"timeout"`
	SampleRate float64       `mapstructure:"sample_rate"`
	Name       string        `mapstructure:"name"`
}
