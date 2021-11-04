package proxy

import (
	"context"
	"search_proxy/internal/model/objs"
	"search_proxy/internal/util/idgenerator"
	"search_proxy/internal/util/log"
	"search_proxy/internal/util/tools"
	"time"
)

var px *proxy

func NewPx(groupMasters []string, groupSlaves [][]string, groupTimeout int, routerMode string) {
	px = newProxy(groupMasters, groupSlaves, groupTimeout, routerMode)
	idgenerator.NewIdGenerator()
}

func RetrieveDoc(ctx context.Context, routerKey string, uri string, body []byte) (objs.RecallPostingList, int, string) {
	defer func(cost func() time.Duration) {
		log.Infof("trackid:%v, cost: %.3f ms", ctx.Value("trackid"), float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	repl, count, errString := px.retrieveDoc(ctx, routerKey, uri, body)
	return repl, count, errString
}

func AddDoc(ctx context.Context, routerKey string, uri string, body []byte) ([]byte, error) {
	defer func(cost func() time.Duration) {
		log.Infof("trackid:%v, cost: %.3f ms", ctx.Value("trackid"), float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	retByte, err := px.addDoc(ctx, routerKey, uri, body)
	return retByte, err
}

func DelDoc(ctx context.Context, routerKey string, uri string) ([]byte, error) {
	defer func(cost func() time.Duration) {
		log.Infof("trackid:%v, cost: %.3f ms", ctx.Value("trackid"), float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	retByte, err := px.delDoc(ctx, routerKey, uri)
	return retByte, err
}

func DocIsDel(ctx context.Context, routerKey string, uri string) ([]byte, error) {
	defer func(cost func() time.Duration) {
		log.Infof("trackid:%v, cost: %.3f ms", ctx.Value("trackid"), float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	retByte, err := px.docIsDel(ctx, routerKey, uri)
	return retByte, err
}
