// Package metrics Provide a simple intuitive use of prometheus metrics
package metrics

import (
	"net/http"
	"time"

	"github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rcrowley/go-metrics"
)

// HTTPHandler Return a fully configured http handler.
func HTTPHandler(
	metricsToRegister map[string]interface{},
	envName,
	serviceName string,
) http.Handler {
	// register the metrics
	registry := Register(metricsToRegister)

	// create a new prometheus providerthat updates every second with the
	// given parameters
	provider := prometheusmetrics.NewPrometheusProvider(
		registry,
		envName, serviceName,
		prometheus.DefaultRegisterer,
		1*time.Second,
	)

	// start the background updates
	go provider.UpdatePrometheusMetrics()

	return promhttp.Handler()
}

// Register Register required prometheus metrics and return a registry.
func Register(metricsToRegister map[string]interface{}) metrics.Registry {
	// iterate through the given metrics that should be registered
	for name, metric := range metricsToRegister {
		// register each metric
		metrics.GetOrRegister(name, metric)
	}

	// return the default registry
	return metrics.DefaultRegistry
}

// Inc Increase the metric with the given name.
func Inc(name string) {
	// fetch the metric by its name
	metric := metrics.Get(name)

	// resolve the metric type
	switch metric := metric.(type) {
	case metrics.Counter:
		// increase the counter if its a counter
		metric.Inc(1)
	case metrics.Gauge:
		// if the value is a gauge, first the current value needs to be
		// fetched...
		value := metric.Value()

		// ...and then increased by one
		metric.Update(value + 1)
	}
}

// Dec Decrease the metric with the given name.
func Dec(name string) {
	// fetch the metric with by its name
	metric := metrics.Get(name)

	// resolve the metric type
	switch metric := metric.(type) {
	case metrics.Counter:
		// decrease the counter if its a counter
		metric.Dec(1)
	case metrics.Gauge:
		// if the value is a gauge, first the current value needs to be
		// fetched...
		value := metric.Value()

		// ...and then decreased by one
		metric.Update(value - 1)
	}
}
