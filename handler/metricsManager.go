package handler

import (
	"log"
	"strconv"
	"strings"
	"yan.site/ts_server/model"
)

type MetricsManager struct {
	CrawAndMetricsChan chan model.Record
}

func NewMetricsManager(crawAndMetricsChan chan model.Record) *MetricsManager {
	return &MetricsManager{CrawAndMetricsChan: crawAndMetricsChan,}
}

func (c *MetricsManager) Start() {
	log.Println("metrics data: start deal metrics data")
	for temp := range c.CrawAndMetricsChan {
		c.dealMetrics(temp)
	}
}

func (c *MetricsManager) dealMetrics(record model.Record) {

	name := strings.Split(record.Name, ".")
	recordType := strings.ToLower(name[0])
	baseLabels := map[string]string{"serverName": record.ServerName, "stage": record.Stage, "error": strconv.FormatBool(record.Error)}
	switch recordType {
	case "http":
		{
			baseLabels["method"] = name[1]
			baseLabels["path"] = record.AdditionalPair["path"]
			baseLabels["statusCode"] = record.AdditionalPair["status code"]
			model.WebApiResponseCount.With(baseLabels).Inc()
			model.WebApiResponseCost.With(baseLabels).Observe(float64(record.DurationTime))
			log.Println("metrics data: add a record to metrics,type: http,id: " + record.Id)
			break
		}
	case "client":
		{
			baseLabels["method"] = name[1]
			baseLabels["path"] = record.AdditionalPair["path"]
			baseLabels["statusCode"] = record.AdditionalPair["status code"]
			baseLabels["remoteServer"] = record.AdditionalPair["remote server"]
			model.HttpClientApiCount.With(baseLabels).Inc()
			model.HttpClientApiCost.With(baseLabels).Observe(float64(record.DurationTime))
			log.Println("metrics data: add a record to metrics,type: client,id: " + record.Id)
			break
		}
	case "mysql":
		{
			baseLabels["method"] = name[1]
			baseLabels["mysqlServer"] = record.AdditionalPair["mysql name"]
			model.MysqlRequestCount.With(baseLabels).Inc()
			model.MysqlRequestCost.With(baseLabels).Observe(float64(record.DurationTime))
			log.Println("metrics data: add a record to metrics,type: mysql,id: " + record.Id)
			break
		}
	}
}
