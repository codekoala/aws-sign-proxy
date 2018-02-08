package aws_sign_proxy

import "github.com/prometheus/client_golang/prometheus"

var (
	totalRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests proxied.",
	})
	requestTime = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "request_time",
		Help: "Total request time.",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})
)

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestTime)
}
