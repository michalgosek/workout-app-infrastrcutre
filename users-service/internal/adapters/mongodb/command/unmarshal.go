package command

import (
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/domain"
)

func UnmarshalToUserWriteModel(u domain.User) documents.UserWriteModel {
	doc := documents.UserWriteModel{
		UUID:           u.UUID(),
		Active:         u.Active(),
		Role:           u.Role(),
		Name:           u.Name(),
		Email:          u.Email(),
		LastActiveDate: u.LastActiveDate(),
	}
	return doc
}

func UnmarshalToUser(d documents.UserWriteModel) domain.User {
	u := domain.UnmarshalUserFromDatabase(domain.DatabaseUser{
		UUID:           d.UUID,
		Active:         d.Active,
		Role:           d.Role,
		Name:           d.Name,
		Email:          d.Email,
		LastActiveDate: d.LastActiveDate,
	})
	return u
}
