module github.com/michalgosek/workout-app-infrastrcutre/trainings-service

go 1.17

require github.com/michalgosek/workout-app-infrastrcutre/service-discovery-cli v0.0.0

require (
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)

replace github.com/michalgosek/workout-app-infrastrcutre/service-discovery-cli => ../service-discovery-cli
