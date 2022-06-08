package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"service-discovery/internal/rest"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr         string
	ShutdownTime time.Duration
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Server struct {
	logger Logger
	http   *http.Server
	router http.Handler
	cfg    Config
}

func (s *Server) Config() Config {
	return s.cfg
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe() error {
	s.logger.Infof("Server starts listen on port addr: %s\n", s.http.Addr)
	err := s.http.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve failed %v", err)
	}
	return nil
}

func (s *Server) GracefulShutdown(quit <-chan os.Signal, done chan<- struct{}) {
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTime)
	defer cancel()
	s.logger.Info("Shutting down service")
	err := s.http.Shutdown(ctx)
	if err != nil {
		s.logger.Errorf("shutting down the service failed; err: %v", err)
	}
	done <- struct{}{}
}

type Option func(s *Server)

func WithConfig(c Config) Option {
	return func(s *Server) {
		if len(c.Addr) > 0 {
			s.cfg.Addr = c.Addr
		}
		if c.ShutdownTime > 0 {
			s.cfg.ShutdownTime = c.ShutdownTime
		}
	}
}

func NewServer(opts ...Option) *Server {
	c := Config{
		ShutdownTime: 5 * time.Second,
		Addr:         "localhost:8080",
	}

	l := logrus.New()
	restAPI := rest.New()

	s := Server{
		logger: l,
		cfg:    c,
		http: &http.Server{
			Addr:           c.Addr,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        restAPI,
		},
	}

	for _, o := range opts {
		o(&s)
	}

	return &s
}
