package rest_test

import (
	"testing"
)

func TestServiceRegisterHandlerShouldReturnHTTPStatusBadRequestForEmptyPayload(t *testing.T) {
	// assert := assert.New(t)

	// // given:
	// recoder := httptest.NewRecorder()
	// dummyEndpoint := "http://localhost:9090/"
	// request := httptest.NewRequest(http.MethodGet, dummyEndpoint, nil)
	// registryService := mocks.ServiceRegistry{}
	// opts := []rest.RegisterHandlerOption{
	// 	rest.WithRegisterHandlerRegistryService(&registryService),
	// }

	// registerHandler := rest.NewRegisterHandler(opts...)

	// SUT := http.HandlerFunc(registerHandler.ServiceRegistry)

	// expectedResponse := rest.JSONResponse{
	// 	Message: "specifeid empty payload",
	// 	Code:    http.StatusOK,
	// }

	// // when:
	// SUT.ServeHTTP(recoder, request)

	// // then:
	// actualResponse, err := convertToJSONResponse(recoder.Body)
	// assert.Nil(err)
	// assert.Equal(actualResponse, expectedResponse)
	// mock.AssertExpectationsForObjects(t, registryService)

}
