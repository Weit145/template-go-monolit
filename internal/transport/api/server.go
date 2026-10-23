package api

import (
	"net/http"
	"net/netip"
	"time"
)

type Settings struct {
	Addr              string
	TrustedProxies    []netip.Prefix
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func NewServer(handler http.Handler, settings Settings) *http.Server {
	return &http.Server{
		Addr:              settings.Addr,
		Handler:           handler,
		ReadHeaderTimeout: settings.ReadHeaderTimeout,
		ReadTimeout:       settings.ReadTimeout,
		WriteTimeout:      settings.WriteTimeout,
		IdleTimeout:       settings.IdleTimeout,
	}
}
