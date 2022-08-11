package documents

import "time"

type NotificationWriteModel struct {
	UUID         string    `bson:"_id"`
	UserUUID     string    `bson:"user_uuid"`
	TrainingUUID string    `bson:"training_uuid"`
	Title        string    `bson:"title"`
	Trainer      string    `bson:"trainer"`
	Content      string    `bson:"content"`
	Date         time.Time `bson:"date"`
}
