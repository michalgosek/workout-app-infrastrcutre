package client_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	client "github.com/michalgosek/workout-app-infrastrcutre/service-discovery-cli"
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery-cli/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldRegisterServiceInstanceWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	HTTPCli := mocks.HTTPClient{}
	logger := mocks.Logger{}
	opts := []client.Option{
		client.WithHTTPClient(&HTTPCli),
		client.WithLogger(&logger),
	}
	instance := client.ServiceInstance{
		Component:      "service1",
		Instance:       "node1",
		IP:             "localhost",
		Port:           "8080",
		HealthEndpoint: "/api/v1/health",
	}
	body := client.ServiceRegistryResponse{
		Message: "OK",
		Code:    "OK",
	}

	res := newHTTPResponse(t, body)
	HTTPCli.EXPECT().Do(mock.Anything).Return(res, nil)

	SUT := client.NewServiceRegistry(opts...)

	// when:
	err := SUT.Register(instance)

	// then:
	assert.Nil(err)
	HTTPCli.AssertExpectations(t)
}

func TestShouldReturnErrorWhenRegisterServiceInstanceUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	HTTPCli := mocks.HTTPClient{}
	logger := mocks.Logger{}
	opts := []client.Option{
		client.WithHTTPClient(&HTTPCli),
		client.WithLogger(&logger),
	}
	instance := client.ServiceInstance{
		Component:      "service1",
		Instance:       "node1",
		IP:             "localhost",
		Port:           "8080",
		HealthEndpoint: "/api/v1/health",
	}
	body := client.ServiceRegistryResponse{
		Message: "Something's gone wrong",
		Code:    "500",
	}
	expectedErr := errors.New(body.Message)

	res := newHTTPResponse(t, body)
	res.StatusCode = http.StatusInternalServerError
	HTTPCli.EXPECT().Do(mock.Anything).Return(res, nil)

	SUT := client.NewServiceRegistry(opts...)

	// when:
	err := SUT.Register(instance)

	// then:
	assert.Equal(err, expectedErr)
	HTTPCli.AssertExpectations(t)
}

func newHTTPResponse(t *testing.T, body client.ServiceRegistryResponse) *http.Response {
	bb, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal failed: %s", err)
	}
	res := http.Response{
		Status:     "200",
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(bb)),
	}
	return &res
}
