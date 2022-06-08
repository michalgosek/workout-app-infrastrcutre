package http_test

import (
	"os"
	"service-discovery/internal/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_GracefulShutdown_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	interruptSigFunc := func(quit chan<- os.Signal) {
		time.Sleep(1 * time.Second)
		quit <- os.Interrupt
	}

	SUT := http.NewServer()
	quit := make(chan os.Signal)
	done := make(chan struct{})
	go SUT.GracefulShutdown(quit, done)
	go interruptSigFunc(quit)

	// when:
	err := SUT.ListenAndServe()

	// then:
	assert.Nil(err)
	assert.Empty(quit, done)
}
