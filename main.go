package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"yan.site/ts_server/api"
	"yan.site/ts_server/config"
	"yan.site/ts_server/dao"
	"yan.site/ts_server/handler"
	"yan.site/ts_server/model"
)

func main() {
	const version = "3.22.12.00"
	log.Println("VERSION TIMESTAMP: " + version)

	// Register prometheus metrics
	for _, c := range model.Metrics {
		prometheus.MustRegister(c)
	}

	// New data channel
	crawAndStorageChan := make(chan model.Record, config.GetConf().TsServerConfig.Crawl.BufferSize)
	crawAndMetricsChan := make(chan model.Record, config.GetConf().TsServerConfig.Crawl.BufferSize)
	defer close(crawAndStorageChan)

	// Crawl data
	crawlManager := handler.NewCrawlManager(crawAndStorageChan, crawAndMetricsChan)

	// Storage data
	mysqlStorage := dao.NewMysqlStorage(crawAndStorageChan)

	// metrics
	metricsManager := handler.NewMetricsManager(crawAndMetricsChan)

	// MQ
	receiveMq := handler.NewReceiveMQ(mysqlStorage, crawlManager)

	// Web api
	apiManager := api.NewApiManager(mysqlStorage, crawlManager)

	go crawlManager.Start()
	go apiManager.Start()
	go metricsManager.Start()
	go receiveMq.Start()
	// Register prometheus api according port in 57
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":57", nil)
	log.Println("http listen: prometheus metrics api has started")
	mysqlStorage.Start()
}
