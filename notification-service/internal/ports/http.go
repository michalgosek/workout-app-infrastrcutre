package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"notification-service/internal/application"
	"notification-service/internal/application/authorization"
	"notification-service/internal/application/command"
	"notification-service/internal/application/server"
)

type POST struct {
	UserUUID     string `json:"user_uuid"`
	TrainingUUID string `json:"training_uuid"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Trainer      string `json:"trainer"`
	Date         string `json:"date"`
}

type HTTP struct {
	app *application.Application
}

func (h *HTTP) CreateNotificationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := authorization.CreateUserClaimsFromToken(r.Context())
		if err != nil {
			server.SendJSONResponse(w, server.JSONResponse{Message: http.StatusText(http.StatusUnauthorized)}, http.StatusUnauthorized)
			return
		}
		if !claims.HasScope(authorization.NotifyParticipants) {
			server.SendJSONResponse(w, server.JSONResponse{Message: "Insufficient scope."}, http.StatusUnauthorized)
			return
		}

		var payload POST
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&payload)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = h.app.InsertNotificationHandler.Do(r.Context(), command.InsertNotificationCommand{
			UserUUID:     payload.UserUUID,
			TrainingUUID: payload.TrainingUUID,
			Title:        payload.Title,
			Content:      payload.Content,
			Trainer:      payload.Trainer,
			Date:         payload.Date,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *HTTP) GetParticipantNotificationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := authorization.CreateUserClaimsFromToken(r.Context())
		if err != nil {
			server.SendJSONResponse(w, server.JSONResponse{Message: http.StatusText(http.StatusUnauthorized)}, http.StatusUnauthorized)
			return
		}
		if !claims.HasScope(authorization.ViewParticipantNotifications) {
			server.SendJSONResponse(w, server.JSONResponse{Message: "Insufficient scope."}, http.StatusUnauthorized)
			return
		}

		userUUID := chi.URLParam(r, "userUUID")
		if userUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing userUUID in path"}, http.StatusBadRequest)
			return
		}
		notifications, err := h.app.AllNotificationsHandler.Do(r.Context(), userUUID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Println(notifications, userUUID)
		server.SendJSONResponse(w, notifications, http.StatusOK)
	}
}

func (h *HTTP) ClearParticipantNotificationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := authorization.CreateUserClaimsFromToken(r.Context())
		if err != nil {
			server.SendJSONResponse(w, server.JSONResponse{Message: http.StatusText(http.StatusUnauthorized)}, http.StatusUnauthorized)
			return
		}
		if !claims.HasScope(authorization.DeleteParticipantNotifications) {
			server.SendJSONResponse(w, server.JSONResponse{Message: "Insufficient scope."}, http.StatusUnauthorized)
			return
		}
		userUUID := chi.URLParam(r, "userUUID")
		if userUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing userUUID in path"}, http.StatusBadRequest)
			return
		}

		err = h.app.ClearNotificationsHandler.Do(r.Context(), userUUID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func NewHTTP(app *application.Application) (*HTTP, error) {
	if app == nil {
		return nil, errors.New("nil application service implementation")
	}
	h := HTTP{
		app: app,
	}
	return &h, nil
}

func (h *HTTP) NewAPI() chi.Router {
	r := server.NewRouter()
	r.Route("/api/v1/notifications", func(r chi.Router) {
		r.Use(authorization.ValidateJWT())
		r.Route("/{userUUID}", func(r chi.Router) {
			r.Get("/", h.GetParticipantNotificationsHandler())
			r.Post("/", h.CreateNotificationHandler())
			r.Delete("/", h.ClearParticipantNotificationsHandler())
		})
	})
	return r
}
