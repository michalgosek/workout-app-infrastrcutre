package query

import "context"

type AllTrainingGroupReadModel interface {
	AllTrainingGroups(ctx context.Context) ([]TrainingWorkoutGroup, error)
}

type AllTrainingGroupsHandler struct {
	read AllTrainingGroupReadModel
}

func (t *AllTrainingGroupsHandler) Do(ctx context.Context) ([]TrainingWorkoutGroup, error) {
	gg, err := t.read.AllTrainingGroups(ctx)
	if err != nil {
		return nil, err
	}
	return gg, nil
}

func NewAllTrainingGroupsHandler(r AllTrainingGroupReadModel) *AllTrainingGroupsHandler {
	if r == nil {
		panic("nil all training group read model")
	}
	h := AllTrainingGroupsHandler{read: r}
	return &h
}
