package inmem

import (
	"errors"
	"sync"

	"github.com/ssouthcity/sweeper"
)

type EventRepository struct {
	mtx    sync.RWMutex
	events map[sweeper.Snowflake]*sweeper.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		events: make(map[sweeper.Snowflake]*sweeper.Event),
	}
}

func (r *EventRepository) Store(evt *sweeper.Event) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.events[evt.ID] = evt
	return nil
}

func (r *EventRepository) Find(id sweeper.Snowflake) (*sweeper.Event, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if evt, ok := r.events[id]; ok {
		return evt, nil
	}
	return nil, errors.New("event does not exist")
}

func (r *EventRepository) FindAll() map[sweeper.Snowflake]*sweeper.Event {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	e := make(map[sweeper.Snowflake]*sweeper.Event, len(r.events))
	for key, val := range r.events {
		e[key] = val
	}
	return e
}

func (r *EventRepository) Remove(id sweeper.Snowflake) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.events, id)
	return nil
}
