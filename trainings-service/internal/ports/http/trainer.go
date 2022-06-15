package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type TrainerService interface {
	CreateWorkoutGroup(ctx context.Context, args application.TrainerSchedule) (string, error)
	GetWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) (trainer.WorkoutGroup, error)
	AssignCustomer(ctx context.Context, args application.WorkoutRegistration) error
	GetWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
	DeleteWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) error
	DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error
}

type TrainerWorkoutGroups struct {
	service TrainerService
	format  string
}

func (h *TrainerWorkoutGroups) CreateTrainerWorkoutGroup(w http.ResponseWriter, r *http.Request) {
	var payload TrainerWorkoutGroupHTTPRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&payload)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Internal Server Error", Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	date, err := time.Parse(h.format, payload.Date)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Internal Server Error", Code: http.StatusInternalServerError}, http.StatusInternalServerError)
	}
	UUID, err := h.service.CreateWorkoutGroup(r.Context(), application.TrainerSchedule{
		TrainerUUID: payload.TrainerUUID,
		Name:        payload.Name,
		Desc:        payload.Desc,
		Date:        date,
	})
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: err.Error(), Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	res := rest.JSONResponse{
		Message: fmt.Sprintf("WorkoutSchedule created with UUID: %s", UUID),
		Code:    http.StatusOK,
	}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func (h *TrainerWorkoutGroups) GetTrainerWorkoutGroup(w http.ResponseWriter, r *http.Request) {
	workoutGroupUUID := r.URL.Query().Get("workoutGroupUUID")
	if workoutGroupUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing workoutGroupUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	schedule, err := h.service.GetWorkoutGroup(r.Context(), workoutGroupUUID, trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Internal Server Error", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	res := TrainerWorkoutGroupHTTPResponse{
		UUID:          schedule.UUID(),
		CustomerUUIDs: schedule.CustomerUUIDs(),
		Name:          schedule.Name(),
		Desc:          schedule.Desc(),
		Limit:         schedule.Limit(),
		Date:          schedule.Date().Format(h.format),
	}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func (h *TrainerWorkoutGroups) GetTrainerWorkoutGroups(w http.ResponseWriter, r *http.Request) {
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	schedules, err := h.service.GetWorkoutGroups(r.Context(), trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Internal Server Error", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	res := newTrainerTrainerWorkoutGroupsHTTPResponse(h.format, trainerUUID, schedules...)
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func (h *TrainerWorkoutGroups) DeleteWorkoutGroup(w http.ResponseWriter, r *http.Request) {
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	workoutGroupUUID := r.URL.Query().Get("workoutGroupUUID")
	if workoutGroupUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing workoutGroupUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	err := h.service.DeleteWorkoutGroup(r.Context(), workoutGroupUUID, trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Bad request :(", Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}
	res := DeleteTrainerWorkoutGroupHTTPResponse{UUID: workoutGroupUUID}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func (h *TrainerWorkoutGroups) DeleteWorkoutGroups(w http.ResponseWriter, r *http.Request) {
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	err := h.service.DeleteWorkoutGroups(r.Context(), trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Bad request :(", Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}
	res := DeleteTrainerWorkoutGroupsHTTPResponse{TrainerUUID: trainerUUID}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func NewTrainerWorkoutGroupsHTTP(service TrainerService, format string) *TrainerWorkoutGroups {
	return &TrainerWorkoutGroups{
		service: service,
		format:  format,
	}
}

func newTrainerTrainerWorkoutGroupsHTTPResponse(format, trainerUUID string, groups ...trainer.WorkoutGroup) TrainerWorkoutGroupsHTTPResponse {
	var trainerWorkoutGroups []TrainerWorkoutGroupHTTPResponse
	for _, s := range groups {
		trainerWorkoutGroups = append(trainerWorkoutGroups, TrainerWorkoutGroupHTTPResponse{
			UUID:          s.UUID(),
			CustomerUUIDs: s.CustomerUUIDs(),
			Date:          s.Date().Format(format),
			Name:          s.Name(),
			Desc:          s.Desc(),
			Limit:         s.Limit(),
		})
	}
	return TrainerWorkoutGroupsHTTPResponse{
		WorkoutGroups: trainerWorkoutGroups,
		TrainerUUID:   trainerUUID,
	}
}
