package server_test

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/server"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldShutdownGracefullyShutdownIntegration(t *testing.T) {
	assertions := assert.New(t)
	// given:
	router := server.NewRouter()
	cfg := server.DefaultHTTPConfig("localhost:8080", "test-server")
	srv := server.NewHTTP(router, cfg)

	time.AfterFunc(2*time.Second, func() {
		srv.Terminate()
	})

	// when:
	srv.StartHTTPServer()

	// then:
	assertions.True(srv.ConnsClosed())
}
