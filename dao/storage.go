package dao

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"yan.site/ts_server/config"
	"yan.site/ts_server/model"
)

type MysqlStorage struct {
	DataChan chan model.Record
	Db       *sql.DB
	config   config.TsServerConfig
}

func NewMysqlStorage(dadaChan chan model.Record) *MysqlStorage {
	return &MysqlStorage{
		DataChan: dadaChan,
	}
}

func (s *MysqlStorage) Start() {
	s.init()
	s.storageData()
}

func (s *MysqlStorage) init() {
	s.config = config.GetConf().TsServerConfig

	db, err := sql.Open(s.config.DataBase.DrvierName, s.getDataSourceName())
	if err != nil {
		log.Println("database: connect failure")
		panic("crash")
	} else {
		db.SetMaxOpenConns(int(s.config.DataBase.SetMaxOpenConns))
		db.SetMaxIdleConns(int(s.config.DataBase.SetMaxIdleConns))
		db.Ping()
	}
	s.Db = db
	log.Println("database: init success")
}

func (s *MysqlStorage) storageData() {
	for temp := range s.DataChan {
		go s.SaveRecord(temp)
	}
}

func (s *MysqlStorage) SaveRecord(record model.Record) bool {

	jsonTextByte, err := json.Marshal(record)
	if err != nil {
		log.Println("database: generate jsonText failure")
		return false
	}

	rows, err2 := s.Db.Query("INSERT INTO record(id,trace_id,parent_id,json_text) VALUES (?,?,?,?)", record.Id, record.TraceId, record.ParentId, string(jsonTextByte))
	defer rows.Close()
	if err2 != nil {
		log.Println("database: save record failure,id: " + record.Id)
		println(err2.Error())
		return false
	}
	log.Println("database: save a record,id: " + record.Id)
	return true
}

func (s *MysqlStorage) getDataSourceName() string {

	user := s.config.DataBase.User
	pwd := s.config.DataBase.Pwd
	url := s.config.DataBase.Url
	db := s.config.DataBase.Db
	return user + ":" + pwd + "@tcp(" + url + ")/" + db + "?charset=utf8"
}

func (s *MysqlStorage) GetRecordByTraceId(traceId string) ([]model.Record, bool) {
	rows, err := s.Db.Query("SELECT * FROM record WHERE trace_id = ?", traceId)
	if err != nil {
		log.Println("database: get record error," + err.Error())
		return nil, false
	}
	var recordResultSet []model.Record
	for rows.Next() {
		var id string
		var traceId string
		var parentId string
		var jsonText string
		err = rows.Scan(&id, &traceId, &parentId, &jsonText)
		if err != nil {
			log.Println("database: prase return data error," + err.Error())
			return nil, false
		}
		var record model.Record
		err2 := json.Unmarshal([]byte(jsonText), &record)
		if err2 != nil {
			log.Println("database: json.Unmarshal return data error," + err.Error())
			return nil, false
		}
		recordResultSet = append(recordResultSet, record)
	}
	return recordResultSet, true
}
