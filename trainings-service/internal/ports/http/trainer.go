package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	command2 "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"

	"github.com/go-chi/chi"

	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
)

const (
	InternalMessageErrorMsg = "Internal Message Error"
)

type TrainerWorkoutGroups struct {
	app    *application.Application
	format string
}

func (h *TrainerWorkoutGroups) CreateTrainerWorkoutGroup() http.HandlerFunc {
	type HTTPRequestBody struct {
		TrainerUUID string `json:"trainer_uuid"`
		Name        string `json:"name"`
		Desc        string `json:"desc"`
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
		UUID, err := h.app.Commands.CreateTrainerWorkout.Do(r.Context(), command2.WorkoutGroup{
			TrainerUUID: payload.TrainerUUID,
			Name:        payload.Name,
			Desc:        payload.Desc,
			Date:        date,
		})
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		res := rest.JSONResponse{Message: fmt.Sprintf("workout group created with UUID: %s", UUID)}
		rest.SendJSONResponse(w, res, http.StatusCreated)
	}
}

func (h *TrainerWorkoutGroups) AssignCustomer() http.HandlerFunc {
	type HTTPRequestBody struct {
		TrainerUUID  string `json:"trainer_uuid"`
		UUID         string `json:"workout_group_uuid"`
		CustomerUUID string `json:"customer_uuid"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var payload HTTPRequestBody
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		err = h.app.Commands.AssignCustomer.Do(r.Context(), command2.WorkoutRegistration{
			TrainerUUID:  payload.TrainerUUID,
			GroupUUID:    payload.UUID,
			CustomerUUID: payload.CustomerUUID,
		})
		if errors.Is(err, application.ErrScheduleNotOwner) {
			res := rest.JSONResponse{Message: fmt.Sprintf("workout group created with UUID: %s", payload.UUID)}
			rest.SendJSONResponse(w, res, http.StatusBadRequest)
			return
		}
		if errors.Is(err, application.ErrRepositoryFailure) {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}

		res := rest.JSONResponse{
			Message: fmt.Sprintf("Customer UUID: %s assgined to wrokout group UUID: %s", payload.CustomerUUID, payload.UUID),
		}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *TrainerWorkoutGroups) GetTrainerWorkoutGroup() http.HandlerFunc {
	type HTTPResponseBody struct {
		UUID          string   `json:"workout_group_uuid"`
		CustomerUUIDs []string `json:"customer_uuids"`
		Date          string   `json:"date"`
		Name          string   `json:"name"`
		Desc          string   `json:"desc"`
		Limit         int      `json:"limit"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		workoutUUID := chi.URLParam(r, "workoutUUID")
		if workoutUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing workoutUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}

		group, err := h.app.Queries.GetTrainerWorkout.Do(r.Context(), workoutUUID, trainerUUID)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		res := HTTPResponseBody{
			UUID:          group.UUID(),
			CustomerUUIDs: group.CustomerUUIDs(),
			Name:          group.Name(),
			Desc:          group.Desc(),
			Limit:         group.Limit(),
			Date:          group.Date().Format(h.format),
		}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *TrainerWorkoutGroups) GetTrainerWorkoutGroups() http.HandlerFunc {
	type WorkoutGroup struct {
		UUID          string   `json:"workout_group_uuid"`
		CustomerUUIDs []string `json:"customer_uuids"`
		Date          string   `json:"date"`
		Name          string   `json:"name"`
		Desc          string   `json:"desc"`
		Limit         int      `json:"limit"`
	}

	type HTTPResponseBody struct {
		TrainerUUID   string         `json:"trainer_uuid"`
		WorkoutGroups []WorkoutGroup `json:"workout_groups"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		groups, err := h.app.Queries.GetTrainerWorkouts.Do(r.Context(), trainerUUID)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}

		resp := HTTPResponseBody{TrainerUUID: trainerUUID}
		for _, g := range groups {
			resp.WorkoutGroups = append(resp.WorkoutGroups, WorkoutGroup{
				UUID:          g.UUID(),
				CustomerUUIDs: g.CustomerUUIDs(),
				Date:          g.Date().Format(h.format),
				Name:          g.Name(),
				Desc:          g.Desc(),
				Limit:         g.Limit(),
			})
		}
		rest.SendJSONResponse(w, resp, http.StatusOK)
	}
}

func (h *TrainerWorkoutGroups) DeleteWorkoutGroup() http.HandlerFunc {
	type HTTPResponseBody struct {
		UUID string `json:"workout_group_uuid"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		workoutUUID := chi.URLParam(r, "workoutUUID")
		if workoutUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing workoutUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		err := h.app.Commands.DeleteTrainerWorkout.Do(r.Context(), workoutUUID, trainerUUID)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		res := HTTPResponseBody{UUID: workoutUUID}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *TrainerWorkoutGroups) DeleteWorkoutGroups() http.HandlerFunc {
	type HTTPResponseBody struct {
		TrainerUUID string `json:"trainer_uuid_uuid"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		err := h.app.Commands.DeleteTrainerWorkouts.Do(r.Context(), trainerUUID)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		res := HTTPResponseBody{TrainerUUID: trainerUUID}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func NewTrainerWorkoutGroupsHTTP(app *application.Application, format string) *TrainerWorkoutGroups {
	return &TrainerWorkoutGroups{
		app:    app,
		format: format,
	}
}
