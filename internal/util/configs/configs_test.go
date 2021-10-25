package configs

import (
	"testing"

	"github.com/spf13/viper"
)

func TestAll(t *testing.T) {
	viper.SetConfigName("configs")
	viper.SetConfigType("toml") // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")    // 查找配置文件所在的路径
	viper.SetDefault("redis.port", 6381)
	t.Log("redis port: ", viper.Get("redis.port"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	t.Log("app_name:", viper.Get("app_name"))
	t.Log("log_level:", viper.Get("log_level"))
	t.Log("redis ip: ", viper.Get("redis.ip"))
	t.Log("redis port: ", viper.Get("redis.port"))
}
