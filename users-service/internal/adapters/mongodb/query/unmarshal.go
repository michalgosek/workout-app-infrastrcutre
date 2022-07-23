package query

import (
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/query"
)

func UnmarshalToQueryUser(d documents.UserWriteModel) query.User {
	g := query.User{
		Role:  d.Role,
		Name:  d.Name,
		Email: d.Email,
	}
	return g
}
