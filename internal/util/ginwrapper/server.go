package ginwrapper

import (
	"context"
	"net/http"
	"strconv"
)

type Server struct {
	http.Server
}

func NewServer(ip string, port int, r http.Handler, opts ...Option) *Server {
	opts = append(opts, WithAddr(ip+":"+strconv.Itoa(port)), WithHandler(r))
	withOptions(opts)

	return &Server{
		http.Server{
			Handler:      sc.handler,
			ReadTimeout:  sc.readTimeout,
			WriteTimeout: sc.writeTimeout,
			IdleTimeout:  sc.idleTimeout,
			Addr:         sc.addr,
		},
	}
}

func (s *Server) start() error {
	if sc.enableTLS {
		return s.ListenAndServeTLS(sc.certFile, sc.keyFile)
	}
	return s.ListenAndServe()
}

func (s *Server) stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}
