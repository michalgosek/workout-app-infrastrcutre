package ports

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/v1/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/v1/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/api-gateway/internal/application/v1/trainer/query"
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
	type HTTPRequestBody struct {
		TrainerUUID string `json:"trainer_uuid"`
		TrainerName string `json:"trainer_name"`
		GroupName   string `json:"group_name"`
		GroupDesc   string `json:"group_desc"`
		Date        string `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var payload HTTPRequestBody
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		err = t.trainerAPI.Commands.Do(r.Context(), command.PlanTraining{
			TrainerUUID: payload.TrainerUUID,
			TrainerName: payload.TrainerName,
			GroupName:   payload.GroupName,
			GroupDesc:   payload.GroupDesc,
			Date:        payload.Date,
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
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		res, err := t.trainerAPI.Queries.Do(r.Context(), query.Training{
			UUID:        trainingUUID,
			TrainerUUID: trainerUUID,
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
