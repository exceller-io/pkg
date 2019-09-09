package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusHandler Bootstraps prometheus for metrics collection
func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
