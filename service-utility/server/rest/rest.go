package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	HealthEndpoint  = "/api/v1/health"
	VersionEndponit = "/api/v1/version"
)

// NewRouter returns chi.Router with basic middlewares setup, health check route.
func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get(HealthEndpoint, healthHandler)
	return r
}

type JSONResponse struct {
	Message string
	Code    int
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	SendJSONResponse(w, JSONResponse{Message: "OK", Code: http.StatusOK}, http.StatusOK)
}

func SendJSONResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")
		err := enc.Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
