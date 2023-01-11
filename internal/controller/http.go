package controller

import (
	"search_proxy/internal/model/proxy"
	"search_proxy/internal/objs"
	"search_proxy/internal/util/ginwrapper"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func StartNet(config objs.ServerConfig, closeFunc func()) error {
	opts, err := ginwrapper.SetOpts(config)
	if err != nil {
		return err
	}
	return ginwrapper.GinServer(config.IP, config.Port, router(config), closeFunc, opts...)
}

func router(config objs.ServerConfig) *gin.Engine {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(ginwrapper.Recovery())
	apiGroup := r.Group("/api/v1")
	{
		apiGroup.POST("/add_doc", proxy.AddDoc)
		apiGroup.GET("/del_doc", proxy.DelDoc)
		apiGroup.GET("/doc_isdel", proxy.DocIsDel)
		apiGroup.POST("/retrieve", proxy.RetrieveDoc)
	}
	if config.Debug {
		pprof.Register(r)
	}
	return r
}

/*
func addDoc(ctx *gin.Context) {
	remoteIp := ctx.ClientIP()
	uri := ctx.Request.RequestURI
	bodyReader := ctx.Request.Body
	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)
	body := buf.Bytes()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	retByte, err := proxy.AddDoc(newCtx, remoteIp, uri, body)

	if err != nil {
		respData := make(map[string]interface{})
		respData["code"] = -1
		respData["message"] = err.Error()
		ctx.JSON(http.StatusOK, respData)
	} else {
		ctx.Data(http.StatusOK, "application/json", retByte)
	}
}

func delDoc(ctx *gin.Context) {
	uri := ctx.Request.RequestURI
	remoteIp := ctx.ClientIP()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	retByte, err := proxy.DelDoc(newCtx, remoteIp, uri)

	if err != nil {
		respData := make(map[string]interface{})
		respData["code"] = -1
		respData["message"] = err.Error()
		ctx.JSON(http.StatusOK, respData)
	} else {
		ctx.Data(http.StatusOK, "application/json", retByte)
	}
}

func docIsDel(ctx *gin.Context) {
	uri := ctx.Request.RequestURI
	remoteIp := ctx.ClientIP()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	retByte, err := proxy.DocIsDel(newCtx, remoteIp, uri)

	if err != nil {
		respData := make(map[string]interface{})
		respData["code"] = -1
		respData["message"] = err.Error()
		ctx.JSON(http.StatusOK, respData)
	} else {
		ctx.Data(http.StatusOK, "application/json", retByte)
	}
}

func retrieveDoc(ctx *gin.Context) {
	remoteIp := ctx.ClientIP()
	uri := ctx.Request.RequestURI
	bodyReader := ctx.Request.Body
	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)
	body := buf.Bytes()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)

	repl, count, errString := proxy.RetrieveDoc(newCtx, remoteIp, uri, body)
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

	if len(errString) != 0 {
		respData := make(map[string]interface{})
		respData["code"] = -1
		respData["message"] = errString
		respData["count"] = count
		respData["result"] = repl
		ctx.JSON(http.StatusOK, respData)
	} else {
		respData := make(map[string]interface{})
		respData["code"] = 0
		respData["message"] = "ok"
		respData["count"] = replLen
		respData["result"] = repl
		ctx.JSON(http.StatusOK, respData)
	}
}*/
