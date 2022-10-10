package jobs

import (
	"goproxy/pkg/config"
	"github.com/tidwall/gjson"
)

//初始化任务
func init(){
	//5秒种获取一次IP地址
	//Cron("*/5 * * * * ?", Example)
}


//获取IP列表
func Example() {
	url := "http://www.example.com/getProxy.json"	
	
	if Count() < config.GetInt64("check.lowest") {
		resp := Http("Get", url, nil)

		if gjson.ValidBytes(resp) && gjson.GetBytes(resp, "code").Int() == 10001 {
			data := gjson.GetBytes(resp, "data.proxy_list")
			data.ForEach(func(key, value gjson.Result) bool {
				// println(value.Get("ip").String()) 
				Save(value.Get("ip").String() +":"+ value.Get("port").String(), value.String(), value.Get("timeout").Int())
				return true
			})
		}
	}	
}
