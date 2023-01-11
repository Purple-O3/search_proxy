package request

import (
	"context"
	"search_proxy/internal/util/tools"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	ctx := context.Background()
	timeout := time.Duration(100)
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
	timeout := time.Duration(100)
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
