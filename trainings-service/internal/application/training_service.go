package application

type WorkoutSessionsRepoistory interface {
}

type TraningsService struct {
	repository WorkoutSessionsRepoistory
}
