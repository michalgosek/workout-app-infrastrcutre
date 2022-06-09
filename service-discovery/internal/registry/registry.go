package registry

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type ServiceHealthEndpoints map[string]string

type ServiceInstance struct {
	Component string
	Name      string
	IP        string
	Port      string
	Endpoint  string
	healthy   bool
}

func (s *ServiceInstance) SetHealth(v bool) {
	s.healthy = v
}

type ServiceCluster struct {
	Name      string
	Instances map[string]ServiceInstance
}

type Repository interface {
	Register(ss ...ServiceInstance) error
	QueryInstances(name string) ([]ServiceInstance, error)
	UpdateStatus(s ServiceInstance) error
	ListClusters() []ServiceCluster
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

type Config struct {
	HeartBeat time.Duration
}

type Service struct {
	cfg        Config
	stop       chan struct{}
	healthz    ServiceHealthEndpoints
	repository Repository
	http       HTTPClient
	logger     Logger
}

func (s *Service) verifyInstances(ss ...ServiceInstance) error {
	if len(ss) == 0 {
		return nil
	}
	for _, si := range ss {
		if si.Component == "" {
			return fmt.Errorf("specified empty component name: %w", ErrMissingData)
		}
		if si.Name == "" {
			return fmt.Errorf("specified empty service name: %w", ErrMissingData)
		}

		_, err := url.ParseRequestURI(si.Endpoint)
		if err != nil {
			return fmt.Errorf("%v: %w", err, ErrMalformedData)
		}
		if si.IP != "localhost" {
			ip := net.ParseIP(si.IP)
			if ip == nil {
				return fmt.Errorf("parsing IP addres failure for %s: %w", si.IP, ErrInvalidDataFormat)
			}
		}
		port, err := strconv.ParseUint(si.Port, 10, 32)
		if err != nil {
			return fmt.Errorf("parse uint for port %s str failed: %v", si.Port, err)
		}
		isValidPort := port < 1 || port > 65535
		if isValidPort {
			return fmt.Errorf("port number range must be in (1-65535) for %s: %w", si.Name, ErrMissingData)
		}
	}
	return nil
}

func (s *Service) QueryInstances(name string) ([]ServiceInstance, error) {
	return s.repository.QueryInstances(name)
}

func (s *Service) Register(ss ...ServiceInstance) error {
	err := s.verifyInstances(ss...)
	if err != nil {
		return fmt.Errorf("instance verification failed: %w", err)
	}
	err = s.repository.Register(ss...)
	if err != nil {
		return fmt.Errorf("instances registration failed: %w", ErrRepositoryFailure)
	}
	return nil
}

func (s *Service) StopHeartBeat() {
	close(s.stop)
}

func (s *Service) HeartBeat() {
	defer s.logger.Infof("HeartBeat thread stopped.")
	s.logger.Infof("HeartBeat thread started.")
	for {
		select {
		case <-time.After(s.cfg.HeartBeat):
			clusters := s.repository.ListClusters()
			s.ProcessClusters(clusters...)
		case <-s.stop:
			return
		}
	}
}

func (s *Service) ProcessClusters(clusters ...ServiceCluster) {
	for _, c := range clusters {
		for _, v := range c.Instances {
			addr := fmt.Sprintf("http://%s:%s%s", v.IP, v.Port, v.Endpoint)
			resp, err := s.http.Get(addr)
			if err != nil || resp.StatusCode != http.StatusOK {
				s.logger.Errorf("[HealthCheckError][Cluster: %s, Node: %s]: %v", c.Name, v.Name, err)
				continue
			}
			healthy := true
			if resp.StatusCode != http.StatusOK {
				healthy = false
			}
			v.SetHealth(healthy)

			err = s.repository.UpdateStatus(v)
			if err != nil {
				s.logger.Errorf("Service-registry update: %v", err)
				continue
			}
			s.logger.Infof("[HealthCheck][Cluster: %s, Node: %s]: %s - status update to: %v", c.Name, v.Name, healthy)
		}
	}
}

type RegistryServiceOption func(s *Service)

func WithHTTPClient(cli HTTPClient) RegistryServiceOption {
	return func(s *Service) {
		if cli != nil {
			s.http = cli
		}
	}
}

func WithConfig(c Config) RegistryServiceOption {
	return func(s *Service) {
		s.cfg = c
	}
}

func WithRepository(r Repository) RegistryServiceOption {
	return func(s *Service) {
		if r != nil {
			s.repository = r
		}
	}
}

func WithHealthz(h ServiceHealthEndpoints) RegistryServiceOption {
	return func(s *Service) {
		if h != nil {
			s.healthz = h
		}
	}
}

func WithLogger(l Logger) RegistryServiceOption {
	return func(s *Service) {
		if l != nil {
			s.logger = l
		}
	}
}

func NewService(opts ...RegistryServiceOption) *Service {
	s := Service{
		http:       http.DefaultClient,
		repository: NewCacheRepository(),
		logger:     logrus.StandardLogger(),
		cfg: Config{
			HeartBeat: 5 * time.Second,
		},
		stop: make(chan struct{}),
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

var ErrMissingData = errors.New("provided data not required")
var ErrRepositoryFailure = errors.New("service registry repository failure")
var ErrMalformedData = errors.New("malformed data")
var ErrInvalidDataFormat = errors.New("invalid data format")
