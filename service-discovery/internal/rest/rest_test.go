package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"service-discovery/internal/rest"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckShouldReturnHTTPStatusOKUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	recoder := httptest.NewRecorder()
	cfg := rest.RegisterHandlerConfig{}
	handlerOpt := rest.WithRegisterHandleConfig(cfg)
	handler := rest.NewRegisterHandler(handlerOpt)
	SUT := rest.NewAPI(handler)
	SUT.SetEndpoints()

	expectedResponse := rest.JSONResponse{
		Message: "OK",
		Code:    http.StatusOK,
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
