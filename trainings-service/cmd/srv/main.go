package main

import (
	"log"

	cli "github.com/michalgosek/workout-app-infrastrcutre/service-discovery-cli"
)

func main() {
	instance := cli.ServiceInstance{
		Component:      "service1",
		Instance:       "dd",
		IP:             "localhost",
		Port:           "9060",
		HealthEndpoint: "/v1/api/health/",
	}

	cli := cli.NewServiceRegistry()

	err := cli.Register(instance)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")
}
