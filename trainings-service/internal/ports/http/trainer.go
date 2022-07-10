package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
)

const (
	InternalMessageErrorMsg = "Internal Message Error."
	ResourceNotFoundMsg     = "Resource not found."
	ServiceUnavailable      = "Service currently unavailable."
)

type TrainerWorkoutGroups struct {
	app    *application.Application
	format string
}

func (h *TrainerWorkoutGroups) CreateTrainerWorkoutGroup() http.HandlerFunc {
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
		date, err := time.Parse(h.format, payload.Date)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		err = h.app.Commands.CreateTrainerWorkout.Do(r.Context(), command.ScheduleWorkout{
			TrainerUUID: payload.TrainerUUID,
			GroupName:   payload.GroupName,
			GroupDesc:   payload.GroupDesc,
			TrainerName: payload.TrainerName,
			Date:        date,
		})

		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *TrainerWorkoutGroups) UnassignCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupUUID := chi.URLParam(r, "groupUUID")
		if groupUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing groupUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		customerUUID := chi.URLParam(r, "customerUUID")
		if customerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing customerUUID in path"}, http.StatusBadRequest)
			return
		}
		err := h.app.Commands.UnassignCustomer.Do(r.Context(), command.UnassignCustomer{
			CustomerUUID: customerUUID,
			GroupUUID:    groupUUID,
			TrainerUUID:  trainerUUID,
		})
		if errors.Is(err, command.ErrResourceNotFound) {
			http.Error(w, ResourceNotFoundMsg, http.StatusNotFound)
			return
		}
		if errors.Is(err, command.ErrWorkoutGroupNotOwner) {
			http.Error(w, ResourceNotFoundMsg, http.StatusNotFound)
			return
		}
		if errors.Is(err, command.ErrRepositoryFailure) {
			http.Error(w, ServiceUnavailable, http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *TrainerWorkoutGroups) GetTrainerWorkoutGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupUUID := chi.URLParam(r, "groupUUID")
		if groupUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing groupUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		res, err := h.app.Queries.GetTrainerWorkout.Do(r.Context(), trainerUUID, groupUUID)
		if errors.Is(err, query.ErrWorkoutGroupNotOwner) {
			http.Error(w, ResourceNotFoundMsg, http.StatusNotFound)
			return
		}
		if errors.Is(err, query.ErrRepositoryFailure) {
			http.Error(w, ServiceUnavailable, http.StatusServiceUnavailable)
		}
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *TrainerWorkoutGroups) GetTrainerWorkoutGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		res, err := h.app.Queries.GetTrainerWorkouts.Do(r.Context(), trainerUUID)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *TrainerWorkoutGroups) DeleteWorkoutGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupUUID := chi.URLParam(r, "groupUUID")
		if groupUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing groupUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		err := h.app.Commands.DeleteTrainerWorkout.Do(r.Context(), command.CancelWorkoutArgs{
			GroupUUID:   groupUUID,
			TrainerUUID: trainerUUID,
		})
		if errors.Is(err, command.ErrWorkoutGroupNotOwner) {
			http.Error(w, ResourceNotFoundMsg, http.StatusNotFound)
			return
		}
		if errors.Is(err, command.ErrRepositoryFailure) {
			http.Error(w, ServiceUnavailable, http.StatusServiceUnavailable)
		}
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *TrainerWorkoutGroups) DeleteWorkoutGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		err := h.app.Commands.DeleteTrainerWorkouts.Do(r.Context(), trainerUUID)
		if errors.Is(err, command.ErrRepositoryFailure) {
			http.Error(w, ServiceUnavailable, http.StatusServiceUnavailable)
		}
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func NewTrainerWorkoutGroupsHTTP(app *application.Application, format string) *TrainerWorkoutGroups {
	return &TrainerWorkoutGroups{
		app:    app,
		format: format,
	}
}
