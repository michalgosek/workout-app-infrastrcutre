package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/command"
	"net/http"
)

type HTTP struct {
	addr string
	app  *application.Application
}

func (h *HTTP) CreateUser() http.HandlerFunc {
	type HTTPRequestBody struct {
		Role  string `json:"role"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var payload HTTPRequestBody
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		UUID := uuid.NewString()
		err = h.app.Commands.RegisterUser.Do(r.Context(), command.RegisterUser{
			UUID:  UUID,
			Role:  payload.Role,
			Name:  payload.Name,
			Email: payload.Email,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location:", fmt.Sprintf("%s/v1/users/%s", h.addr, UUID))
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *HTTP) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UUID := chi.URLParam(r, "UUID")
		if UUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing userUUID in path"}, http.StatusBadRequest)
			return
		}
		user, err := h.app.Queries.User.Do(r.Context(), UUID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		rest.SendJSONResponse(w, user, http.StatusOK)
	}
}

func NewHTTP(app *application.Application, addr string) *HTTP {
	if app == nil {
		panic("nil application")
	}
	h := HTTP{
		app:  app,
		addr: addr,
	}
	return &h
}
