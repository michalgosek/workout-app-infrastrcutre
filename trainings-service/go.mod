module github.com/michalgosek/workout-app-infrastrcutre/trainings-service

go 1.17

require github.com/michalgosek/workout-app-infrastrcutre/service-utility v0.0.0

require (
	github.com/go-chi/chi v1.5.4 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20220608164250-635b8c9b7f68 // indirect
)

replace github.com/michalgosek/workout-app-infrastrcutre/service-utility => ../service-utility
