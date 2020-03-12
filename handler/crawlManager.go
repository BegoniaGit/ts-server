package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"yan.site/ts_server/config"
	"yan.site/ts_server/model"
)

var (
	timeout = time.Second * 60
)

type CrawlManager struct {
	CrawAndStorageChan chan model.Record
	CrawAndMetricsChan chan model.Record
}

func NewCrawlManager(crawAndStorageChan chan model.Record, crawAndMetricsChan chan model.Record) *CrawlManager {
	return &CrawlManager{
		CrawAndStorageChan: crawAndStorageChan,
		CrawAndMetricsChan: crawAndMetricsChan,
	}
}

func (m *CrawlManager) Start() {
	log.Println("crawl record: start crawl data")
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			m.crawlData()
		}
	}
}

func (m *CrawlManager) crawlData() {
	HostList := config.GetConf().AppServerConfig.Hosts
	data, success := m.getTraceData(HostList)
	if success {
		for _, v := range data {
			m.CrawAndStorageChan <- v
			m.CrawAndMetricsChan <- v
		}
		log.Println("crawl record: " + strconv.Itoa(len(data)) + " records")
	} else {

	}
}

func (m *CrawlManager) getTraceData(hostList []config.Host) ([]model.Record, bool) {

	var recordResultSet []model.Record
	for _, v := range hostList {
		records, ok := m.getTraceDataByHost(v)
		if !ok {
			log.Println(v.Ip + "crawl record: get trace data failure")
		} else {
			recordResultSet = append(recordResultSet, records...)
		}
	}
	return recordResultSet, true
}

func (m *CrawlManager) getTraceDataByHost(host config.Host) ([]model.Record, bool) {
	url := "http://" + host.Ip + ":" + strconv.Itoa(host.Port) + "/trace"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Println("crawl record: create req failed", err)
		return nil, false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false
	}
	if resp.StatusCode >= 300 {
		log.Println("crawl record: httpRequest failed", "code: ", resp.StatusCode, "msg: ", resp.Status)
	}
	var respData model.TraceReceive
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("crawl record: read body failed", err)
		return nil, false
	}
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		log.Println("crawl record: deserialize body failed", err)
		return nil, false
	}
	if respData.Code != 0 {
		return nil, false
	}
	return respData.Data, true
}
