package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"service-discovery/internal/rest"
	"service-discovery/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldReturnHTTPStatusBadRequestForEmptyPayload(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.ErrMissingRequestBody.Error(),
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

func TestShouldReturnHTTPStatusBadRequestForEmptyNameValue(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Port: "9090",
		IP:   "localhost",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.ErrMissingHostName.Error(),
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

func TestShouldReturnHTTPStatusBadRequestForEmptyIPValue(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		Name: "dummy",
		Port: "9090",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.ErrMissingHostIPValue.Error(),
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

func TestShouldReturnHTTPStatusBadRequestForPortValue(t *testing.T) {
	assert := assert.New(t)

	// given:
	requestBody := rest.ServiceRegistryRequest{
		IP:   "localhost",
		Name: "dummy",
	}
	request := createHTTPrequestWithBody(http.MethodPost, requestBody)
	recoder := httptest.NewRecorder()
	expectedResponse := rest.JSONResponse{
		Message: rest.ErrMissingHostPortValue.Error(),
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
