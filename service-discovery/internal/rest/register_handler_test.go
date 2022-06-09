package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/registry"
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/rest/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceRegistryHandlerShouldReturnHTTPStatus500WhenDecodingRequestBodyFailure(t *testing.T) {
	assert := assert.New(t)

	// given:
	request := createHTTPrequestWithoutBody(http.MethodPost)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.InternalServiceErr,
		Code:    http.StatusInternalServerError,
	}

	registryService := mocks.RegistryService{}
	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.AssertNotCalled(t, "ServiceRegistry", mock.Anything)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func TestServiceRegistryHandlerShouldReturnHTTPStatusBadRequestForEmptyRequestBodyUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingRequestBody,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.RegistryService{}

	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.AssertNotCalled(t, "Register", mock.Anything)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func TestServiceRegistryHandlerShouldReturnHTTPStatusBadRequestForEmptyInstanceNameUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Component:      "service1",
		HealthEndpoint: "localhost:8080/v1/health",
		Port:           "9090",
		IP:             "localhost",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingHostName,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.RegistryService{}

	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.AssertNotCalled(t, "Register", mock.Anything)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func TestServiceRegistryHandlerShouldReturnHTTPStatusBadRequestForEmptyHealthEndpointUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		HealthEndpoint: "",
		Component:      "service1",
		Instance:       "node1",
		IP:             "localhost",
		Port:           "9090",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingHostHealthEndpoint,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.RegistryService{}

	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.AssertNotCalled(t, "Register", mock.Anything)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func TestServiceRegistryHandlerShouldReturnHTTPStatusBadRequestForEmptyIPUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Component:      "service",
		Instance:       "node1",
		HealthEndpoint: "localhost:8080/v1/health",
		Port:           "9090",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingHostIP,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.RegistryService{}

	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.AssertNotCalled(t, "Register", mock.Anything)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func TestServiceRegistryHandlerShouldReturnHTTPStatusBadRequestForEmptyComponentUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Component:      "",
		IP:             "localhost",
		HealthEndpoint: "localhost/v1/health",
		Port:           "8080",
		Instance:       "node1",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingComponent,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.RegistryService{}
	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.AssertNotCalled(t, "Register", mock.Anything)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func TestServiceRegistryHandlerShouldReturnHTTPStatusBadRequestForEmptyPortUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Component:      "service1",
		IP:             "localhost",
		HealthEndpoint: "localhost:8080/v1/health",
		Instance:       "node1",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingHostPort,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.RegistryService{}
	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.AssertNotCalled(t, "Register", mock.Anything)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func TestServiceRegistryHandlerShouldReturnHTTPStatusOKForSucessfulRegistrationUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		IP:             "localhost",
		Instance:       "node1",
		Component:      "service",
		HealthEndpoint: "localhost:8080/v1/health",
		Port:           "8080",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()

	registryService := mocks.RegistryService{}
	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.EXPECT().Register(mock.Anything).Return(nil)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse.Code, http.StatusOK)
	registryService.AssertExpectations(t)
}

func TestServiceRegistryHandlerShouldReturnInternalServiceErrorStatusWhenRegistrationFailureUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		IP:             "localhost",
		Component:      "service1",
		Instance:       "node1",
		HealthEndpoint: "localhost:8080/v1/health",
		Port:           "8080",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.InternalServiceErr,
		Code:    http.StatusInternalServerError,
	}

	registryService := mocks.RegistryService{}
	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	registryService.EXPECT().Register(mock.Anything).Return(errors.New("service is down"))

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func TestQueryInstancesHandlerShouldReturnHTTPStatusBadRequestForMissingComponentQueryParam(t *testing.T) {
	assert := assert.New(t)

	// given:
	request := createHTTPrequestWithoutBody(http.MethodGet)
	recoder := httptest.NewRecorder()

	registryService := mocks.RegistryService{}
	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.QueryInstances)

	expectedResponse := rest.JSONResponse{
		Code:    http.StatusBadRequest,
		Message: rest.MissingComponentQueryParam,
	}

	registryService.AssertNotCalled(t, "QueryInstances")

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(expectedResponse, actualResponse)
	registryService.AssertExpectations(t)
}
func TestQueryInstancesHandlerShouldReturnHTTPStatusOKForSucessfulQueryUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	component := "service1"
	request := createHTTPrequestWithoutBody(http.MethodGet)
	params := request.URL.Query()
	params.Add(rest.ComponentQueryParam, component)
	request.URL.RawQuery = params.Encode()

	recoder := httptest.NewRecorder()

	registryService := mocks.RegistryService{}
	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.QueryInstances)

	instances := []registry.ServiceInstance{
		{
			Name: "node1",
			IP:   "localhost",
			Port: "8080",
		},
		{
			Name: "node2",
			IP:   "localhost",
			Port: "8090",
		},
	}
	expectedResponse := rest.QueryInstancesRespone{
		Code:      http.StatusOK,
		Name:      component,
		Instances: instances,
	}

	registryService.EXPECT().QueryInstances(mock.Anything).Return(instances, nil)

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJQueryInstanceResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(expectedResponse, actualResponse)
	registryService.AssertExpectations(t)
}

func TestQueryInstancesHandlerShouldReturnInternalServiceErrorStatusWhenQueryFailureUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	component := "service1"
	request := createHTTPrequestWithoutBody(http.MethodGet)
	params := request.URL.Query()
	params.Add(rest.ComponentQueryParam, component)
	request.URL.RawQuery = params.Encode()
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.InternalServiceErr,
		Code:    http.StatusInternalServerError,
	}

	registryService := mocks.RegistryService{}
	opts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerRegistryService(&registryService),
	}
	registerHandler := rest.NewRegisterHandler(opts...)
	SUT := http.HandlerFunc(registerHandler.QueryInstances)

	registryService.EXPECT().QueryInstances(mock.Anything).Return(nil, errors.New("service is down"))

	// when:
	SUT.ServeHTTP(recoder, request)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
	registryService.AssertExpectations(t)
}

func createHTTPrequestWithoutBody(method string) *http.Request {
	return httptest.NewRequest(method, "http://localhost:9090/", nil)
}

func createHTTPrequestWithBody(method string, v interface{}) *http.Request {
	bb, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(bb)
	request := httptest.NewRequest(method, "http://localhost:9090/", body)
	return request
}
