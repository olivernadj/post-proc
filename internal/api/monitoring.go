package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"strconv"
	"time"
)

func BuildSummaryVec(metricName string, metricHelp string, service string) *prometheus.SummaryVec {
	summaryVec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:      metricName,
			Help:      metricHelp,
			ConstLabels: prometheus.Labels{"service":service},
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"handler", "code"},
	)
	err := prometheus.Register(summaryVec)
	if err != nil {
		log.Println(err)
	}
	return summaryVec
}

// WithMonitoring optionally adds a middleware that stores request duration and response size into the supplied
// summaryVec
func WithMonitoring(next http.Handler, route Route, summary *prometheus.SummaryVec) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		lrw := NewMonitoringResponseWriter(rw)
		next.ServeHTTP(lrw, req)
		statusCode := lrw.statusCode
		duration := time.Since(start)

		// Store duration of request
		summary.WithLabelValues(route.Name, strconv.FormatInt(int64(statusCode), 10)).Observe(duration.Seconds() * 1000)
	})
}

type monitoringResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewMonitoringResponseWriter(w http.ResponseWriter) *monitoringResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &monitoringResponseWriter{w, http.StatusOK}
}

func (lrw *monitoringResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
