package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"yan.site/ts_server/dao"
	"yan.site/ts_server/handler"
	"yan.site/ts_server/model"
)

func NewRecordApiResp(code int8, msg string, data []model.Record) *RecordApiResp {
	return &RecordApiResp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

type RecordApiResp struct {
	Code int8
	Msg  string
	Data []model.Record
}

type ApiManager struct {
	MysqlStorage *dao.MysqlStorage
	CrawlManager *handler.CrawlManager
}

func NewApiManager(mysqlStorage *dao.MysqlStorage, crawlManager *handler.CrawlManager) *ApiManager {
	return &ApiManager{MysqlStorage: mysqlStorage,
		CrawlManager: crawlManager}
}

func (a *ApiManager) Start() {
	r := gin.Default()

	// router '/' to front end page
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/graph")
	})

	// front end page router
	r.StaticFile("/graph", "./static/graph.html")

	// get record data
	r.GET("/api/records", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("server-name", "golang server")
		traceId := c.Query("traceId")
		if traceId == "" {
			c.PureJSON(200, NewRecordApiResp(-1, "Incoming parameters are incorrect", nil))
		} else {
			data, ok := a.MysqlStorage.GetRecordByTraceId(traceId)
			if !ok {
				c.PureJSON(200, NewRecordApiResp(1, "No data found", nil))
			} else {
				c.PureJSON(200, NewRecordApiResp(0, "ok", data))
			}
		}
	})

	// report records
	r.POST("/report", func(c *gin.Context) {
		var records []model.Record
		err := c.BindJSON(&records)
		if err != nil {
			log.Println("api: get body error")
			c.PureJSON(500, NewRecordApiResp(-1, "server has error", nil))
		} else {
			a.CrawlManager.SaveData(records...)
			log.Println("api: report success,count: " + strconv.Itoa(len(records)))
		}
	})

	// listen and serve on 0.0.0.0:56
	r.Run(":56")
}
