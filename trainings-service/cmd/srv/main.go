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
	API := rest.NewAPI()
	API.SetEndpoints()

	cfg := server.DefaultHTTPConfig()
	cfg.Addr = "localhost:8090"

	srv := server.NewHTTP(API, cfg)
	srv.ListenAndServe()
	return nil
}
