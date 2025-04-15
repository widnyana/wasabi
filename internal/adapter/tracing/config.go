package tracing

import "time"

type Config struct {
	Addr       string        `envconfig:"addr"`
	Secure     bool          `envconfig:"secure"`
	Timeout    time.Duration `envconfig:"timeout"`
	SampleRate float64       `envconfig:"sample_rate"`
}
