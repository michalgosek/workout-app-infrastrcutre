package config_test

import (
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/config"
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/rest"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrorAfterInvalidConfigPathUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const path = "./random_file.yml"

	// when:
	actualCfg, err := config.New(path)

	// then:
	assert.ErrorIs(err, config.ErrConfigFileNotFound)
	assert.Nil(actualCfg)
}

func TestShouldReadConfigFileWithoutError(t *testing.T) {
	assert := assert.New(t)

	// given:
	const path = "./example_cfg.yml"
	expectedConfig := &config.Config{
		Server: config.ServerHTTP{
			Addr:         "localhost:8090",
			ShutdownTime: 10 * time.Second,
		},
		RegisterHandler: rest.RegisterHandlerConfig{
			Endpoints: map[string]string{
				"client-registry":  "/v1/api/query",
				"service-registry": "/v1/api/registry",
			},
		},
	}

	// when:
	actualCfg, err := config.New(path)

	// then:
	assert.Nil(err)
	assert.Equal(expectedConfig, actualCfg)
}
