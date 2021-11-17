package proxy

import (
	"context"
	"encoding/json"
	"search_proxy/internal/util/log"
	"search_proxy/internal/util/tools"
	"strconv"
	"testing"
)

func TestAll(t *testing.T) {
	level := "debug"
	filePath := "../../../logs/proxy.log"
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
	body := `{"Ident":"88.199.1/bbb.def","Modified":"北京市首都机场","Saled":"成都","Num":13,"CreatedAt":"2021-11-02T16:42:21.199502+08:00"}`
	bodyByte := tools.Str2Bytes(body)

	retByte, _ := AddDoc(cxt, remoteIP, uri, bodyByte)
	var resp map[string]interface{}
	json.Unmarshal(retByte, &resp)
	t.Log(resp)

	uri = "/retrieve"
	retriveBody := `{"RetreiveTerms":[{"FieldName":"Modified","Term":"北京","TermCompareType":1,"Operator":"must"}],"Offset":0,"Limit":10}
`
	retriveBodyByte := tools.Str2Bytes(retriveBody)
	repl, count, errString := RetrieveDoc(cxt, remoteIP, uri, retriveBodyByte)
	t.Log(repl, count, errString)

	id := uint64(resp["docid"].(float64))
	uri = "/del_doc?docid=" + strconv.FormatUint(id, 10)
	retByte, _ = DelDoc(cxt, remoteIP, uri)
	json.Unmarshal(retByte, &resp)
	t.Log(resp)

	uri = "/retrieve"
	repl, count, errString = RetrieveDoc(cxt, remoteIP, uri, retriveBodyByte)
	t.Log(repl, count, errString)
}
