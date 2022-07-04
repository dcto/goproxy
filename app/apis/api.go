package apis

import (
	"context"
	"goproxy/pkg/cache"
	"goproxy/pkg/config"
	"log"

	"github.com/gin-gonic/gin"
)

/**
 * 随机代理
 * @author dc
 * @version 20220620
 */
func Get(c *gin.Context) {

	data, err := cache.Redis.HRandField(context.Background(),  config.GetString("cache.key"), 1, false).Result()

	if err != nil {
		log.Println(err.Error())
	}
	c.String(200, data[0])
}

/**
 * 弹出单个并删除
 * @author dc
 * @version 20220704
 */
func Pop(c *gin.Context) {	
	data, err := cache.Redis.HRandField(context.Background(),  config.GetString("cache.key"), 1, false).Result()

	if err != nil {
		log.Println(err.Error())
	}

	if err := cache.Redis.HDel(context.Background(), config.GetString("cache.key"), data[0]).Err(); err != nil {
		log.Println(err.Error())
	}
	c.String(200, data[0])
}

/**
 * 获取全部
 * @author dc
 * @version 20220704
 */
func All(c *gin.Context) {
	data, err := cache.Redis.HKeys(context.Background(), config.GetString("cache.key")).Result()
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(200, data)
}


/**
 * 获取原始数据
 * @author dc
 * @version 20220704
 */
func Raw(c *gin.Context) {
	d, _ := cache.Redis.HGetAll(context.Background(), config.GetString("cache.key")).Result()
	c.JSON(200, d)
}