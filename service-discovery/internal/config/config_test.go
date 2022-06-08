package config_test

import (
	"service-discovery/internal/config"

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
	}

	// when:
	actualCfg, err := config.New(path)

	// then:
	assert.Nil(err)
	assert.Equal(expectedConfig, actualCfg)
}
