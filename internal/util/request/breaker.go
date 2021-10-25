package request

import (
	"context"
	"errors"
	"search_proxy/internal/util/log"
	"strings"
	"sync"
	"time"
)

//conf
const retryInterval = 500
const retryCnt = 50
const failedThreshold = 20

type state int

const (
	stateClose state = iota
	stateOpen
	stateHalfOpen
)

type breaker struct {
	lock            sync.Mutex
	state           state
	totalReqCnt     int64
	totalSuccessCnt int64
	totalFailedCnt  int64
	startTimer      bool
	timer           *time.Timer
	t               int64
}

var bk *breaker

func NewBreaker() {
	bk = new(breaker)
	bk.state = stateClose
}

func DoCall(reqName string, args ...interface{}) ([]byte, error) {
	var respByte []byte
	getState()
	//Tips:防止beforeCall和afterCall两次compare&set加锁造成状态前后不一致
	st, err := bk.beforeCall()
	if err != nil {
		return nil, err
	}
	reqName = strings.ToLower(reqName)
	if reqName == "get" {
		ctx := args[0].(context.Context)
		url := args[1].(string)
		timeout := args[2].(int)
		respByte, err = Get(ctx, url, timeout)
	} else if reqName == "post" {
		ctx := args[0].(context.Context)
		url := args[1].(string)
		contentType := args[2].(string)
		body := args[3].([]byte)
		timeout := args[4].(int)
		respByte, err = Post(ctx, url, contentType, body, timeout)
	}
	bk.afterCall(st, err)
	return respByte, err
}

func getState() state {
	switch bk.state {
	case stateClose:
		log.Debugf("breaker state: stateClose")
	case stateOpen:
		log.Debugf("breaker state: stateOpen")
	case stateHalfOpen:
		log.Debugf("breaker state: stateHalfOpen")
	}
	return bk.state
}

func (bk *breaker) beforeCall() (state, error) {
	bk.lock.Lock()
	defer bk.lock.Unlock()
	switch bk.state {
	case stateClose:
		return bk.closedBefore()
	case stateOpen:
		return bk.openBefore()
	case stateHalfOpen:
		return bk.halfOpenBefore()
	default:
		return bk.state, errors.New("break state error: default")
	}
}

func (bk *breaker) afterCall(st state, err error) {
	bk.lock.Lock()
	defer bk.lock.Unlock()
	if st != bk.state {
		return
	}
	switch bk.state {
	case stateClose:
		bk.closedAfter(err)
	case stateOpen:
	case stateHalfOpen:
		bk.halfOpenAfter(err)
	default:
	}
}

func (bk *breaker) closedBefore() (state, error) {
	return bk.state, nil
}

func (bk *breaker) closedAfter(err error) {
	now := time.Now().Unix() / 60
	if bk.t != now {
		bk.totalReqCnt = 0
		bk.totalSuccessCnt = 0
		bk.totalFailedCnt = 0
		bk.t = now
	}

	if err != nil {
		bk.totalFailedCnt += 1
	}
	bk.totalReqCnt += 1
	if bk.totalFailedCnt*100/bk.totalReqCnt > failedThreshold {
		bk.totalReqCnt = 0
		bk.totalSuccessCnt = 0
		bk.totalFailedCnt = 0
		bk.state = stateOpen
	}
}

func (bk *breaker) openBefore() (state, error) {
	if !bk.startTimer {
		bk.timer = time.NewTimer(retryInterval * time.Millisecond)
		bk.startTimer = true
	}
	select {
	case <-bk.timer.C:
		bk.startTimer = false
		bk.state = stateHalfOpen
		return bk.state, errors.New("breaker state open: refuse req")
	default:
		return bk.state, errors.New("breaker state open: refuse req")
	}
}

func (bk *breaker) halfOpenBefore() (state, error) {
	if bk.totalReqCnt >= retryCnt {
		return bk.state, errors.New("breaker state breaker state half_open: too many req")
	}
	bk.totalReqCnt += 1
	return bk.state, nil
}

func (bk *breaker) halfOpenAfter(err error) {
	if err != nil {
		bk.totalReqCnt = 0
		bk.totalSuccessCnt = 0
		bk.totalFailedCnt = 0
		bk.state = stateOpen
		return
	}

	bk.totalSuccessCnt += 1
	if bk.totalSuccessCnt >= retryCnt {
		bk.totalReqCnt = 0
		bk.totalSuccessCnt = 0
		bk.totalFailedCnt = 0
		bk.state = stateClose
	}
}
