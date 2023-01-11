package ginwrapper

import "time"

type Config struct {
	Name         string
	IP           string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Debug        bool
	Tls          TLS
}

type TLS struct {
	Enable   bool
	CertFile string
	KeyFile  string
}
