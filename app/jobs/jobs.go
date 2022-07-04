package jobs

import (
	"context"
	"goproxy/pkg/cache"
	"goproxy/pkg/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/robfig/cron/v3"
)

var Crontab *cron.Cron = cron.New(cron.WithSeconds())

/**
 * 任务列表
 */
type Jobs interface {
	tong() map[string]int64
	cron() error
}

type Tong struct {}

/**
 * 保存IP
 * IP格式：IP:Port, 源数据, ttl过期时间（开发中）
 * example: 127.0.0.1:123456, origin string,  TTL exiprie time
 */
func Save(k string, v string, ttl int64){
	if err := cache.Redis.HSet(context.Background(), config.GetString("cache.key"), k, v).Err(); err != nil {
		log.Println("Cache to redis Error: "+ err.Error())
	}
}



/**
 * 检测IP
 */
func Ping() {
	var err error
	var timeout =  time.Duration(config.GetInt64("check.timeout")) * time.Second
	var proxy *url.URL
	var resp *http.Response
	proxies, err := cache.Redis.HGetAll(context.Background(), config.GetString("cache.key")).Result()	

	if err != nil {
		log.Println(err.Error())
	}

	for k, _ := range proxies {
		// log.Printf("Testing Proxy:`%s`", k)
		proxy, err = url.Parse("http://"+k)
		if err != nil {
			if err = cache.Redis.HDel(context.Background(), config.GetString("cache.key"), k).Err(); err != nil {
				log.Printf("Remove Invalid Proxy:`%s` Error: %s", k, err.Error())
			}else {
				log.Printf("Remove Invalid Proxy:`%s` Success", k)
			}
			continue
		}
		
		Client := &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}, 
			Timeout: timeout,
		}

		resp, err = Client.Get(config.GetString("check.http"))

		if err != nil || resp.StatusCode != http.StatusOK {
			resp, err = Client.Get(config.GetString("check.https")) 
			if err != nil || resp.StatusCode != http.StatusOK {
				c, err := cache.Redis.HDel(context.Background(), config.GetString("cache.key"), k).Result()
				log.Printf("Tested Proxy `%s` Error Removed [%d] Context: %v", k, c, err)
				continue
			}
		}

		log.Printf("Tested Proxy: `%s` Success", k)

		time.Sleep(1 * time.Second)
	}


    // // 解析代理地址
    // proxy, err := url.Parse(proxy_addr)
    // //设置网络传输
    // netTransport := &http.Transport{
    //     Proxy:                 http.ProxyURL(proxy),
    //     MaxIdleConnsPerHost:   10,
    //     ResponseHeaderTimeout: time.Second * time.Duration(5),
    // }
    // // 创建连接客户端
    // httpClient := &http.Client{
    //     Timeout:   time.Second * 10,
    //     Transport: netTransport,
    // }
    // begin := time.Now() //判断代理访问时间
    // // 使用代理IP访问测试地址
    // res, err := httpClient.Get(pingUrl)

    // if err != nil {
    //     log.Println(err)
    //     return
    // }
    // defer res.Body.Close()
    // speed := int(time.Now().Sub(begin).Nanoseconds() / 1000 / 1000) //ms
    // //判断是否成功访问，如果成功访问StatusCode应该为200
    // if res.StatusCode != http.StatusOK {
    //     log.Println(err)
    //     return
    // }
    // return speed, res.StatusCode

}


func Count() int64 {
	count, err := cache.Redis.HLen(context.Background(), config.GetString("cache.key")).Result()
	if err != nil {
		log.Println(err.Error())
	}
	return count
}


// 增加定时任务
// 每隔5秒执行一次：*/5 * * * * ?
// 每隔1分钟执行一次: 0 */1 * * * ?
// 每天23点执行一次: 0 0 23 * * ?
// 每天凌晨1点执行一次: 0 0 1 * * ?
// 每月1号凌晨1点执行一次: 0 0 1 1 * ?
// 每天的0点、13点、18点、21点都执行一次: 0 0 0,13,18,21 * * ?
func Cron(spec string, f func()) (cron.EntryID, error) {
	return Crontab.AddFunc(spec, f)
}



/**
 * HTTP请求
 * url: https://www.baidu.com
 * params: map[string]string
 */
 func Http(method string, url string, header map[string][]string) []byte {
	var b []byte
	var err error
	
	req, err := http.NewRequest(method, url, nil)
	req.Header = header
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Error Request From URL: %s", url)
		log.Println(err)
		return nil
	}
	
	b, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return b	
}