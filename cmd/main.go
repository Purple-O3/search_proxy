package main

import (
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"search_proxy/internal/controller/customnet"
	"search_proxy/internal/model/proxy"
	"search_proxy/internal/util/log"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"
)

var cn customnet.Net

func init() {
	filePath := "../configs/proxy.toml"
	fileName := filepath.Base(filePath)
	fileNames := strings.Split(fileName, ".")
	if len(fileNames) != 2 {
		panic(errors.New("fileNames len not equal 2"))
	}

	vp := viper.New()
	vp.SetConfigName(fileNames[0])
	vp.SetConfigType(fileNames[1])
	vp.AddConfigPath(filepath.Dir(filePath))
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	logLevel := viper.GetString("log.level")
	logFilePath := viper.GetString("log.file_path")
	logMaxSize := viper.GetInt("log.max_size")
	logMaxBackups := viper.GetInt("log.max_backups")
	logMaxAge := viper.GetInt("log.max_age")
	logCompress := viper.GetBool("log.compress")
	log.InitLogger(logLevel, logFilePath, logMaxSize, logMaxBackups, logMaxAge, logCompress)

	routerMode := viper.GetString("router.mode")
	groupTimeout := viper.GetInt("group.timeout")
	groupMasters := make([]string, 0)
	groupSlaves := make([][]string, 0)
	i := 0
	for {
		group := "group" + strconv.Itoa(i)
		if viper.IsSet(group) {
			master := viper.GetString(group + ".master")
			groupMasters = append(groupMasters, master)
			slave := viper.GetStringSlice(group + ".slave")
			groupSlaves = append(groupSlaves, slave)
		} else {
			break
		}
		i++
	}
	proxy.NewPx(groupMasters, groupSlaves, groupTimeout, routerMode)

	ip := viper.GetString("server.ip")
	port := viper.GetString("server.port")
	cn = customnet.NetFactory("http")
	cn.StartNet(ip, port)
	log.Infof("server start!!!")
}

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}

func closeServer() {
	log.CloseLogger()
	cn.Shutdown()
}

func main() {
	listenSignal()
	closeServer()
}
