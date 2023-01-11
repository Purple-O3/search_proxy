package controller

import (
	"search_proxy/internal/model/proxy"
	"search_proxy/internal/objs"
	"search_proxy/internal/util/log"
	"testing"
)

func TestAll(t *testing.T) {
	config := objs.Config{
		objs.ServerConfig{IP: "127.0.0.1", Port: 7788, Debug: false, ReadTimeout: 1000, WriteTimeout: 1000, IdleTimeout: 1000},
		objs.LogConfig{},
		objs.RouterConfig{Model: "hash"},
		objs.GroupConfig{Timeout: 30, Masters: []string{"127.0.0.1:7788", "127.0.0.1:7799"}, Slaves: [][]string{{"127.0.0.1:7788", "127.0.0.1:7788"}, {"127.0.0.1:7799", "127.0.0.1:7799"}}},
	}

	proxy.NewProxyWrap(config.Group, config.Router)
	if err := StartNet(config.Server, closeFunc); err != nil {
		panic(err)
	}
}

func closeFunc() {
	log.CloseLogger()
}

/*
GET 127.0.0.1:7070/api/v1/del_doc?docid=0
GET 127.0.0.1:7070/api/v1/doc_isdel?docid=3
POST 127.0.0.1:7070/api/v1/add_doc
{"body": "浪漫巴黎土耳其", "title": "五零班", "price": 5.00}
{"body": "明朝那些事儿", "title": "五一班", "price": 5.10}
{"body": "银河英雄传说", "title": "五二班", "price": 5.20}
{"body": "中国万里长城", "title": "五三班", "price": 5.30}
{"body": "埃及金字塔", "title": "五四班", "price": 5.40}
POST 127.0.0.1:7070/api/v1/retrieve
{"retreive_terms":[{"term":"英雄","union":true,"inter":false},{"term":"埃及","union":true,"inter":false},{"term":"长城","union":false,"inter":true}],"title_must":"五三班","price_start":5.1,"price_end":5.5}
*/
