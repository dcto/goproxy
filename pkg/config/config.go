package config

import (
	"log"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func init(){

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	log.Println("load config file with :" + viper.ConfigFileUsed())
}


func Get(path string, defaultValue ...interface{}) interface{} {
    if !viper.IsSet(path) {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return nil
    }
    return viper.Get(path)
}

// Env 读取环境变量，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
    if len(defaultValue) > 0 {
        return Get(envName, defaultValue[0])
    }
    return Get(envName)
}

// Add 新增配置项
func Set(name string, configuration map[string]interface{}) {
    viper.Set(name, configuration)
}


// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
    return cast.ToString(Get(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
    return cast.ToInt(Get(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
    return cast.ToInt64(Get(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
    return cast.ToUint(Get(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
    return cast.ToBool(Get(path, defaultValue...))

}