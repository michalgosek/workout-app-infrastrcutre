module github.com/michalgosek/workout-app-infrastrcutre/api-gateway

go 1.17

require (
	github.com/auth0/go-jwt-middleware/v2 v2.0.1
	github.com/go-chi/chi v1.5.4
	github.com/michalgosek/workout-app-infrastrcutre/service-utility v0.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd // indirect
	golang.org/x/sys v0.0.0-20220608164250-635b8c9b7f68 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/michalgosek/workout-app-infrastrcutre/service-utility => ../service-utility
