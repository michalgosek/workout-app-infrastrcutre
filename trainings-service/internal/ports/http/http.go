package http

import (
	"encoding/json"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/authorization"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/server"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"net/http"
	"time"
)

// fixme: add claims on auth0
type Trainings struct {
	addr   string
	app    *application.Application
	format string
	router chi.Router
}

func (h *Trainings) UpdateTrainingGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainingUUID := chi.URLParam(r, "trainingUUID")
		if trainingUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing groupUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}

		var payload UpdateTrainingGroupPost
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		date, err := time.Parse(h.format, payload.Date)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		updateCMD := command.UpdateTrainingGroup{
			TrainerUUID:  trainerUUID,
			TrainingUUID: trainingUUID,
			Name:         payload.GroupName,
			Description:  payload.GroupDesc,
			Date:         date,
		}
		err = updateCMD.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.app.Commands.UpdateTrainingGroup.Do(r.Context(), updateCMD)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Trainings) CreateTrainingGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*authorization.CustomClaims)
		if !claims.HasScope("create:training-group") {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message":"Insufficient scope."}`))
			return
		}

		var payload TrainingGroupPost
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		t, err := trainings.NewTrainer(payload.User.UUID, payload.User.Name)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		date, err := time.Parse(h.format, payload.Date)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = h.app.Commands.PlanTrainingGroup.Do(r.Context(), command.PlanTrainingGroup{
			UUID:        uuid.NewString(),
			Trainer:     t,
			Name:        payload.GroupName,
			Description: payload.GroupDesc,
			Date:        date,
		})

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Trainings) GetTrainerGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainingUUID := chi.URLParam(r, "trainingUUID")
		if trainingUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing groupUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		res, err := h.app.Queries.TrainerGroup.Do(r.Context(), trainingUUID, trainerUUID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		server.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *Trainings) AssignParticipant() http.HandlerFunc {
	type HTTPRequestBody struct {
		ParticipantUUID string `json:"participant_uuid"`
		ParticipantName string `json:"participant_name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		trainingUUID := chi.URLParam(r, "trainingUUID")
		if trainingUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing groupUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}

		var payload HTTPRequestBody
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&payload)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		p, err := trainings.NewParticipant(payload.ParticipantUUID, payload.ParticipantName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.app.Commands.AssignParticipant.Do(r.Context(), command.AssignParticipant{
			TrainerUUID:  trainerUUID,
			TrainingUUID: trainingUUID,
			Participant:  p,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Trainings) UnassignParticipant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainingUUID := chi.URLParam(r, "trainingUUID")
		if trainingUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainingUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		participantUUID := chi.URLParam(r, "participantUUID")
		if participantUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing participantUUID in path"}, http.StatusBadRequest)
			return
		}
		err := h.app.Commands.UnassignParticipant.Do(r.Context(), command.UnassignParticipant{
			ParticipantUUID: participantUUID,
			TrainingUUID:    trainingUUID,
			TrainerUUID:     trainerUUID,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Trainings) GetTrainerGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		res, err := h.app.Queries.TrainerGroups.Do(r.Context(), trainerUUID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		server.SendJSONResponse(w, res, http.StatusOK)
	}
}

func (h *Trainings) DeleteTrainerGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainingUUID := chi.URLParam(r, "trainingUUID")
		if trainingUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainingUUID in path"}, http.StatusBadRequest)
			return
		}
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		err := h.app.Commands.CancelTrainingGroup.Do(r.Context(), command.CancelWorkoutGroup{
			TrainingUUID: trainingUUID,
			TrainerUUID:  trainerUUID,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Trainings) DeleteTrainerGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainerUUID := chi.URLParam(r, "trainerUUID")
		if trainerUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing trainerUUID in path"}, http.StatusBadRequest)
			return
		}
		err := h.app.Commands.CancelTrainingGroups.Do(r.Context(), trainerUUID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Trainings) GetAllTrainingGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := h.app.AllTrainingGroups.Do(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		server.SendJSONResponse(w, groups, http.StatusOK)
	}
}

func (h *Trainings) GetParticipantGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		participantUUID := chi.URLParam(r, "participantUUID")
		if participantUUID == "" {
			server.SendJSONResponse(w, server.JSONResponse{Message: "missing participantUUID in path"}, http.StatusBadRequest)
			return
		}
		res, err := h.app.Queries.ParticipantGroups.Do(r.Context(), participantUUID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		server.SendJSONResponse(w, res, http.StatusOK)
	}
}

func NewTrainingsHTTP(app *application.Application, addr string) *Trainings {
	return &Trainings{
		app:    app,
		addr:   addr,
		format: query.UIFormat,
		router: server.NewRouter(),
	}
}

func (h *Trainings) NewAPI() chi.Router {
	h.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/trainings", func(r chi.Router) {
			r.Get("/", h.GetAllTrainingGroups())
		})
		r.Route("/participants", h.participantRoutes())
		r.Route("/trainers", h.trainerRoutes())
	})
	return h.router
}

func (h *Trainings) participantRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(authorization.EnsureValidToken())
		r.Route("/{participantUUID}", func(r chi.Router) {
			r.Route("/trainings", func(r chi.Router) {
				r.Get("/", h.GetParticipantGroups())
			})
		})
	}
}

func (h *Trainings) trainerRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(authorization.EnsureValidToken())
		r.Route("/{trainerUUID}", func(r chi.Router) {
			r.Route("/trainings", func(r chi.Router) {
				r.Post("/", h.CreateTrainingGroup())
				r.Get("/", h.GetTrainerGroups())
				r.Delete("/", h.DeleteTrainerGroups())
				r.Route("/{trainingUUID}", func(r chi.Router) {
					r.Put("/", h.UpdateTrainingGroup())
					r.Get("/", h.GetTrainerGroup())
					r.Delete("/", h.DeleteTrainerGroup())
					r.Route("/participants", func(r chi.Router) {
						r.Post("/", h.AssignParticipant())
						r.Route("/{participantUUID}", func(r chi.Router) {
							r.Delete("/", h.UnassignParticipant())
						})
					})
				})
			})
		})
	}
}
