package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	client "github.com/michalgosek/workout-app-infrastrcutre/service-discovery-cli"
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/rest"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckShouldReturnHTTPStatusOKUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	recoder := httptest.NewRecorder()
	handler := rest.NewRegisterHandler()
	SUT := rest.NewAPI(handler)
	SUT.SetEndpoints()

	expectedResponse := client.ServiceRegistryResponse{
		Message: "OK",
		Code:    "OK",
	}

	req := httptest.NewRequest(http.MethodGet, rest.HealthEndpoint, nil)

	// when:
	SUT.ServeHTTP(recoder, req)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assert.Nil(err)
	assert.Equal(actualResponse, expectedResponse)
}

func convertToJSONResponse(body *bytes.Buffer) (rest.JSONResponse, error) {
	var res rest.JSONResponse
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		return rest.JSONResponse{}, fmt.Errorf("decode failed: %v", err)
	}
	return res, nil
}

func convertToJQueryInstanceResponse(body *bytes.Buffer) (rest.QueryInstancesRespone, error) {
	var res rest.QueryInstancesRespone
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		return rest.QueryInstancesRespone{}, fmt.Errorf("decode failed: %v", err)
	}
	return res, nil
}
