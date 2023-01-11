package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"search_proxy/internal/objs"
	"search_proxy/internal/util/ginwrapper"
	"search_proxy/internal/util/idgenerator"
	"search_proxy/internal/util/log"
	"search_proxy/internal/util/tools"
	"time"

	"github.com/gin-gonic/gin"
)

type proxyWrapper struct {
	*proxy
	ginwrapper.Base
}

var pxw *proxyWrapper

func NewProxyWrap(groupConfig objs.GroupConfig, routerConfig objs.RouterConfig) {
	pxw = new(proxyWrapper)
	pxw.proxy = newProxy(groupConfig, routerConfig)
}

func RetrieveDoc(ctx *gin.Context) {
	remoteIp := ctx.ClientIP()
	uri := ctx.Request.RequestURI
	bodyReader := ctx.Request.Body
	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)
	body := buf.Bytes()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	defer func(cost func() time.Duration) {
		log.Infof("trackid:%v, cost: %.3f ms", newCtx.Value("trackid"), float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	repl, count, err := pxw.retrieveDoc(newCtx, remoteIp, uri, body)
	replLen := len(repl)
	m := make(map[string]interface{})
	json.Unmarshal(body, &m)
	offset := int(m["Offset"].(float64))
	limit := int(m["Limit"].(float64))
	end := offset + limit
	if replLen < offset {
		repl = make(objs.RecallPostingList, 0)
	} else if replLen < end {
		repl = repl[offset:]
	} else if replLen >= end {
		repl = repl[offset:end]
	}

	if err != nil {
		pxw.ErrMsg(ctx, err, objs.RetreiveDocResp{Count: count, Result: repl})
	} else {
		pxw.SucMsg(ctx, objs.RetreiveDocResp{Count: replLen, Result: repl})
	}
}

func AddDoc(ctx *gin.Context) {
	remoteIp := ctx.ClientIP()
	uri := ctx.Request.RequestURI
	bodyReader := ctx.Request.Body
	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)
	body := buf.Bytes()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	defer func(cost func() time.Duration) {
		log.Infof("trackid:%v, cost: %.3f ms", newCtx.Value("trackid"), float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	retByte, err := pxw.addDoc(newCtx, remoteIp, uri, body)
	if err != nil {
		pxw.ErrMsg(ctx, err)
	} else {
		pxw.SucMsg(ctx, retByte)
	}
}

func DelDoc(ctx *gin.Context) {
	uri := ctx.Request.RequestURI
	remoteIp := ctx.ClientIP()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	defer func(cost func() time.Duration) {
		log.Infof("trackid:%v, cost: %.3f ms", newCtx.Value("trackid"), float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	retByte, err := pxw.delDoc(newCtx, remoteIp, uri)
	if err != nil {
		pxw.ErrMsg(ctx, err)
	} else {
		pxw.SucMsg(ctx, retByte)
	}
}

func DocIsDel(ctx *gin.Context) {
	uri := ctx.Request.RequestURI
	remoteIp := ctx.ClientIP()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	defer func(cost func() time.Duration) {
		log.Infof("trackid:%v, cost: %.3f ms", newCtx.Value("trackid"), float64(cost().Microseconds())/1000.0)
	}(tools.TimeCost())

	retByte, err := pxw.docIsDel(newCtx, remoteIp, uri)
	if err != nil {
		pxw.ErrMsg(ctx, err)
	} else {
		pxw.SucMsg(ctx, retByte)
	}
}
