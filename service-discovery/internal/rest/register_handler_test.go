package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"service-discovery/internal/rest"
	"service-discovery/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldReturnHTTPStatusBadRequestForEmptyPayloadUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingRequestBodyMsg,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.ServiceRegistry{}

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

func TestShouldReturnHTTPStatusBadRequestForEmptyInstanceNameUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Component: "service1",
		Port:      "9090",
		IP:        "localhost",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingHostNameMsg,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.ServiceRegistry{}

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

func TestShouldReturnHTTPStatusBadRequestForEmptyIPUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Component: "service",
		Instance:  "node1",
		Port:      "9090",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingHostIPMsg,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.ServiceRegistry{}

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

func TestShouldReturnHTTPStatusBadRequestForEmptyComponentUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Component: "",
		IP:        "localhost",
		Port:      "8080",
		Instance:  "node1",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingComponentMsg,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.ServiceRegistry{}
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

func TestShouldReturnHTTPStatusBadRequestForEmptyPortUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Component: "service1",
		IP:        "localhost",
		Instance:  "node1",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.MissingHostPortMsg,
		Code:    http.StatusBadRequest,
	}

	registryService := mocks.ServiceRegistry{}
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

func TestShouldReturnHTTPStatusOKForSucessfulServiceInstanceRegistry(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		IP:        "localhost",
		Instance:  "node1",
		Component: "service",
		Port:      "8080",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()

	registryService := mocks.ServiceRegistry{}
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

func TestShouldReturnInternalServiceErrorStatusWhenServiceInstanceRegistryFailure(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		IP:        "localhost",
		Component: "service1",
		Instance:  "node1",
		Port:      "8080",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.InternalServiceErrMsg,
		Code:    http.StatusInternalServerError,
	}

	registryService := mocks.ServiceRegistry{}
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

func createHTTPrequestWithoutBody(method string) *http.Request {
	return httptest.NewRequest(method, "http://localhost:9090/", nil)
}

func createHTTPrequestWithBody(method string, v rest.ServiceRegistryRequest) *http.Request {
	bb, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(bb)
	request := httptest.NewRequest(method, "http://localhost:9090/", body)
	return request
}
