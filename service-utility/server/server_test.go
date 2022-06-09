package server_test

import (
	"os"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/stretchr/testify/assert"
)

func TestShouldShutdownGracefullyShutdownIntegration(t *testing.T) {
	assert := assert.New(t)

	interruptSigFunc := func(quit chan<- os.Signal) {
		time.Sleep(1 * time.Second)
		quit <- os.Interrupt
	}
	API := rest.NewAPI()
	s := server.NewHTTP(API, server.DefaultHTTPConfig())
	quit := make(chan os.Signal)
	done := make(chan struct{})
	errc := make(chan error)

	go s.GracefulShutdown(quit, done, errc)
	go interruptSigFunc(quit)

	err := s.ListenAndServe()
	assert.Nil(err)
	assert.Empty(errc)
	assert.Empty(quit)
	assert.Empty(done)
}
