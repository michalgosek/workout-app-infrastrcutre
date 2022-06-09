package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/registry"
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

type RegistryService interface {
	Register(ss ...registry.ServiceInstance) error
	QueryInstances(name string) ([]registry.ServiceInstance, error)
}

type RegisterHandler struct {
	logger    Logger
	service   RegistryService
	endpoints map[string]string
}

type ServiceRegistryRequest struct {
	HealthEndpoint string
	Component      string
	Instance       string
	IP             string
	Port           string
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
	clientRegistry  = "clientRegistry"
	serviceRegistry = "serviceRegistry"
)

const (
	InternalServiceErr         = "Internal Service Error"
	MissingRequestBody         = "Missing request body."
	MissingComponent           = "Missing component name value in the request body."
	MissingHostIP              = "Missing service instance IP value in the request body."
	MissingHostName            = "Missing service instance Name in the request body."
	MissingHostPort            = "Missing service instance Port value in the request body."
	MissingHostHealthEndpoint  = "Missing service instance Health Endpoint value in the request body."
	MissingComponentQueryParam = "Missing component param in the request query."
)

func (s *ServiceRegistryRequest) Verify() string {
	var empty ServiceRegistryRequest
	recv := ServiceRegistryRequest{
		Component:      strings.TrimSpace(s.Component),
		Instance:       strings.TrimSpace(s.Instance),
		IP:             strings.TrimSpace(s.IP),
		Port:           strings.TrimSpace(s.Port),
		HealthEndpoint: strings.TrimSpace(s.HealthEndpoint),
	}
	if recv == empty {
		return MissingRequestBody
	}
	switch {
	case recv.Component == "":
		return MissingComponent
	case recv.HealthEndpoint == "":
		return MissingHostHealthEndpoint
	case recv.IP == "":
		return MissingHostIP
	case recv.Port == "":
		return MissingHostPort
	case recv.Instance == "":
		return MissingHostName
	}
	return ""
}

func (h *RegisterHandler) ServiceRegistry(w http.ResponseWriter, r *http.Request) {
	var payload ServiceRegistryRequest
	err := payload.Decode(r.Body)
	if err != nil {
		response(w, JSONResponse{Message: InternalServiceErr, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	errMsg := payload.Verify()
	if errMsg != "" {
		response(w, JSONResponse{Message: errMsg, Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}

	instance := registry.ServiceInstance{
		Component: payload.Component,
		Name:      payload.Instance,
		IP:        payload.IP,
		Port:      payload.Port,
		Endpoint:  payload.HealthEndpoint,
	}
	err = h.service.Register(instance)
	if err != nil {
		response(w, JSONResponse{Message: InternalServiceErr, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("Cluster %s - Instance of service %s registered successfully", payload.Component, payload.Instance)
	response(w, JSONResponse{Message: msg, Code: http.StatusOK}, http.StatusOK)
}

type QueryInstancesRespone struct {
	Code      int
	Name      string
	Instances []registry.ServiceInstance
}

const ComponentQueryParam = "component"

func (h *RegisterHandler) QueryInstances(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	componentName := query.Get(ComponentQueryParam)
	if componentName == "" {
		response(w, JSONResponse{Message: MissingComponentQueryParam, Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}
	instances, err := h.service.QueryInstances(componentName)
	if err != nil {
		response(w, JSONResponse{Message: InternalServiceErr, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	response(w, QueryInstancesRespone{Code: http.StatusOK, Name: componentName, Instances: instances}, http.StatusOK)
}

func (r *RegisterHandler) ServiceRegiststryEndpoint() string {
	return r.endpoints[serviceRegistry]
}

func (r *RegisterHandler) ClientRegistryEndpoint() string {
	return r.endpoints[clientRegistry]
}

type RegisterHandlerOption func(r *RegisterHandler)

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
		endpoints: map[string]string{
			serviceRegistry: "/api/v1/services/register",
			clientRegistry:  "/api/v1/services/query",
		},
	}

	for _, o := range opts {
		o(&h)
	}
	return &h
}
