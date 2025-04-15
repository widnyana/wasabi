package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	lblIP     = "ip"
	lblMethod = "method"
	lblPath   = "path"
	lblStatus = "status"
)

// PromProbe is a probe that collects metrics for HTTP requests.
type PromProbe struct {
	req *prometheus.CounterVec
}

type PromHabitProbe struct{}

// NewPromProbe creates a new PromProbe.
func NewPromProbe() *PromProbe {
	return &PromProbe{
		req: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_request_total",
				Help: "Total number of application requests",
			},
			[]string{lblIP, lblMethod, lblPath, lblStatus},
		),
	}
}

// Middleware is a middleware that logs HTTP requests.
func (probe *PromProbe) Middleware(ctx *fiber.Ctx) error {
	defer probe.LogReq(ctx)
	return ctx.Next()
}

// LogReq logs an HTTP request.
func (probe *PromProbe) LogReq(ctx *fiber.Ctx) {
	labels := prometheus.Labels{
		lblIP:     ctx.IP(),
		lblMethod: ctx.Route().Method,
		lblPath:   ctx.Route().Path,
		lblStatus: strconv.Itoa(ctx.Response().StatusCode()),
	}

	probe.req.With(labels).Inc()
}
