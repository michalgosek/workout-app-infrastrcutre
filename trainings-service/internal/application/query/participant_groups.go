package query

import "context"

type ParticipantGroupsReadModel interface {
	ParticipantGroups(ctx context.Context, UUID string) ([]ParticipantGroup, error)
}

type ParticipantGroupsHandler struct {
	read ParticipantGroupsReadModel
}

func (t *ParticipantGroupsHandler) Do(ctx context.Context, UUID string) ([]ParticipantGroup, error) {
	gg, err := t.read.ParticipantGroups(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return gg, nil
}

func NewParticipantGroupsHandler(r ParticipantGroupsReadModel) *ParticipantGroupsHandler {
	if r == nil {
		panic("nil participant groups read model")
	}
	h := ParticipantGroupsHandler{read: r}
	return &h
}
