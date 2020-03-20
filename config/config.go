package config

import (
	"github.com/winjeg/go-commons/conf"
	"sync"
)

const configFile = "config.yml"

// 配置文件结构
type Config struct {
	TsServerConfig  TsServerConfig  `yaml:"ts_server"`
	AppServerConfig AppServerConfig `yaml:"app_server"`
}

// ts服务配置
type TsServerConfig struct {
	Crawl    Crawl    `yaml:"crawl"`
	DataBase DataBase `yaml:"data_base"`
}
type Crawl struct {
	TimeInterval int64 `yaml:"time_interval"`
	BufferSize   int64 `yaml:"buffer_size"`
}

type DataBase struct {
	DrvierName      string `yaml:"drvier_name"`
	Url             string `yaml:"url"`
	User            string `yaml:"user"`
	Pwd             string `yaml:"pwd"`
	Db              string `yaml:"db"`
	SetMaxOpenConns int8   `yaml:"set_max_open_conns"`
	SetMaxIdleConns int8   `yaml:"set_max_idle_conns"`
}

// 需要被抓取的应用
type AppServerConfig struct {
	Hosts []Host `yaml:"host"`
}
type Host struct {
	Ip          string `yaml:"ip"`
	Port        int    `yaml:"port"`
	ProjectName string `yaml:"project_name"`
}

// 全局变量
var (
	once      sync.Once
	configure *Config
)

func GetConf() *Config {
	if configure != nil {
		return configure
	} else {
		once.Do(getConf)
	}
	return configure
}

func getConf() {
	configure = new(Config)
	err := conf.Yaml2Object(configFile, &configure)
	if err != nil {
		panic(err)
	}
}
