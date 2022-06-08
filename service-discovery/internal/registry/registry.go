package registry

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ServiceHealthEndpoints map[string]string

type ServiceInstance struct {
	name    string
	ip      string
	port    uint
	healthy bool
}

type ServiceCluster struct {
	Name      string
	Instances []ServiceInstance
}

type ServiceRegistry interface {
	Register(ss ...ServiceInstance) error
	QueryInstances(name string) ([]ServiceInstance, error)
	UpdateStatus(s ServiceInstance) error
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

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Service struct {
	healthz  ServiceHealthEndpoints
	registry ServiceRegistry
	http     HTTPClient
	logger   Logger
}

func (s *Service) verifyInstances(ss ...ServiceInstance) error {
	if len(ss) == 0 {
		return ErrEmptyServiceInstances
	}
	for _, si := range ss {
		if si.name == "" {
			return fmt.Errorf("specified empty service name: %w", ErrMalformedData)
		}
		if si.ip != "localhost" {
			ip := net.ParseIP(si.ip)
			if ip == nil {
				return fmt.Errorf("parsing IP addres failure for %s: %w", si.name, ErrMalformedData)
			}
		}
		if si.port == 0 || si.port > 65535 {
			return fmt.Errorf("port number range must be in (1-65535) for %s: %w", si.name, ErrMalformedData)
		}
	}
	return nil
}

func (s *Service) QueryInstances(name string) ([]ServiceInstance, error) {
	return s.registry.QueryInstances(name)
}

func (s *Service) Register(ss ...ServiceInstance) error {
	err := s.verifyInstances(ss...)
	if err != nil {
		return fmt.Errorf("instance verification failed: %w", err)
	}
	err = s.registry.Register(ss...)
	if err != nil {
		return fmt.Errorf("instances registration failed: %w", ErrRepositoryFailure)
	}
	return nil
}

func (s *Service) HeartBeat() {
	for service, addr := range s.healthz {
		instances, err := s.registry.QueryInstances(service)
		if err != nil {
			s.logger.Errorf("Service-registry query instance: %v", err)
			continue
		}
		err = s.updateClusterInstancesStatus(addr, instances...)
		if err != nil {
			s.logger.Errorf("Service-registry update: %v", err)
			continue
		}
	}
}

func (s *Service) updateClusterInstancesStatus(addr string, instances ...ServiceInstance) error {
	for _, ins := range instances {
		resp, err := s.http.Get(addr)
		if err != nil {
			s.logger.Errorf("HTTP-CLI: %v", err)
			continue
		}
		healthy := true
		if resp.StatusCode != http.StatusOK {
			healthy = false
		}
		ins.SetHealth(healthy)

		err = s.registry.UpdateStatus(ins)
		if err != nil {
			s.logger.Errorf("Service-registry update: %v", err)
			continue
		}
	}
	return nil
}

type Option func(s *Service)

func WithHTTPClient(cli HTTPClient) Option {
	return func(s *Service) {
		if cli != nil {
			s.http = cli
		}
	}
}

func WithRegistry(r ServiceRegistry) Option {
	return func(s *Service) {
		if r != nil {
			s.registry = r
		}
	}
}

func WithHealthz(h ServiceHealthEndpoints) Option {
	return func(s *Service) {
		if h != nil {
			s.healthz = h
		}
	}
}

func WithLogger(l Logger) Option {
	return func(s *Service) {
		if l != nil {
			s.logger = l
		}
	}
}

func New(opts ...Option) *Service {
	s := Service{
		registry: NewCacheServiceRegistry(),
		logger:   logrus.StandardLogger(),
		healthz: ServiceHealthEndpoints{
			"users-service":     "http://localhost:8030/api/v1/health",
			"trainer-service":   "http://localhost:8040/api/v1/health",
			"trainings-service": "http://localhost:8050/api/v1/health",
		},
	}
	for _, o := range opts {
		o(&s)
	}
	return &s
}

var ErrEmptyServiceInstances = errors.New("empty service instances")
var ErrMalformedData = errors.New("malformed data")
var ErrRepositoryFailure = errors.New("service registry repository failure")
