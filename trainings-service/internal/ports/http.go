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

type TrainerScheduleHTTPRequestBody struct {
	TrainerUUID string `json:"trainer_uuid"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Date        string `json:"date"`
}

type TrainerScheduleHTTPResponse struct {
	UUID          string   `json:"schedule_uuid"`
	CustomerUUIDs []string `json:"customer_uuids"`
	Date          string   `json:"date"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Limit         int      `json:"limit"`
}

type DeleteTrainerScheduleHTTPResponse struct {
	UUID string `json:"schedule_uuid"`
}

type DeleteSchedulesHTTPResponse struct {
	UUID string `json:"trainer_uuid"`
}

type TrainerSchedulesHTTPResponse struct {
	Schedules []TrainerScheduleHTTPResponse `json:"schedules"`
}

func NewTrainerSchedulesResponse(format string, schedules ...trainer.TrainerSchedule) TrainerSchedulesHTTPResponse {
	var trainerSchedules []TrainerScheduleHTTPResponse
	for _, s := range schedules {
		trainerSchedules = append(trainerSchedules, TrainerScheduleHTTPResponse{
			UUID:          s.UUID(),
			CustomerUUIDs: s.CustomerUUIDs(),
			Date:          s.Date().Format(format),
			Name:          s.Name(),
			Desc:          s.Desc(),
			Limit:         s.Limit(),
		})
	}
	return TrainerSchedulesHTTPResponse{
		Schedules: trainerSchedules,
	}
}

const (
	internalServerErrMSG = "Internal Server Error"
)

type TrainerService interface {
	CreateSchedule(ctx context.Context, args application.TrainerSchedule) (string, error)
	GetSchedule(ctx context.Context, scheduleUUID, trainerUUID string) (trainer.TrainerSchedule, error)
	GetSchedules(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error)
	AssingCustomer(ctx context.Context, customerUUID, scheduleUUID, trainerUUID string) error
	DeleteSchedule(ctx context.Context, scheduleUUID, trainerUUID string) error
	DeleteSchedules(ctx context.Context, trainerUUID string) error
}
type HTTP struct {
	trainer TrainerService
	format  string
}

func (h *HTTP) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var payload TrainerScheduleHTTPRequestBody
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
	UUID, err := h.trainer.CreateSchedule(r.Context(), application.TrainerSchedule{
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
		Message: fmt.Sprintf("Schedule created with UUID: %s", UUID),
		Code:    http.StatusOK,
	}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func (h *HTTP) GetSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleUUID := r.URL.Query().Get("scheduleUUID")
	if scheduleUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing scheduleUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	schedule, err := h.trainer.GetSchedule(r.Context(), scheduleUUID, trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: internalServerErrMSG, Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	res := TrainerScheduleHTTPResponse{
		UUID:          schedule.UUID(),
		CustomerUUIDs: schedule.CustomerUUIDs(),
		Name:          schedule.Name(),
		Desc:          schedule.Desc(),
		Limit:         schedule.Limit(),
		Date:          schedule.Date().Format(h.format),
	}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func (h *HTTP) GetSchedules(w http.ResponseWriter, r *http.Request) {
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	schedules, err := h.trainer.GetSchedules(r.Context(), trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: internalServerErrMSG, Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	res := NewTrainerSchedulesResponse(h.format, schedules...)
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func (h *HTTP) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	scheduleUUID := r.URL.Query().Get("scheduleUUID")
	if scheduleUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing scheduleUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	err := h.trainer.DeleteSchedule(r.Context(), scheduleUUID, trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Bad request :(", Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}
	res := DeleteTrainerScheduleHTTPResponse{UUID: scheduleUUID}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func (h *HTTP) DeleteSchedules(w http.ResponseWriter, r *http.Request) {
	trainerUUID := r.URL.Query().Get("trainerUUID")
	if trainerUUID == "" {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID query param", Code: http.StatusInternalServerError}, http.StatusBadRequest)
		return
	}
	err := h.trainer.DeleteSchedules(r.Context(), trainerUUID)
	if err != nil {
		rest.SendJSONResponse(w, rest.JSONResponse{Message: "Bad request :(", Code: http.StatusBadRequest}, http.StatusBadRequest)
		return
	}
	res := DeleteSchedulesHTTPResponse{UUID: trainerUUID}
	rest.SendJSONResponse(w, res, http.StatusOK)
}

func NewHTTP(trainer TrainerService, format string) *HTTP {
	h := HTTP{
		trainer: trainer,
		format:  format,
	}
	return &h
}
