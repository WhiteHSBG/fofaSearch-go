package tools

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)



type fofa struct {
	username string
	key      string
	query    string
	size     string
	fields   string
	page     string
	UrlList  []string
}

type rst struct {
	Err bool `json:"error"`
	Model string `json:"mode"`
	Page int `json:"page"`
	Query string `json:"query"`
	Relst []result `json:"results"`
	Size int `json:"size"`
}

type result []string

func NewFofa() fofa  {
	return fofa{
		username: username ,
		key:      key,
		query:    "",
		size:     "10000",
		fields:   "ip,port,title,cert",
		page:     "1",
	}
}


func (f *fofa) Query(queryString string) {

	n:=time.Now()


	for i := 1; i <= mounthCount; i++ {
		b:=n.AddDate(0,-i+1,0).Format("2006-01-02")
		a:=n.AddDate(0,-i,0).Format("2006-01-02")
		f.parseResult(f.httpClient(queryString,a,b))//将返回的json字符串转为结构体并将urls存储到f.urls里
	}
	if  len(f.UrlList)>0 {
		fmt.Printf("数据查询完毕\n")
	}else {
		fmt.Printf("未查询到数据\n")
		os.Exit(0)
	}

}


func (f *fofa) parseResult(data []byte) *rst  {
	var r1 rst
	err:=json.Unmarshal(data,&r1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("获取数据: %v 条\n",len(r1.Relst))

	for _,rlst:= range r1.Relst{
		ip:=rlst[0]
		port:=rlst[1]
		title:=rlst[2]
		cert:=rlst[3]
		var strBuf string
		if cert == ""{
			strBuf="http://"+ip+":"+port+"|"+title+"\n"
		}else {
			strBuf="https://"+ip+":"+port+"|"+title+"\n"
		}
		f.UrlList =append(f.UrlList, strBuf)
	}

	return &r1
}

func (f *fofa) httpClient(queryString string,after string,befor string) []byte {
	f.query=fmt.Sprintf(queryString+"&& after=\"%s\" && before=\"%s\" && is_honeypot=false",after,befor)
	px,err:=url.Parse("http://127.0.0.1:8099")
	hc:=http.Client{
		Transport:     &http.Transport{
			Proxy: http.ProxyURL(px),
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10*time.Second,
	}

	//url结构体创建
	u,err:=url.Parse(baseUrl)
	if err != nil {
		fmt.Println("baseUrl err:",err)
		panic(err)
	}

	//组装参数
	data:=url.Values{}
	data.Set("email",f.username)
	data.Set("key",f.key)
	data.Set("size",f.size)
	data.Set("page",f.page)
	data.Set("fields",f.fields)
	//data.Set("after",after)
	//data.Set("before",befor)
	//data.Set("is_honeypot","false")
	if f.query == "" {
		panic(errors.New("query string is nil"))
	}
	data.Set("qbase64",f.queryDecode(f.query))
	u.RawQuery=data.Encode()//参数传入结构体
	fmt.Printf("--------------------------------------------------------------------------\n")
	fmt.Printf("请求接口与参数：\n%s\n",u.String())
	resp,err:=hc.Get(u.String())
	if err != nil {
		fmt.Println("fofa request err:",err)
		panic(err)
	}
	b,err:=ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return b
}

func (f fofa) queryDecode(queryString string) string {
	//return base64.URLEncoding.EncodeToString([]byte(queryString))
	return base64.StdEncoding.EncodeToString([]byte(queryString))
}
