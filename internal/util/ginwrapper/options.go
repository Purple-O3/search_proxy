package ginwrapper

import (
	"net/http"
	"time"
)

type serverConfig struct {
	addr         string
	handler      http.Handler
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
	certFile     string
	keyFile      string
	enableTLS    bool
}

type Option func()

var sc serverConfig

func init() {
	sc = serverConfig{
		readTimeout:  500 * time.Millisecond,
		writeTimeout: 500 * time.Millisecond,
		idleTimeout:  500 * time.Millisecond,
		handler:      http.DefaultServeMux,
	}
}

func withOptions(opts []Option) {
	for _, opt := range opts {
		opt()
	}
}

func WithAddr(addr string) Option {
	return func() {
		sc.addr = addr
	}
}

func WithHandler(handler http.Handler) Option {
	return func() {
		sc.handler = handler
	}
}

func WithTLSConfig(enable bool, cert string, key string) Option {
	return func() {
		sc.enableTLS = enable
		sc.certFile = cert
		sc.keyFile = key
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func() {
		sc.readTimeout = timeout * time.Millisecond
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func() {
		sc.writeTimeout = timeout * time.Millisecond
	}
}

func WithIdleTimeout(timeout time.Duration) Option {
	return func() {
		sc.idleTimeout = timeout * time.Millisecond
	}
}
