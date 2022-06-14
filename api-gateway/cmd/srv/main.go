package main

import (
	"log"

	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	API := rest.NewRouter()
	serverCfg := server.DefaultHTTPConfig("localhost:8080", "api-gateway")
	srv := server.NewHTTP(API, serverCfg)
	srv.StartHTTPServer()
	return nil
}
