package api

import (
	"github.com/gin-gonic/gin"
	"yan.site/ts_server/dao"
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
}

func NewApiManager(mysqlStorage *dao.MysqlStorage) *ApiManager {
	return &ApiManager{MysqlStorage: mysqlStorage}
}

func (a *ApiManager) Start() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("server-name","golang http")
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
	r.Run(":56") // listen and serve on 0.0.0.0:56
}
