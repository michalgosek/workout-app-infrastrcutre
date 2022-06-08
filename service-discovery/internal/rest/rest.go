package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	HealthEndpoint  = "/v1/health"
	VersionEndponit = "/v1/version"
)

type Config struct {
	MiddlewareTimeout time.Duration
}

type API struct {
	router          chi.Router
	registerHandler *RegisterHandler
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *API) SetEndpoints() {
	a.router.Get(HealthEndpoint, healthHandler)
	a.router.Get(a.registerHandler.ClientRegistryEndpoint(), a.registerHandler.QueryInstances)
	a.router.Post(a.registerHandler.ServiceRegiststryEndpoint(), a.registerHandler.ServiceRegistry)
}

func NewAPI(h *RegisterHandler) *API {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	a := API{
		router:          r,
		registerHandler: h,
	}
	return &a
}

type JSONResponse struct {
	Message string
	Code    int
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
