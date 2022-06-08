package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type ServiceRegistryRequest struct {
	Name string
	IP   string
	Port string
}

func (s *ServiceRegistryRequest) Decode(body io.ReadCloser) error {
	dec := json.NewDecoder(body)
	err := dec.Decode(&s)
	if err != nil {
		return fmt.Errorf("decode failed: %v", err)
	}
	return nil
}

type RequestFieldValueErr struct {
	err   error
	value string
}

func (s *ServiceRegistryRequest) Verify() error {
	var empty ServiceRegistryRequest
	if *s == empty {
		return ErrMissingRequestBody
	}

	m := map[string]RequestFieldValueErr{
		"IP":   {value: s.IP, err: ErrMissingHostIPValue},
		"name": {value: s.Name, err: ErrMissingHostName},
		"port": {value: s.Port, err: ErrMissingHostPortValue},
	}
	for _, v := range m {
		if v.value == "" {
			return v.err
		}
	}
	return nil
}

var (
	ErrMissingRequestBody   = errors.New("missing request body")
	ErrMissingHostIPValue   = errors.New("missing service instance IP value in the request body")
	ErrMissingHostName      = errors.New("missing service instance Name in the request body")
	ErrMissingHostPortValue = errors.New("missing service instance Port value in the request body")
)

const (
	ErrInternalServiceErrMsg = "Internal Service Error"
)

func (h *RegisterHandler) ServiceRegistry(w http.ResponseWriter, r *http.Request) {
	var payload ServiceRegistryRequest
	err := payload.Decode(r.Body)
	if err != nil {
		response(w, JSONResponse{Message: "Internal Failure", Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	err = payload.Verify()
	if err != nil {
		response(w, JSONResponse{Message: err.Error(), Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}

	instance := registry.NewServiceInstance(payload.Name, payload.IP, payload.Port)
	err = h.service.Register(instance)
	if err != nil {
		response(w, JSONResponse{Message: ErrInternalServiceErrMsg, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Instance of service %s registered successfully", payload.Name)
	response(w, JSONResponse{Message: msg, Code: http.StatusOK}, http.StatusOK)
}

func (h *RegisterHandler) QueryInstances(w http.ResponseWriter, r *http.Request) {
	response(w, JSONResponse{Message: "OK", Code: http.StatusOK}, http.StatusOK)
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
