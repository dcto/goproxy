# Goproxy

* <strong>Git</strong>: <a href="https://gitee.com/sdox/goproxy">https://gitee.com/sdox/goproxy</a>

* <strong>Issues</strong>: <a href="https://gitee.com/sdox/goproxy/issues">https://gitee.com/sdox/goproxy/issues</a>

* <strong>License</strong>: MIT

* <strong>IRC</strong>: #Goproxy on freenode

___

#### 介绍
Goproxy Pool 
Golang 代理池系统，自动抓取代理，检测代理

<br />
启动：

```go run main.go```

<br />
编译：

```go build main.go```

<br />

>监听：0.0.0.0:5010

>访问：http://127.0.0.1:5010


#### 配置文件

```yaml
check:
  time: "*/10 * * * * ?"        #校验周期10秒
  http: http://www.166.com      #http效验地址
  https: https://www.baidu.com  #https效验地址
  lowest: 30                    #代理池伐值，少于30个
  timeout: 3                    #代理检测超时
redis:
  key: "proxies"                #Redis队列的KEY
  addr: 192.168.0.224:6379      #Redis连接地址

server:
  addr: 0.0.0.0:5010            #WEB服务监听地址:端口
  logs: ./logs                  #Log存储路径
```

  

  
#### 目录说明
```bash
├── app
│   ├── apis
│   │   └── api.go
│   └── jobs
│       ├── example.go
│       └── jobs.go
├── config
│   └── config.yaml
├── go.mod
├── LICENSE
├── main.go
├── pkg
│   ├── cache
│   │   └── redis.go
│   └── config
│       └── config.go
└── README.md
```

#### JOBS目录说明
* jobs.go >核心调度文件不建议修改
* example.go >任务抓取示例文件：

```go

//初始化任务
func init(){
	//增加计划任务，5秒种获取一次IP代理
	Cron("*/5 * * * * ?", Example)
}


//获取IP列表方法
func Example() {
	url := "http://www.example.com/getProxy.json"	
	
	if Count() < config.GetInt64("check.lowest") {
		resp := Http("Get", url, nil)

		if gjson.ValidBytes(resp) && gjson.GetBytes(resp, "code").Int() == 10001 {
			data := gjson.GetBytes(resp, "data.proxy_list")
			data.ForEach(func(key, value gjson.Result) bool {
				println(value.Get("ip").String()) 
				Save(value.Get("ip").String() +":"+ value.Get("port").String(), value.String(), value.Get("timeout").Int())
				return true
			})
		}
	}	
}

```

#### APIs目录说明
Web APIs的访问控制器，用户通过http获取代理，包含以下几个方法:  <br />

```javascript
GET    /                         --> Goproxy 代理池 Welcome 检测路径
GET    /get                      --> Goproxy/app/apis.Get 随机获取代理
GET    /pop                      --> Goproxy/app/apis.Pop 随机获取代理（并删除）
GET    /all                      --> Goproxy/app/apis.All 获取全部代理
GET    /raw                      --> Goproxy/app/apis.Raw 获取代理数据
```
