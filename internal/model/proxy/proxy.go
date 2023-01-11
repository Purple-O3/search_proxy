package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"search_proxy/internal/model/router"
	"search_proxy/internal/objs"
	"search_proxy/internal/util/log"
	"search_proxy/internal/util/request"
	"search_proxy/internal/util/tools"
	"sort"
	"time"
)

type proxy struct {
	masters []string
	slaves  [][]string
	timeout time.Duration
	mrt     router.Router
	srt     router.Router
}

func newProxy(groupConfig objs.GroupConfig, routerConfig objs.RouterConfig) *proxy {
	px := new(proxy)
	px.masters = groupConfig.Masters
	px.slaves = groupConfig.Slaves
	px.timeout = groupConfig.Timeout
	px.mrt = router.RouterFactory(routerConfig.Model, len(groupConfig.Masters))
	px.srt = router.RouterFactory(routerConfig.Model, len(groupConfig.Slaves[0]))
	request.NewBreaker()
	return px
}

func (px *proxy) retrieveDoc(ctx context.Context, routerKey string, uri string, body []byte) (objs.RecallPostingList, int, error) {
	index := px.srt.LoadBalance(routerKey)
	errString := "nil"
	slavesLen := len(px.slaves)
	type goRet struct {
		repl  objs.RecallPostingList
		count int
		err   error
	}
	retChan := make(chan goRet, slavesLen)
	for _, slave := range px.slaves {
		go func(slave []string) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("%v", err)
				}
			}()

			url := "http://" + slave[index] + uri
			log.Infof("trackid:%v, url:%s, body:%s", ctx.Value("trackid"), url, tools.Bytes2Str(body))
			//retByte, err := request.DoCall("Post", ctx, url, "application/json", body, px.timeout)
			retByte, err := request.Post(ctx, url, "application/json", body, px.timeout)
			if err != nil {
				retChan <- goRet{nil, 0, err}
				return
			}
			var retData objs.RetData
			err = json.Unmarshal(retByte, &retData)
			if err != nil {
				retChan <- goRet{nil, 0, err}
				return
			}
			if retData.Code != 0 {
				retChan <- goRet{nil, 0, errors.New("engine return err")}
				return
			}
			retChan <- goRet{retData.Result, retData.Count, nil}
		}(slave)
	}

	errCnt := 0
	totalCount := 0
	totalRepl := make(objs.RecallPostingList, 0)
	for i := 0; i < slavesLen; i++ {
		ret := <-retChan
		repl := ret.repl
		count := ret.count
		err := ret.err
		if err != nil {
			errCnt++
			errString = err.Error()
		} else {
			totalCount += count
			totalRepl = append(totalRepl, repl...)
		}
	}
	if errCnt == slavesLen {
		errString = "All server err: " + errString
	} else if errCnt > 0 {
		errString = "Some server err: " + errString
	}
	sort.Sort(totalRepl)
	log.Infof("trackid:%v, repl:%v, err:%s", ctx.Value("trackid"), totalRepl, errString)
	return totalRepl, totalCount, errors.New(errString)
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
	log.Infof("trackid:%v, url:%s, body:%s, err:%s", ctx.Value("trackid"), url, tools.Bytes2Str(body), errString)
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
	log.Infof("trackid:%v, url:%s, err:%s", ctx.Value("trackid"), url, errString)
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
	log.Infof("trackid:%v, url:%s, err:%s", ctx.Value("trackid"), url, errString)
	return retByte, err
}
