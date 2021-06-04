package tools

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type checker struct {
	targetList   *[]string
	resultStList []resultSt
	targetCh     chan string
	OutPutFile   string
	outDetail    bool
}

type resultSt struct {
	url string
	code int
	title string
}

func NewChecker(urlList *[]string) checker {
	return checker{
		targetList: urlList,
		resultStList: []resultSt{},
		targetCh: make(chan string,len(*urlList)),
	}
}

func (c *checker) StartCheck() {

	tmpUrlList:=[]string{}//去重后的url

	//url去重
	TmpUrls := make(map[string]bool)
	for _, url := range *c.targetList {
		_,ok:=TmpUrls[url]//如果TmpUrls中不存在url这个key，则ok为false
		if !ok{
			TmpUrls[url] = true
			if url != "" {
				tmpUrlList=append(tmpUrlList, url)//将去重后的url存入
			}
		}
	}
	fmt.Printf("共查询到数据：%d条\n",len(tmpUrlList))

	//启动线程池开始访问，将访问结果存入checker.resultStList
	for _,url := range tmpUrlList {
		c.targetCh <- url//url存入通道
	}
	close(c.targetCh)
	for i := 0; i < thread; i++ {
		wg.Add(1)
		go c.curl()
	}
	wg.Wait()
	c.writeToOutput()
}

func (c checker) writeToOutput()  {
	var fileName string
	if c.OutPutFile ==""{
		fileName=time.Now().Format("200601021504")+".txt"
	}else {
		fileName=c.OutPutFile
	}
	pt,_:=os.Getwd()
	pt=path.Join(pt,fileName)
	fmt.Printf("文件输出路径：%s",pt)
	f,err:=os.OpenFile(pt,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if !resultdetail {
		for _, st := range c.resultStList {
			stBuff:=fmt.Sprintf("%s\n",st.url)
			fmt.Fprintf(f,stBuff)
		}
	}else {
		for _, st := range c.resultStList {
			stBuff:=fmt.Sprintf("%s|%d|%s\n",st.url,st.code,st.title)
			fmt.Fprintf(f,stBuff)
		}
	}


}

func (c *checker) curl() {
	defer wg.Done()

	//制造一个client
	var hc http.Client
	if proxy.Host!="" {
		hc=http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
			Timeout:   time.Duration(timeout)*time.Second,
		}
	}else {
		hc=http.Client{
			Transport: &http.Transport{
				IdleConnTimeout:     15 * time.Second,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				TLSHandshakeTimeout: 5 * time.Second,
				DisableKeepAlives:   false,
			},
			Timeout:   time.Duration(timeout)*time.Second,
		}
	}



	for true {
		target:=<-c.targetCh//从关闭的通道取值，如果通道空了就会返回对应值的空值
		if target == "" {//空了就退出
			break
		}
		url:=strings.Split(target,"|")[0]
		title:=strings.Split(target,"|")[1]
		title=strings.Replace(title,"\n","",-1)//去除所有\n

		resp,err:=hc.Get(url)
		if err != nil {
			fmt.Printf("check:%s:fail\n",url)
			//fmt.Printf("\nrequest %s err: %s\n",,err)
			continue
		}
		fmt.Printf("check:%s:success\n",url)
		//b,err:=ioutil.ReadAll(resp.Body)

		c.resultStList=append(c.resultStList, resultSt{
			url:   url,
			code:  resp.StatusCode,
			title: title,
		})
	}
}