package napnap

import (
	"net/http"
	"time"
)

type Config struct {
	Addr         string
	TLSCertFile  string
	TLSKeyFile   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	*http.Server
	Config *Config
}

func NewHttpEngine(addr string) *Server {
	s := &Server{
		Server: new(http.Server),
	}
	s.Addr = addr
	return s
}

func NewHttpEngineWithConfig(c *Config) *Server {
	s := &Server{
		Server: new(http.Server),
		Config: c,
	}
	s.Addr = c.Addr
	s.ReadTimeout = c.ReadTimeout
	s.WriteTimeout = c.WriteTimeout
	return s
}
