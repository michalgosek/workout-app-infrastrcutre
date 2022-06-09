package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type ServiceRegistryResponse struct {
	Message string
	Code    string
}

type ServiceInstance struct {
	Instance       string
	Component      string
	IP             string
	Port           string
	HealthEndpoint string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
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

type ServiceRegistry struct {
	addr    string
	cli     HTTPClient
	logger  Logger
	timeout time.Duration
}

func (s *ServiceRegistry) Register(instance ServiceInstance) error {
	bb, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}
	body := bytes.NewReader(bb)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addr, body)
	if err != nil {
		return fmt.Errorf("creating http request failed: %v", err)
	}
	resp, err := s.cli.Do(req)
	if err != nil {
		return fmt.Errorf("http request do failed: %v", err)
	}

	dec := json.NewDecoder(resp.Body)
	var dst ServiceRegistryResponse
	err = dec.Decode(&dst)
	if err != nil {
		return fmt.Errorf("decoding response failed %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(dst.Message)
	}
	return nil
}

type Option func(s *ServiceRegistry)

func WithServiceRegistryEndpoint(addr string) Option {
	return func(s *ServiceRegistry) {
		s.addr = addr
	}
}

func WithHTTPClient(c HTTPClient) Option {
	return func(s *ServiceRegistry) {
		s.cli = c
	}
}

func WithTimeout(d time.Duration) Option {
	return func(s *ServiceRegistry) {
		s.timeout = d
	}
}

func WithLogger(l Logger) Option {
	return func(s *ServiceRegistry) {
		s.logger = l
	}
}

func NewServiceRegistry(opts ...Option) *ServiceRegistry {
	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	sr := ServiceRegistry{
		addr:    "http://localhost:8080/api/v1/services/register",
		cli:     http.DefaultClient,
		logger:  l,
		timeout: 10 * time.Second,
	}
	for _, o := range opts {
		o(&sr)
	}
	return &sr
}
