package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr           string
	ShutdownTime   time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

func DefaultHTTPConfig() Config {
	cfg := Config{
		Addr:           "localhost:8080",
		ShutdownTime:   5 * time.Second,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return cfg
}

type HTTP struct {
	srv          http.Server
	shutdownTime time.Duration
}

func NewHTTP(h http.Handler, cfg Config) *HTTP {
	srv := HTTP{
		shutdownTime: cfg.ShutdownTime,
		srv: http.Server{
			Addr:           cfg.Addr,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.ReadTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			Handler:        h,
		},
	}
	return &srv
}

func (s *HTTP) ListenAndServe() error {
	logrus.Infof("Server starts listen on port addr: %s\n", s.srv.Addr)
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve failed %v", err)
	}
	return nil
}

func (s *HTTP) GracefulShutdown(quit <-chan os.Signal, done chan<- struct{}, errc chan<- error) {
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTime)
	defer cancel()
	logrus.Info("Shutting down service")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		errc <- err
	}
	done <- struct{}{}
}
