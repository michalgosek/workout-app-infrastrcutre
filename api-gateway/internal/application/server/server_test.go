package server_test

import (
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/stretchr/testify/assert"
)

func TestShouldShutdownGracefullyShutdownIntegration(t *testing.T) {
	assert := assert.New(t)
	// given:
	router := rest.NewRouter()
	cfg := server.DefaultHTTPConfig("localhost:8080", "test-server")
	srv := server.NewHTTP(router, cfg)

	time.AfterFunc(2*time.Second, func() {
		srv.Terminate()
	})

	// when:
	srv.StartHTTPServer()

	// then:
	assert.True(srv.ConnsClosed())
}
