package proxy

import (
	"context"
	"encoding/json"
	"search_proxy/internal/model/objs"
	"search_proxy/internal/model/router"
	"search_proxy/internal/util/log"
	"search_proxy/internal/util/request"
	"search_proxy/internal/util/tools"
)

type proxy struct {
	masters []string
	slaves  [][]string
	timeout int
	mrt     router.Router
	srt     router.Router
}

func newProxy(groupMasters []string, groupSlaves [][]string, groupTimeout int, routerMode string) *proxy {
	px := new(proxy)
	px.masters = groupMasters
	px.slaves = groupSlaves
	px.timeout = groupTimeout
	px.mrt = router.RouterFactory(routerMode, len(groupMasters))
	px.srt = router.RouterFactory(routerMode, len(groupSlaves[0]))
	request.NewBreaker()
	return px
}

func (px *proxy) retrieveDoc(ctx context.Context, routerKey string, uri string, body []byte) (objs.RecallPostingList, string) {
	index := px.srt.LoadBalance(routerKey)
	errString := "nil"
	slavesLen := len(px.slaves)
	type goRet struct {
		repl objs.RecallPostingList
		err  error
	}
	retChan := make(chan goRet, slavesLen)
	for _, slave := range px.slaves {
		go func(slave []string) {
			url := "http://" + slave[index] + uri
			log.Infof("trackid:%d, url:%s, body:%s", ctx.Value("trackid").(uint64), url, tools.Bytes2Str(body))
			//retByte, err := request.DoCall("Post", ctx, url, "application/json", body, px.timeout)
			retByte, err := request.Post(ctx, url, "application/json", body, px.timeout)
			if err != nil {
				retChan <- goRet{nil, err}
				return
			}
			var resp objs.RespData
			err = json.Unmarshal(retByte, &resp)
			if err != nil {
				retChan <- goRet{nil, err}
				return
			}
			repl := resp.Result.Repl
			retChan <- goRet{repl, nil}
		}(slave)
	}

	errCnt := 0
	totalRepl := make(objs.RecallPostingList, 0)
	for i := 0; i < slavesLen; i++ {
		ret := <-retChan
		repl := ret.repl
		err := ret.err
		if err != nil {
			errCnt++
			errString = err.Error()
		} else {
			totalRepl = append(totalRepl, repl...)
		}
	}
	if errCnt == slavesLen {
		errString = "All server err: " + errString
	} else if errCnt > 0 {
		errString = "Some server err: " + errString
	}
	log.Infof("trackid:%d, repl:%v, err:%s", ctx.Value("trackid").(uint64), totalRepl, errString)
	return totalRepl, errString
}

func (px *proxy) addDoc(ctx context.Context, routerKey string, uri string, body []byte) ([]byte, error) {
	errString := "nil"
	index := px.mrt.LoadBalance(routerKey)
	url := "http://" + px.masters[index] + uri
	//retByte, err := request.DoCall("Post", ctx, url, "application/json", body, px.timeout)
	retByte, err := request.Post(ctx, url, "application/json", body, px.timeout)
	if err != nil {
		errString = err.Error()
	}
	log.Infof("trackid:%d, url:%s, body:%s, err:%s", ctx.Value("trackid").(uint64), url, tools.Bytes2Str(body), errString)
	return retByte, err
}

func (px *proxy) delDoc(ctx context.Context, routerKey string, uri string) ([]byte, error) {
	errString := "nil"
	index := px.mrt.LoadBalance(routerKey)
	url := "http://" + px.masters[index] + uri
	//retByte, err := request.DoCall("Get", ctx, url, px.timeout)
	retByte, err := request.Get(ctx, url, px.timeout)
	if err != nil {
		errString = err.Error()
	}
	log.Infof("trackid:%d, url:%s, err:%s", ctx.Value("trackid").(uint64), url, errString)
	return retByte, err
}

func (px *proxy) docIsDel(ctx context.Context, routerKey string, uri string) ([]byte, error) {
	errString := "nil"
	index := px.mrt.LoadBalance(routerKey)
	url := "http://" + px.masters[index] + uri
	//retByte, err := request.DoCall("Get", ctx, url, px.timeout)
	retByte, err := request.Get(ctx, url, px.timeout)
	if err != nil {
		errString = err.Error()
	}
	log.Infof("trackid:%d, url:%s, err:%s", ctx.Value("trackid").(uint64), url, errString)
	return retByte, err
}
