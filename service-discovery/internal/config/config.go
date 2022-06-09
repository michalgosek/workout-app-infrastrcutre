package config

import (
	"errors"
	"time"

	"fmt"
	"io/fs"

	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/registry"
	"github.com/spf13/viper"
)

type ServerHTTP struct {
	Addr           string
	ShutdownTime   time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

func DefaultServerHTTP() ServerHTTP {
	c := ServerHTTP{
		Addr:           "localhost:8080",
		ShutdownTime:   5 * time.Second,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return c
}

type Config struct {
	Server   ServerHTTP
	Registry registry.Config
}

func New(path string) (*Config, error) {
	var (
		cfg    Config
		target *fs.PathError
	)
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil && errors.As(err, &target) {
		return nil, ErrConfigFileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("reading config failed: %v", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file %s failed: %v", path, err)
	}
	return &cfg, nil
}

var ErrConfigFileNotFound = errors.New("config file not exist or path to the file is invalid")
