package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"service-discovery/internal/registry"
	"strings"

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
	QueryInstances(name string) ([]registry.ServiceInstance, error)
}

type RegisterHandler struct {
	cfg     RegisterHandlerConfig
	logger  Logger
	service RegistryService
}

type ServiceRegistryRequest struct {
	Component string
	Instance  string
	IP        string
	Port      string
}

func (s *ServiceRegistryRequest) Decode(body io.ReadCloser) error {
	dec := json.NewDecoder(body)
	err := dec.Decode(&s)
	if err != nil {
		return fmt.Errorf("decode failed: %v", err)
	}
	return nil
}

const (
	InternalServiceErrMsg = "Internal Service Error"
	MissingRequestBodyMsg = "Missing request body."
	MissingComponentMsg   = "Missing component name value in the request body."
	MissingHostIPMsg      = "Missing service instance IP value in the request body."
	MissingHostNameMsg    = "Missing service instance Name in the request body."
	MissingHostPortMsg    = "Missing service instance Port value in the request body."
)

type RequestFieldErrorMsg struct {
	field  string
	errMsg string
}

func (s *ServiceRegistryRequest) Verify() string {
	var empty ServiceRegistryRequest
	recv := ServiceRegistryRequest{
		Component: strings.TrimSpace(s.Component),
		Instance:  strings.TrimSpace(s.Component),
		IP:        strings.TrimSpace(s.Component),
		Port:      strings.TrimSpace(s.Component),
	}
	if recv == empty {
		return MissingRequestBodyMsg
	}
	m := map[string]RequestFieldErrorMsg{
		"component": {errMsg: MissingComponentMsg, field: recv.Component},
		"IP":        {errMsg: MissingHostIPMsg, field: recv.IP},
		"instance":  {errMsg: MissingHostNameMsg, field: recv.Instance},
		"port":      {errMsg: MissingHostPortMsg, field: recv.Port},
	}
	for _, v := range m {
		if v.field == "" {
			return v.errMsg
		}
	}
	return ""
}

func (h *RegisterHandler) ServiceRegistry(w http.ResponseWriter, r *http.Request) {
	var payload ServiceRegistryRequest
	err := payload.Decode(r.Body)
	if err != nil {
		response(w, JSONResponse{Message: InternalServiceErrMsg, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	errMsg := payload.Verify()
	if errMsg != "" {
		response(w, JSONResponse{Message: errMsg, Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}

	instance := registry.NewServiceInstance(payload.Component, payload.Instance, payload.IP, payload.Port)
	err = h.service.Register(instance)
	if err != nil {
		response(w, JSONResponse{Message: InternalServiceErrMsg, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("Cluster %s - Instance of service %s registered successfully", payload.Component, payload.Instance)
	response(w, JSONResponse{Message: msg, Code: http.StatusOK}, http.StatusOK)
}

type QueryInstancesRequest struct {
	Component string
}

type QueryInstancesRespone struct {
	Code      int
	Name      string
	Instances []registry.ServiceInstance
}

func (q *QueryInstancesRequest) Decode(body io.ReadCloser) error {
	dec := json.NewDecoder(body)
	err := dec.Decode(&q)
	if err != nil {
		return fmt.Errorf("decode failed: %v", err)
	}
	return nil
}

func (q *QueryInstancesRequest) Verify() string {
	recv := QueryInstancesRequest{
		Component: strings.TrimSpace(q.Component),
	}
	var empty QueryInstancesRequest
	if recv == empty {
		return MissingRequestBodyMsg
	}
	m := map[string]RequestFieldErrorMsg{
		"component": {errMsg: MissingComponentMsg, field: strings.TrimSpace(q.Component)},
	}
	for _, v := range m {
		if v.field == "" {
			return v.errMsg
		}
	}
	return ""
}

func (h *RegisterHandler) QueryInstances(w http.ResponseWriter, r *http.Request) {
	var payload QueryInstancesRequest
	err := payload.Decode(r.Body)
	if err != nil {
		response(w, JSONResponse{Message: InternalServiceErrMsg, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	errMsg := payload.Verify()
	if errMsg != "" {
		response(w, JSONResponse{Message: errMsg, Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}
	instances, err := h.service.QueryInstances(payload.Component)
	if err != nil {
		response(w, JSONResponse{Message: InternalServiceErrMsg, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	response(w, QueryInstancesRespone{Code: http.StatusOK, Name: payload.Component, Instances: instances}, http.StatusOK)
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
