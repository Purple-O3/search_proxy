package request

import (
	"context"
	"search_proxy/internal/util/log"
	"search_proxy/internal/util/tools"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	level := "debug"
	filePath := "/Users/wengguan/search_code/search_file/logs/engine.log"
	maxSize := 128
	maxBackups := 100
	maxAge := 60
	compress := true
	log.InitLogger(level, filePath, maxSize, maxBackups, maxAge, compress)

	ctx := context.Background()
	timeout := 100
	ret, err := Get(ctx, "https://www.baidu.com", timeout)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(len(ret))
	}
	ret, err = Get(ctx, "https://www.google.com", timeout)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(len(ret))
	}
	body := `{"body": "浪漫巴黎土耳其", "title": "五零班", "price": 5.00}`
	bodyByte := tools.Str2Bytes(body)
	ret, err = Post(ctx, "http://127.0.0.1:7788/add_doc", "application/json", bodyByte, timeout)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(string(ret))
	}
	ret, err = Post(ctx, "https://www.google.com", "application/json", bodyByte, timeout)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(len(ret))
	}
}

/*
const retryInterval = 10
const retryCnt = 4
const failedThreshold = 8
*/

func TestBreaker(t *testing.T) {
	ctx := context.Background()
	timeout := 100
	NewBreaker()
	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Millisecond)
		if getState() == stateHalfOpen {
			DoCall("Get", ctx, "https://www.baidu.com", timeout)
		} else {
			DoCall("Get", ctx, "https://www.google.com", timeout)
		}
	}
}
