package mongodb

import "time"

type UserWriteModel struct {
	UUID           string    `bson:"_id"`
	Active         bool      `bson:"active"`
	Role           string    `bson:"role"`
	Name           string    `bson:"name"`
	LastActiveDate time.Time `bson:"last_active_date"`
}
