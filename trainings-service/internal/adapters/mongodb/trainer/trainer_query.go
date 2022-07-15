package trainer

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QueryHandlerConfig struct {
	Collection   string
	Database     string
	Format       string
	QueryTimeout time.Duration
}

type QueryHandler struct {
	cli *mongo.Client
	cfg QueryHandlerConfig
}

func NewQueryHandler(cli *mongo.Client, cfg QueryHandlerConfig) *QueryHandler {
	t := QueryHandler{
		cli: cli,
		cfg: cfg,
	}
	return &t
}

func (t *QueryHandler) TrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]query.TrainerWorkoutGroup, error) {
	f := bson.M{"trainer_uuid": trainerUUID}
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	res, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	var docs []WorkoutGroupWriteModel
	err = res.All(ctx, &docs)
	if err != nil {
		return nil, err
	}
	out := unmarshalTrainerWorkoutGroup(docs)
	return out, nil
}

func (t *QueryHandler) TrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) (query.TrainerWorkoutGroup, error) {
	f := bson.M{"trainer_uuid": trainerUUID, "_id": groupUUID}
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)

	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return query.TrainerWorkoutGroup{}, res.Err()
	}

	var dst WorkoutGroupWriteModel
	err := res.Decode(&dst)
	if err != nil {
		return query.TrainerWorkoutGroup{}, err
	}
	out := query.TrainerWorkoutGroup{
		TrainerUUID:  dst.TrainerUUID,
		TrainerName:  dst.TrainerName,
		GroupUUID:    dst.UUID,
		GroupDesc:    dst.Description,
		GroupName:    dst.Name,
		Date:         dst.Date,
		Participants: unmarshalWorkoutGroupParticipants(dst.Participants),
	}
	return out, nil
}

func unmarshalTrainerWorkoutGroup(docs []WorkoutGroupWriteModel) []query.TrainerWorkoutGroup {
	var out []query.TrainerWorkoutGroup
	for _, d := range docs {
		out = append(out, query.TrainerWorkoutGroup{
			TrainerUUID:  d.TrainerUUID,
			TrainerName:  d.TrainerName,
			GroupUUID:    d.UUID,
			GroupDesc:    d.Description,
			GroupName:    d.Name,
			Date:         d.Date,
			Participants: unmarshalWorkoutGroupParticipants(d.Participants),
		})
	}
	return out
}

func unmarshalWorkoutGroupParticipants(docs []WorkoutGroupParticipantWriteModel) []query.TrainerWorkoutGroupParticipant {
	var out []query.TrainerWorkoutGroupParticipant
	for _, d := range docs {
		out = append(out, query.TrainerWorkoutGroupParticipant{
			UUID: d.UUID,
			Name: d.Name,
		})
	}
	return out
}

// ---> refactor ----<
func (t *QueryHandler) queryTrainerWorkoutGroupWithFilter(ctx context.Context, filter bson.M) (trainer.WorkoutGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	res := coll.FindOne(ctx, filter)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return trainer.WorkoutGroup{}, nil
	}
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("find one failed: %v", err)
	}
	var doc WorkoutGroupWriteModel
	err = res.Decode(&doc)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("decoding failed: %v", err)
	}

	date, err := time.Parse(t.cfg.Format, doc.Date)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("parsing date value from document failed: %v", err)
	}
	var cd []customer.Details
	for _, c := range doc.Participants {
		d, err := customer.UnmarshalCustomerDetails(c.UUID, c.Name)
		if err != nil {
			return trainer.WorkoutGroup{}, fmt.Errorf("unmarshal customer details failed: %v", err)
		}
		cd = append(cd, d)
	}
	group, err := trainer.UnmarshalWorkoutGroupFromDatabase(trainer.WorkoutGroupDetails{
		UUID:        doc.UUID,
		TrainerUUID: doc.TrainerUUID,
		TrainerName: doc.TrainerName,
		Name:        doc.Name,
		Description: doc.Description,
		Date:        date,
		Limit:       doc.Limit,
	}, cd)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("unmarshal workout group from database failed: %v", err)
	}

	return group, nil
}

func (t *QueryHandler) QueryCustomerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID, customerUUID string) (trainer.WorkoutGroup, error) {
	f := bson.M{"_id": groupUUID, "trainer_uuid": trainerUUID, "participants.uuid": customerUUID}
	return t.queryTrainerWorkoutGroupWithFilter(ctx, f)
}

func (t *QueryHandler) QueryTrainerWorkoutGroupWithDate(ctx context.Context, trainerUUID string, date time.Time) (trainer.WorkoutGroup, error) {
	f := bson.M{"trainer_uuid": trainerUUID, "date": date.Format(t.cfg.Format)}
	return t.queryTrainerWorkoutGroupWithFilter(ctx, f)
}

func (t *QueryHandler) QueryTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) (trainer.WorkoutGroup, error) {
	f := bson.M{"_id": groupUUID, "trainer_uuid": trainerUUID}
	return t.queryTrainerWorkoutGroupWithFilter(ctx, f)
}

func (t *QueryHandler) QueryTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	f := bson.M{"trainer_uuid": trainerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("find failed: %v", err)
	}

	var docs []WorkoutGroupWriteModel
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %v", err)
	}

	var groups []trainer.WorkoutGroup
	for _, d := range docs { // O(n^2)
		date, err := time.Parse(t.cfg.Format, d.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date value from document failed: %v", err)
		}
		var details []customer.Details
		for _, c := range d.Participants {
			d, err := customer.UnmarshalCustomerDetails(c.UUID, c.Name)
			if err != nil {
				return nil, fmt.Errorf("unmarshal customer details failed: %v", err)
			}
			details = append(details, d)
		}

		group, err := trainer.UnmarshalWorkoutGroupFromDatabase(trainer.WorkoutGroupDetails{
			UUID:        d.UUID,
			TrainerUUID: d.TrainerUUID,
			TrainerName: d.TrainerName,
			Name:        d.Name,
			Description: d.Description,
			Date:        date,
			Limit:       d.Limit,
		}, details)
		if err != nil {
			return nil, fmt.Errorf("unmarshal workout group from database failed: %v", err)
		}
		groups = append(groups, group)
	}
	return groups, nil
}
