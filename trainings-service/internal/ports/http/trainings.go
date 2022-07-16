package http

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/service-utility/server/rest"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"net/http"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
)

const (
	InternalMessageErrorMsg = "Internal Message Error."
	ResourceNotFoundMsg     = "Resource not found."
	ServiceUnavailable      = "Service currently unavailable."
)

type Trainings struct {
	app    *application.Application
	format string
}

func (h *Trainings) CreateTrainerWorkoutGroup() http.HandlerFunc {
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
		err = h.app.Commands.ScheduleTrainerWorkoutGroup.Do(r.Context(), command.ScheduleTrainerWorkoutGroup{
			UUID:        uuid.NewString(),
			TrainerUUID: payload.TrainerUUID,
			TrainerName: payload.TrainerName,
			Name:        payload.GroupName,
			Description: payload.GroupDesc,
			Date:        date,
		})

		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Trainings) GetTrainerWorkoutGroup() http.HandlerFunc {
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
		res, err := h.app.Queries.TrainerWorkoutGroup.Do(r.Context(), groupUUID, trainerUUID)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *Trainings) AssignParticipantToWorkoutGroup() http.HandlerFunc {
	type HTTPRequestBody struct {
		ParticipantUUID string `json:"participant_uuid"`
		ParticipantName string `json:"participant_name"`
	}

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

		var payload HTTPRequestBody
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}

		p, err := trainings.NewParticipant(payload.ParticipantUUID, payload.ParticipantName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.app.Commands.AssignParticipantToWorkoutGroup.Do(r.Context(), command.AssignParticipant{
			TrainerUUID: trainerUUID,
			GroupUUID:   groupUUID,
			Participant: p,
		})
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}

//func (h *Trainings) UnassignCustomer() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		groupUUID := chi.URLParam(r, "groupUUID")
//		if groupUUID == "" {
//			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing groupUUID in path"}, http.StatusBadRequest)
//			return
//		}
//		trainerUUID := chi.URLParam(r, "trainerUUID")
//		if trainerUUID == "" {
//			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
//			return
//		}
//		customerUUID := chi.URLParam(r, "customerUUID")
//		if customerUUID == "" {
//			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing customerUUID in path"}, http.StatusBadRequest)
//			return
//		}
//		err := h.app.Commands.UnassignCustomer.Do(r.Context(), command.UnassignCustomerArgs{
//			CustomerUUID: customerUUID,
//			GroupUUID:    groupUUID,
//			TrainerUUID:  trainerUUID,
//		})
//		if err != nil {
//			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//	}
//}
//

func (h *Trainings) GetTrainerWorkoutGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		res, err := h.app.Queries.TrainerWorkoutGroups.Do(r.Context(), trainerUUID)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		rest.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *Trainings) DeleteTrainerWorkoutGroup() http.HandlerFunc {
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
		err := h.app.Commands.CancelTrainerWorkoutGroup.Do(r.Context(), command.CancelWorkoutGroup{
			GroupUUID:   groupUUID,
			TrainerUUID: trainerUUID,
		})
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

//func (h *Trainings) DeleteWorkoutGroups() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		trainerUUID := chi.URLParam(r, "trainerUUID")
//		if trainerUUID == "" {
//			rest.SendJSONResponse(w, rest.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
//			return
//		}
//		err := h.app.Commands.DeleteTrainerWorkouts.Do(r.Context(), trainerUUID)
//		if err != nil {
//			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusNoContent)
//	}
//}

func NewTrainingsHTTP(app *application.Application) *Trainings {
	return &Trainings{
		app: app,
	}
}
