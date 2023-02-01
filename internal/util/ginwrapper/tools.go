package ginwrapper

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Recovery() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				respData := make(map[string]interface{})
				respData["code"] = -1
				respData["message"] = fmt.Sprintf("%v", err)
				ctx.JSON(http.StatusOK, respData)
				errBuf := make([]byte, 0, 1024)
				errBuf = errBuf[:runtime.Stack(errBuf, false)]
				fmt.Printf("%s\n", string(errBuf))
				return
			}
		}()
		ctx.Next()
	}
}

func SetOpts(config Config) ([]Option, error) {
	var opts []Option
	if config.ReadTimeout != 0 {
		opts = append(opts, WithReadTimeout(config.ReadTimeout))
	}

	if config.WriteTimeout != 0 {
		opts = append(opts, WithWriteTimeout(config.WriteTimeout))
	}

	if config.IdleTimeout != 0 {
		opts = append(opts, WithIdleTimeout(config.IdleTimeout))
	}
	opts = append(opts, WithTLSConfig(config.Tls.Enable, config.Tls.CertFile, config.Tls.KeyFile))
	return opts, nil
}
