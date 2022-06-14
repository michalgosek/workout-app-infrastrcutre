package ports

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

const (
	internalServerErrMSG = "Internal Server Error"
)

type TrainerService interface {
	CreateWorkoutGroup(ctx context.Context, args application.TrainerSchedule) (string, error)
	GetWorkoutGroup(ctx context.Context, scheduleUUID, trainerUUID string) (trainer.WorkoutGroup, error)
	AssignCustomer(ctx context.Context, customerUUID, workoutGroupUUID, trainerUUID string) error
	GetWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
	DeleteWorkoutGroup(ctx context.Context, workoutGroupUUID, trainerUUID string) error
	DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error
}

type HTTP struct {
	trainer TrainerService
	format  string
}

type TrainerWorkoutGroupHTTPRequest struct {
	TrainerUUID string `json:"trainer_uuid"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Date        string `json:"date"`
}

func (h *HTTP) CreateTrainerWorkoutGroup(w http.ResponseWriter, r *http.Request) {
	var payload TrainerWorkoutGroupHTTPRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&payload)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: internalServerErrMSG, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
		return
	}
	date, err := time.Parse(h.format, payload.Date)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: internalServerErrMSG, Code: http.StatusInternalServerError}, http.StatusInternalServerError)
	}
	UUID, err := h.trainer.CreateWorkoutGroup(r.Context(), application.TrainerSchedule{
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

type TrainerWorkoutGroupHTTPResponse struct {
	UUID          string   `json:"workout_group_uuid"`
	CustomerUUIDs []string `json:"customer_uuids"`
	Date          string   `json:"date"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Limit         int      `json:"limit"`
}

func (h *HTTP) GetTrainerWorkoutGroup(w http.ResponseWriter, r *http.Request) {
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
	schedule, err := h.trainer.GetWorkoutGroup(r.Context(), workoutGroupUUID, trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: internalServerErrMSG, Code: http.StatusInternalServerError}, http.StatusBadRequest)
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

type TrainerWorkoutGroupsHTTPResponse struct {
	TrainerUUID   string                            `json:"trainer_uuid"`
	WorkoutGroups []TrainerWorkoutGroupHTTPResponse `json:"workout_groups"`
}

func NewTrainerTrainerWorkoutGroupsHTTPResponse(format, trainerUUID string, groups ...trainer.WorkoutGroup) TrainerWorkoutGroupsHTTPResponse {
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

func (h *HTTP) GetTrainerWorkoutGroups(w http.ResponseWriter, r *http.Request) {
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	schedules, err := h.trainer.GetWorkoutGroups(r.Context(), trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: internalServerErrMSG, Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	res := NewTrainerTrainerWorkoutGroupsHTTPResponse(h.format, trainerUUID, schedules...)
	rest.SendJSONResponse(w, res, http.StatusOK)
}

type DeleteTrainerWorkoutGroupHTTPResponse struct {
	UUID string `json:"workout_group_uuid"`
}

func (h *HTTP) DeleteWorkoutGroup(w http.ResponseWriter, r *http.Request) {
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
	err := h.trainer.DeleteWorkoutGroup(r.Context(), workoutGroupUUID, trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Bad request :(", Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}
	res := DeleteTrainerWorkoutGroupHTTPResponse{UUID: workoutGroupUUID}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

type DeleteTrainerWorkoutGroupsHTTPResponse struct {
	TrainerUUID string `json:"trainer_uuid_uuid"`
}

func (h *HTTP) DeleteWorkoutGroups(w http.ResponseWriter, r *http.Request) {
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	err := h.trainer.DeleteWorkoutGroups(r.Context(), trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Bad request :(", Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}
	res := DeleteTrainerWorkoutGroupsHTTPResponse{TrainerUUID: trainerUUID}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func NewHTTP(trainer TrainerService, format string) *HTTP {
	h := HTTP{
		trainer: trainer,
		format:  format,
	}
	return &h
}
