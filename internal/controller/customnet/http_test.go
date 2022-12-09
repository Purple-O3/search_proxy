package customnet

import (
	"search_proxy/internal/model/proxy"
	"search_proxy/internal/util/log"
	"testing"
)

func TestAll(t *testing.T) {
	masters := []string{"127.0.0.1:7788", "127.0.0.1:7799"}
	slaves := make([][]string, 0)
	slave1 := []string{"127.0.0.1:7788", "127.0.0.1:7788"}
	slaves = append(slaves, slave1)
	slave2 := []string{"127.0.0.1:7799", "127.0.0.1:7799"}
	slaves = append(slaves, slave2)
	mode := "hash"
	timeout := 30
	proxy.NewPx(masters, slaves, timeout, mode)
	ip, port := "", "7070"
	cn := NetFactory("http")
	cn.StartNet(ip, port)
	/*c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c*/
	log.CloseLogger()
	cn.Shutdown()
}

/*
GET 127.0.0.1:7070/del_doc?docid=0
GET 127.0.0.1:7070/doc_isdel?docid=3
POST 127.0.0.1:7070/add_doc
{"body": "浪漫巴黎土耳其", "title": "五零班", "price": 5.00}
{"body": "明朝那些事儿", "title": "五一班", "price": 5.10}
{"body": "银河英雄传说", "title": "五二班", "price": 5.20}
{"body": "中国万里长城", "title": "五三班", "price": 5.30}
{"body": "埃及金字塔", "title": "五四班", "price": 5.40}
POST 127.0.0.1:7070/retrieve
{"retreive_terms":[{"term":"英雄","union":true,"inter":false},{"term":"埃及","union":true,"inter":false},{"term":"长城","union":false,"inter":true}],"title_must":"五三班","price_start":5.1,"price_end":5.5}
*/
