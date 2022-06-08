package rest

import (
	"net/http"
	"service-discovery/internal/registry"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

const (
	serviceRegistry = "service-registry"
	clientRegistry  = "client-registry"
)

type RegisterHandlerConfig struct {
	Endpoints map[string]string
}

type RegistryService interface {
	Register(ss ...registry.ServiceInstance) error
}

type RegisterHandler struct {
	cfg     RegisterHandlerConfig
	logger  Logger
	service RegistryService
}

func (r *RegisterHandler) ServiceRegiststryEndpoint() string {
	return r.cfg.Endpoints[serviceRegistry]
}

func (r *RegisterHandler) ClientRegistryEndpoint() string {
	return r.cfg.Endpoints[clientRegistry]
}

type RegisterHandlerOption func(r *RegisterHandler)

func WithRegisterHandlerConfig(c RegisterHandlerConfig) RegisterHandlerOption {
	return func(r *RegisterHandler) {
		if c.Endpoints != nil {
			r.cfg.Endpoints = c.Endpoints
		}
	}
}

func WithRegisterHandlerRegistryService(s RegistryService) RegisterHandlerOption {
	return func(r *RegisterHandler) {
		if s != nil {
			r.service = s
		}
	}
}

func WithRegisterHandlerLogger(l Logger) RegisterHandlerOption {
	return func(r *RegisterHandler) {
		if l != nil {
			r.logger = l
		}
	}
}

func NewRegisterHandler(opts ...RegisterHandlerOption) *RegisterHandler {
	h := RegisterHandler{
		logger: logrus.New(),
		cfg: RegisterHandlerConfig{
			Endpoints: map[string]string{
				serviceRegistry: "/v1/api/services/register",
				clientRegistry:  "/v1/api/services/query",
			},
		},
	}

	for _, o := range opts {
		o(&h)
	}
	return &h
}

type ServicerRegisterPayload struct {
	Name string
	IP   string
	Port string
}

func (h *RegisterHandler) ServiceRegistry(w http.ResponseWriter, r *http.Request) {
	response(w, JSONResponse{Message: "OK", Code: http.StatusOK}, http.StatusOK)
}

func (h *RegisterHandler) QueryInstances(w http.ResponseWriter, r *http.Request) {
	response(w, JSONResponse{Message: "OK", Code: http.StatusOK}, http.StatusOK)
}
