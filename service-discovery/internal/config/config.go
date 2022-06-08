package config

import (
	"errors"
	"service-discovery/internal/http"

	"fmt"
	"io/fs"

	"github.com/spf13/viper"
)

type Config struct {
	Server http.Config
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
