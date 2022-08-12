package server

import (
	"encoding/json"
	"github.com/go-chi/cors"
	"net/http"
	"notification-service/internal/application/authorization"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	HealthEndpoint = "/api/v1/health"
)

// NewRouter returns chi.Router with basic middlewares setup, health check route.
func NewRouter() chi.Router {
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(authorization.ValidateJWT())
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get(HealthEndpoint, healthHandler)
	return r
}

type JSONResponse struct {
	Message string `json:"message"`
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	SendJSONResponse(w, JSONResponse{Message: "OK"}, http.StatusOK)
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
