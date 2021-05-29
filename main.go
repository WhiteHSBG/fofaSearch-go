package main

import (
	"flag"
	"fofaSearch/tools"
	"os"
)




func main() {
	qu:=flag.String("query","","查询语句参数，建议使用''包起来")
	outPut:=flag.String("output","","输出文件名")
	flag.Parse()
	if *qu == "" {
		flag.Usage()
		os.Exit(0)
	}

	tools.ChickConfig()
	tools.NewConf("./config.yaml")
	fofa:=tools.NewFofa()
	fofa.Query(*qu)
	checker:=tools.NewChecker(&fofa.UrlList)
	if *outPut != "" {
		checker.OutPutFile=*outPut
	}
	checker.StartCheck()
}

