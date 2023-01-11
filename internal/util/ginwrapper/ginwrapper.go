package ginwrapper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// GinServer gin server
func GinServer(ip string, port int, r *gin.Engine, closeFunc func(), opts ...Option) (err error) {
	if r == nil {
		return errors.New("router is nil")
	}

	// new http server
	httpServer := NewServer(ip, port, r, opts...)

	// register quit signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)

	// http server start
	go func() {
		fmt.Printf("server start!!!")
		if err := httpServer.start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("%v", err)
			quit <- syscall.SIGQUIT
		}
	}()

	// wait quit
	<-quit
	fmt.Printf("server stop!!!")

	// exec close func
	closeFunc()

	// http server stop
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpServer.stop(stopCtx)
}
