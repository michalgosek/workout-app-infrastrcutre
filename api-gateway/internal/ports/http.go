package ports

import (
	"encoding/json"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/api/v1/rest/trainer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/auth"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"net/http"
)

const (
	InternalMessageErrorMsg = "Internal Message Error."
)

type TrainerHTTP struct {
	trainerAPI *trainer.Application
}

func (t *TrainerHTTP) CreateTraining() http.HandlerFunc {
	type User struct {
		UserUUID string `json:"uuid"`
		Role     string `json:"role"`
		Name     string `json:"name"`
	}
	type TrainingGroup struct {
		User      User   `json:"user"`
		GroupName string `json:"group_name"`
		GroupDesc string `json:"group_desc"`
		Date      string `json:"date"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*auth.CustomClaims)
		if !claims.HasScope("create:training-group") {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message":"Insufficient scope."}`))
			return
		}

		var payload TrainingGroup
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		err = t.trainerAPI.Commands.Do(r.Context(), command.PlanTrainingCommand{
			User: command.User{
				UUID: payload.User.UserUUID,
				Role: payload.User.Role,
				Name: payload.User.Name,
			},
			GroupName: payload.GroupName,
			GroupDesc: payload.GroupDesc,
			Date:      payload.Date,
		})
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (t *TrainerHTTP) GetTraining() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainingUUID := chi.URLParam(r, "trainingUUID")
		if trainingUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing groupUUID in path"}, http.StatusBadRequest)
			return
		}
		userUUID := chi.URLParam(r, "userUUID")
		if userUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		res, err := t.trainerAPI.Queries.Do(r.Context(), query.Training{
			UserUUID:     userUUID,
			TrainingUUID: trainingUUID,
		})
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func NewTrainerHTTP(t *trainer.Application) *TrainerHTTP {
	h := TrainerHTTP{trainerAPI: t}
	return &h
}
