package mongo

import (
	"context"

	"github.com/ssouthcity/sweeper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type event struct {
	ID           string   `bson:"_id"`
	Activity     int      `bson:"activity"`
	Description  string   `bson:"description"`
	Participants []string `bson:"participants"`
}

type eventRepository struct {
	events *mongo.Collection
	users  sweeper.UserRepository
}

func (r *eventRepository) Store(e *sweeper.Event) error {
	evt := encodeEvent(e)

	filter := bson.M{"_id": evt.ID}
	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": evt}

	_, err := r.events.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *eventRepository) Find(id sweeper.Snowflake) (*sweeper.Event, error) {
	e := &event{}

	cur := r.events.FindOne(context.Background(), bson.M{"_id": string(id)})

	if err := cur.Decode(e); err != nil {
		return nil, err
	}

	return decodeEvent(e, r.users), nil
}

func (r *eventRepository) FindAll() map[sweeper.Snowflake]*sweeper.Event {
	evts := make(map[sweeper.Snowflake]*sweeper.Event)

	cur, err := r.events.Find(context.Background(), nil)
	if err != nil {
		return evts
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		e := &event{}

		if err := cur.Decode(e); err != nil {
			continue
		}

		de := decodeEvent(e, r.users)

		evts[de.ID] = de
	}

	return evts
}

func (r *eventRepository) Remove(id sweeper.Snowflake) error {
	_, err := r.events.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func NewEventRepository(events *mongo.Collection, users sweeper.UserRepository) sweeper.EventRepository {
	return &eventRepository{
		events: events,
		users:  users,
	}
}

func encodeEvent(e *sweeper.Event) *event {
	return &event{
		ID:           string(e.ID),
		Activity:     int(e.Activity),
		Description:  e.Description,
		Participants: encodeParticipants(e),
	}
}

func encodeParticipants(e *sweeper.Event) []string {
	var participants []string

	for _, p := range e.Participants {
		participants = append(participants, string(p.ID))
	}

	return participants
}

func decodeEvent(e *event, users sweeper.UserRepository) *sweeper.Event {
	return &sweeper.Event{
		ID:           sweeper.Snowflake(e.ID),
		Activity:     sweeper.Activity(e.Activity),
		Description:  e.Description,
		Participants: decodeParticipants(e, users),
	}
}

func decodeParticipants(e *event, users sweeper.UserRepository) []*sweeper.User {
	var result []*sweeper.User

	for _, p := range e.Participants {
		u, err := users.Find(sweeper.Snowflake(p))
		if err != nil {
			continue
		}

		result = append(result, u)
	}

	return result
}
