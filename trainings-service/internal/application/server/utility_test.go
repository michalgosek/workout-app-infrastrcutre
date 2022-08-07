package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckShouldReturnHTTPStatusOKUnit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	recoder := httptest.NewRecorder()
	SUT := server.NewRouter()

	expectedResponse := server.JSONResponse{
		Message: "OK",
	}

	req := httptest.NewRequest(http.MethodGet, server.HealthEndpoint, nil)

	// when:
	SUT.ServeHTTP(recoder, req)

	// then:
	actualResponse, err := convertToJSONResponse(recoder.Body)
	assertions.Nil(err)
	assertions.Equal(actualResponse, expectedResponse)
}

func convertToJSONResponse(body *bytes.Buffer) (server.JSONResponse, error) {
	var res server.JSONResponse
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		return server.JSONResponse{}, fmt.Errorf("decode failed: %v", err)
	}
	return res, nil
}
