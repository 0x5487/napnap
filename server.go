package napnap

import (
	"net/http"
	"time"
)

// Config ...
type Config struct {
	Addr          string
	Domain        string // abc123.com, abc456.com
	CertCachePath string
	TLSCertFile   string
	TLSKeyFile    string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

// Server ...
type Server struct {
	*http.Server
	Config *Config
}

// NewHTTPEngine ...
func NewHTTPEngine(addr string) *Server {
	s := &Server{
		Server: new(http.Server),
	}
	s.Addr = addr
	return s
}

// NewHTTPEngineWithConfig ...
func NewHTTPEngineWithConfig(c *Config) *Server {
	s := &Server{
		Server: new(http.Server),
		Config: c,
	}
	s.Addr = c.Addr
	s.ReadTimeout = c.ReadTimeout
	s.WriteTimeout = c.WriteTimeout
	return s
}
