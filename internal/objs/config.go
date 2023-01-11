package objs

import (
	"search_proxy/internal/util/ginwrapper"
	"search_proxy/internal/util/log"
	"time"
)

type LogConfig = log.Config
type ServerConfig = ginwrapper.Config

type Config struct {
	Server ServerConfig
	Log    LogConfig
	Router RouterConfig
	Group  GroupConfig
}

type RouterConfig struct {
	Model string
}

type GroupConfig struct {
	Timeout time.Duration
	Masters []string
	Slaves  [][]string
}
