package main

import (
	"search_proxy/internal/controller"
	"search_proxy/internal/model/proxy"
	"search_proxy/internal/objs"
	"search_proxy/internal/util/log"
	"search_proxy/internal/util/viperwrapper"

	_ "go.uber.org/automaxprocs"
)

func start() {
	configPath := "../configs/proxy.yaml"
	var config objs.Config
	err := viperwrapper.DecodeConfig(configPath, &config)
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
