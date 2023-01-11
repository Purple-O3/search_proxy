package main

import (
	"errors"
	"path/filepath"
	"search_proxy/internal/controller"
	"search_proxy/internal/model/proxy"
	"search_proxy/internal/objs"
	"search_proxy/internal/util/log"
	"strings"

	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"
)

func start() {
	configPath := "../configs/proxy.toml"
	fileName := filepath.Base(configPath)
	fileNames := strings.Split(fileName, ".")
	if len(fileNames) != 2 {
		panic(errors.New("fileNames len not equal 2"))
	}
	var config objs.Config
	vp := viper.New()
	vp.AddConfigPath(filepath.Dir(configPath))
	vp.SetConfigName(fileNames[0])
	vp.SetConfigType(fileNames[1])
	err := vp.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = vp.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	if config.Log.Type == "file" {
		log.InitLogger(config.Log)
	}
	proxy.NewProxyWrap(config.Group, config.Router)
	if err = controller.StartNet(config.Server, closeFunc); err != nil {
		panic(err)
	}
}

func closeFunc() {
	log.CloseLogger()
}

func main() {
	start()
}
