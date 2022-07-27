package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func DefaultHTTPConfig(addr, name string) Config {
	cfg := Config{
		Addr:           addr,
		Name:           name,
		ShutdownTime:   5 * time.Second,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return cfg
}

type Config struct {
	Addr           string
	Name           string
	ShutdownTime   time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type HTTP struct {
	srv            http.Server
	cfg            Config
	sigs           chan os.Signal
	ideConnsClosed chan struct{}
}

func (h *HTTP) ConnsClosed() bool {
	return len(h.ideConnsClosed) == 0
}

func (h *HTTP) Terminate() {
	h.sigs <- syscall.SIGTERM
}

func (h *HTTP) StartHTTPServer() {
	ctx := context.Background()

	go func() {
		signal.Notify(h.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-h.sigs

		ctx, cancel := context.WithTimeout(ctx, h.cfg.ShutdownTime)
		defer cancel()
		err := h.srv.Shutdown(ctx)
		if err != nil {
			logrus.WithError(err).Panic("shutting down the service failed")
		}
		close(h.ideConnsClosed)
	}()
	logrus.WithField("HTTP-SERVER", h.cfg.Name).Info("started listening on addr: ", h.cfg.Addr)
	err := h.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logrus.WithError(err).Panic("cannot start listening")
	}
	<-h.ideConnsClosed

	logrus.WithField("HTTP-SERVER", h.cfg.Name).Info("stopped listening and done gracefully shutdown")
}

func NewHTTP(h http.Handler, cfg Config) *HTTP {
	http := HTTP{
		srv: http.Server{
			Addr:           cfg.Addr,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.ReadTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			Handler:        h,
		},
		cfg:            cfg,
		sigs:           make(chan os.Signal, 1),
		ideConnsClosed: make(chan struct{}),
	}
	return &http
}
