package proxy

import (
	"context"
	"encoding/json"
	"search_proxy/internal/model/objs"
	"search_proxy/internal/util/log"
	"search_proxy/internal/util/tools"
	"strconv"
	"testing"
)

func TestAll(t *testing.T) {
	level := "debug"
	filePath := "/Users/wengguan/search_code/search_file/logs/proxy.log"
	maxSize := 128
	maxBackups := 100
	maxAge := 60
	compress := true
	log.InitLogger(level, filePath, maxSize, maxBackups, maxAge, compress)

	masters := []string{"127.0.0.1:7788", "127.0.0.1:7788"}
	slaves := make([][]string, 0)
	slave1 := []string{"127.0.0.1:7788", "127.0.0.1:7788"}
	slaves = append(slaves, slave1)
	slave2 := []string{"127.0.0.1:7788", "127.0.0.1:7788"}
	slaves = append(slaves, slave2)
	mode := "hash"
	timeout := 30
	NewPx(masters, slaves, timeout, mode)

	cxt := context.Background()
	var trackid uint64 = 12345
	cxt = context.WithValue(cxt, "trackid", trackid)
	remoteIP := "192.168.0.100"
	uri := "/add_doc"
	body := `{"body": "中国万里长城", "title": "五三班", "price": 5.30}`
	bodyByte := tools.Str2Bytes(body)

	retByte, _ := AddDoc(cxt, remoteIP, uri, bodyByte)
	var resp objs.RespData
	json.Unmarshal(retByte, &resp)
	t.Log(resp)

	uri = "/retrieve"
	retriveBody := `{"retreive_terms":[{"term":"长城","union":true,"inter":false}],"title_must":"五三班","price_start":5.1,"price_end":5.5}`
	retriveBodyByte := tools.Str2Bytes(retriveBody)
	repl, errString := RetrieveDoc(cxt, remoteIP, uri, retriveBodyByte)
	t.Log(repl, errString)

	id := resp.Result.Docid
	uri = "/del_doc?docid=" + strconv.FormatUint(id, 10)
	retByte, _ = DelDoc(cxt, remoteIP, uri)
	json.Unmarshal(retByte, &resp)
	t.Log(resp)

	uri = "/retrieve"
	repl, errString = RetrieveDoc(cxt, remoteIP, uri, retriveBodyByte)
	t.Log(repl, errString)
}
