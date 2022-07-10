package http

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"net/http"
)

type CustomerHTTP struct {
	app    *application.Application
	format string
}

func (c *CustomerHTTP) CreateCustomerWorkout() http.HandlerFunc {
	type HTTPRequestBody struct {
		GroupUUID    string `json:"group_uuid"`
		CustomerName string `json:"customer_name"`
		TrainerUUID  string `json:"trainer_uuid"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		customerUUID := chi.URLParam(r, "customerUUID")
		if customerUUID == "" {
			http.Error(w, "Missing customerUUID in the path", http.StatusBadRequest)
			return
		}

		var payload HTTPRequestBody
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		err = c.app.Commands.CustomerScheduleWorkout.Do(r.Context(), command.ScheduleWorkoutArgs{
			CustomerUUID: customerUUID,
			CustomerName: payload.CustomerName,
			GroupUUID:    payload.GroupUUID,
			TrainerUUID:  payload.TrainerUUID,
		})
		if err != nil {
			http.Error(w, InternalMessageErrorMsg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func NewCustomerHTTP(app *application.Application, format string) *CustomerHTTP {
	return &CustomerHTTP{
		app:    app,
		format: format,
	}
}
