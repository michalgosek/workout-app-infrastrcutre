package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/config"
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/registry"
	"github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/rest"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	var Path string
	flag.StringVar(&Path, "c", "./config.yml", "service config file path")
	flag.Parse()
	cfg, err := config.New(Path)
	if err != nil {
		return fmt.Errorf("config file read failed %v", err)
	}
	serverCfg := createHTTPServerCfg(cfg.Server)
	logger := logrus.New()

	repository := registry.NewCacheRepository()
	registryServiceOpts := []registry.RegistryServiceOption{
		registry.WithRepository(repository),
		registry.WithLogger(logger),
	}
	registryService := registry.NewService(registryServiceOpts...)

	registerHandlerOpts := []rest.RegisterHandlerOption{
		rest.WithRegisterHandlerLogger(logger),
		rest.WithRegisterHandlerRegistryService(registryService),
	}
	registerHandler := rest.NewRegisterHandler(registerHandlerOpts...)

	restAPI := rest.NewAPI(registerHandler)
	restAPI.SetEndpoints()

	srv := http.Server{
		Addr:           serverCfg.Addr,
		MaxHeaderBytes: serverCfg.MaxHeaderBytes,
		ReadTimeout:    serverCfg.ReadTimeout,
		WriteTimeout:   serverCfg.ReadTimeout,
		Handler:        restAPI,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		logger.Info("Shutting down service")
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	logger.Infof("Server starts listen on port addr: %s\n", serverCfg.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	logger.Info("service has been gracefully shutdown")
	return nil
}

func createHTTPServerCfg(cfg config.ServerHTTP) config.ServerHTTP {
	defaultCfg := config.DefaultServerHTTP()
	switch {
	case cfg.Addr != "localhost:8080":
		defaultCfg.Addr = cfg.Addr
	case cfg.MaxHeaderBytes > 0:
		defaultCfg.MaxHeaderBytes = cfg.MaxHeaderBytes
	case cfg.ShutdownTime > 0:
		defaultCfg.ShutdownTime = cfg.ShutdownTime
	case cfg.ReadTimeout > 0:
		defaultCfg.ReadTimeout = cfg.ReadTimeout
	case cfg.WriteTimeout > 0:
		defaultCfg.WriteTimeout = cfg.WriteTimeout
	}
	return defaultCfg
}
