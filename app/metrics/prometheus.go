package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/pieterclaerhout/go-log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

type Prometheus struct {
	router      *gin.Engine
	metric      *prometheus.HistogramVec
	MetricsPath string
}

func NewPrometheus(router *gin.Engine) *Prometheus {
	metric := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_requests",
		Help: "Http requests metric",
	},
		[]string{"status", "uri"},
	)

	p := &Prometheus{
		metric:      metric,
		MetricsPath: "/metrics",
		router:      router,
	}

	return p
}

func (p *Prometheus) Use(e *gin.Engine) {
	e.Use(p.handlerFunc())
	p.setMetricsPath()
}

func (p *Prometheus) handlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		diff := time.Since(start)
		status := strconv.Itoa(c.Writer.Status())
		uri := c.FullPath()
		obs, err := p.metric.GetMetricWith(map[string]string{"status": status, "uri": uri})

		if err != nil {
			log.Error(err)
		} else {
			obs.Observe(diff.Seconds())
		}
	}
}

func (p *Prometheus) setMetricsPath() {
	p.router.GET(p.MetricsPath, prometheusHandler())
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
