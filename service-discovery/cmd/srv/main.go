package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service-discovery/internal/config"
	"service-discovery/internal/rest"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	var path string
	flag.StringVar(&path, "c", "./config.yml", "service config file path")
	cfg, err := config.New(path)
	if err != nil {
		return fmt.Errorf("config file read failed %v", err)
	}
	serverCfg := createHTTPServerCfg(cfg.Server)
	runHTTPServer(serverCfg)
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

func runHTTPServer(cfg config.ServerHTTP) {
	logger := logrus.New()
	restAPI := rest.New()
	srv := http.Server{
		Addr:           cfg.Addr,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.ReadTimeout,
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
	logger.Infof("Server starts listen on port addr: %s\n", cfg.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	logger.Info("service has been gracefully shutdown")
}
