package main

import (
	_ "github.com/go-sql-driver/mysql"
	"yan.site/ts_server/api"
	"yan.site/ts_server/config"
	"yan.site/ts_server/crawl"
	"yan.site/ts_server/dao"
	"yan.site/ts_server/model"
)

func main() {

	dataChan := make(chan model.Record, config.GetConf().TsServerConfig.Crawl.BufferSize)
	defer close(dataChan)

	crawlManager := crawl.NewCrawlManager(dataChan)
	mysqlStorage := dao.NewMysqlStorage(dataChan)
	apiManager := api.NewApiManager(mysqlStorage)

	go crawlManager.Start()
	go apiManager.Start()
	mysqlStorage.Start()
}
