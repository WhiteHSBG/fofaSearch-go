# fofaSearch
使用go实现的fofa搜索批量工具 需要高级会员
- 主要功能：
    - fofa搜索
    - 自动探测搜索结果是否存活
## Usage
```
Usage of ./fofaSearch:
  -output string
        输出文件名
  -query string
        查询语句参数，建议使用''包起来
```

```
./fofaSearch -query 'title="beijing"' -output beijing.txt 
```

## config.yaml
第一次运行会自动生成配置文件
```
username: test@qq.com //fofa邮箱
key: APIKEY //fofa apiKey
month: 24 //从当前时间向前查询多少个月的数据
thread: 3000 //探测存活的线程数
timeout: 10 //探测存活超时时间
resultdetail: false //是否输出详细信息
proxy: //探测存活代理
```
