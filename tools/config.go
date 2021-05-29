package tools

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

var (
	username     string
	key          string
	baseUrl      ="https://fofa.so/api/v1/search/all"
	mounthCount  int
	thread       int
	proxy        *url.URL
	timeout      int
	resultdetail bool
)

type Config struct {
	Username string `yaml:username`
	Key      string `yaml:key`
	Thread   int    `yaml:thread`
	Month    int    `yaml:month`
	Proxy    string `yaml:proxy`
	Timeout int `yaml:timeout`
	Resultdetail bool   `yaml:resultdetail`
}

var conf *Config

func NewConf(path string) *Config {
	cfgf,err:=os.Open(path)
	if err != nil {
		fmt.Printf("%s\n",err)
		panic(err)
	}
	by,err:=ioutil.ReadAll(cfgf)
	if err != nil {
		fmt.Printf("%s\n",err)
		panic(err)
	}
	yaml.Unmarshal(by,&conf)

	username=conf.Username
	key=conf.Key
	mounthCount=conf.Month
	thread=conf.Thread
	proxy,_=url.Parse(conf.Proxy)
	timeout= conf.Timeout
	resultdetail=conf.Resultdetail

	return conf
}

func ChickConfig() {
	cupath,_:=os.Getwd()
	cfpath:=path.Join(cupath,"config.yaml")
	_,err:=os.Stat(cfpath)
	if err != nil {
		fmt.Printf("未检测到配置文件，配置文件已写入%s，请修改配置文件内容后再运行。",cfpath)
		file,err1:=os.OpenFile(cfpath,os.O_TRUNC|os.O_APPEND|os.O_WRONLY|os.O_CREATE,0666)
		if err1 != nil {
			fmt.Printf("配置文件写入失败")
		}
		cfgByts,_ := Asset("config.yaml")
		file.Write(cfgByts)
		file.Close()
		os.Exit(0)
	}else {
		fmt.Printf("检测到配置文件：%s\n",cfpath)
	}
}