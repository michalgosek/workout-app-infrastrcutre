package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ServiceRegisteryRequest struct {
	Name string
	IP   string
	Port string
}

const (
	HealthEndpoint  = "/v1/health"
	VersionEndponit = "/v1/version"
)

type API struct {
	r *chi.Mux
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func New() *API {
	r := newRouter()
	r.Get(HealthEndpoint, healthHandler)
	a := API{
		r: r,
	}
	return &a
}

type JSONResponse struct {
	Message string
	Code    int
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response(w, JSONResponse{Message: "OK", Code: http.StatusOK}, http.StatusOK)
}

func response(w http.ResponseWriter, data interface{}, code int) {
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
