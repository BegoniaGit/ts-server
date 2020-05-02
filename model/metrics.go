package model

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DefaultBucket       = []float64{1, 2, 4, 8, 16, 32, 64, 128, 512, 1024}
	ConstantLabels      = []string{"serverName", "stage", "error"}
	WebHttpResponseCost = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ts_http_seconds_bucket",
			Help:    "web http response time in seconds",
			Buckets: DefaultBucket,
		},
		append(ConstantLabels, "method", "path", "statusCode"),
	)
	WebHttpResponseCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts_http_num_count",
			Help: "web http cache result count",
		},
		append(ConstantLabels, "method", "path", "statusCode"),
	)
	HttpClientApiCost = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ts_client_api_seconds_bucket",
			Help:    "ts_http_api_seconds_bucket",
			Buckets: DefaultBucket,
		},
		append(ConstantLabels, "method", "path", "statusCode", "remoteServer"),
	)
	HttpClientApiCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts_client_api_num_count",
			Help: "ts_client_api_num_count",
		},
		append(ConstantLabels, "method", "path", "statusCode", "remoteServer"),
	)
	MysqlRequestCost = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ts_mysql_seconds_bucket",
			Help:    "ts_mysql_seconds_bucket",
			Buckets: DefaultBucket,
		},
		append(ConstantLabels, "method", "mysqlServer"),
	)
	MysqlRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts_mysql_num_count",
			Help: "ts_mysql_num_count",
		},
		append(ConstantLabels, "method", "mysqlServer"),
	)

	SamplingRateGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ts_sampling_rate",
			Help: "ts_sampling_rate",
		},
		[]string{"serverName", "stage"},
	)

	Metrics = []prometheus.Collector{
		WebHttpResponseCost,
		WebHttpResponseCount,
		HttpClientApiCost,
		HttpClientApiCount,
		MysqlRequestCost,
		MysqlRequestCount,
		SamplingRateGauge,
	}
)
