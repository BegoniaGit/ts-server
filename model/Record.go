package model

type TraceReceive struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data []Record `json:"data"`
}

type Record struct {
	TraceId        string            `json:"traceId"`
	ParentId       string            `json:"parentId"`
	LastId         string            `json:"lastId"`
	Id             string            `json:"id"`
	StartTimeStamp int               `json:"startTimeStamp"`
	DurationTime   int               `json:"durationTime"`
	Error          bool              `json:"error"`
	Name           string            `json:"name"`
	ServerName     string            `json:"serverName"`
	Stage          string            `json:"stage"`
	SamplingRate   float32            `json:"samplingRate"`
	NotePair       []NotePair        `json:"notePair"`
	AdditionalPair map[string]string `json:"additionalPair"`
}

type NotePair struct {
	NoteName  string `json:"noteName"`
	TimeStamp int    `json:"timeStamp"`
	Host      Host   `json:"host"`
}

type Host struct {
	ServerName string `json:"serverName"`
	Address    string `json:"address"`
	Port       int    `json:"port"`
}
