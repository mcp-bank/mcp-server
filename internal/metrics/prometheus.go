package metrics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var mcpToolCallsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mcp_tool_calls_total",
		Help: "total number of mcp tool calls",
	},
	[]string{"tool"},
)

var mcpToolErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mcp_tool_errors_total",
		Help: "total number of mcp tool errors",
	},
	[]string{"tool"},
)

var mcpToolDurationSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "mcp_tool_duration_seconds",
		Help:    "mcp tools durations in seconds ",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"tool"},
)

func register() {
	prometheus.MustRegister(
		mcpToolCallsTotal,
		mcpToolErrorsTotal,
		mcpToolDurationSeconds,
	)
}

func NewServer() *http.Server {
	register()
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	server := &http.Server{
		Addr:         ":9091", //TODO fix hardcode
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func RecordToolCall(tool string) {
	mcpToolCallsTotal.WithLabelValues(tool).Inc()
}

func RecordToolCallError(tool string) {
	mcpToolErrorsTotal.WithLabelValues(tool).Inc()
}

func RecordToolDuration(tool string, duration time.Duration) {
	mcpToolDurationSeconds.WithLabelValues(tool).Observe(duration.Seconds())
}
