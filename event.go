package sweeper

import (
	"errors"
)

var (
	ErrNoOpenSpots        = errors.New("no spots are open for this event")
	ErrUserAlreadyJoined  = errors.New("you are already in this event")
	ErrUserNotParticipant = errors.New("user is not a participant")
	ErrUserIsLeader       = errors.New("you are the fireteam leader of this event")
	ErrEventIsCancelled   = errors.New("event is cancelled")
)

type EventStatus int

const (
	EventStatusSearching EventStatus = iota
	EventStatusFull
	EventStatusCancelled
)

type Event struct {
	ID           Snowflake
	Activity     Activity
	Status       EventStatus
	Description  string
	Participants []*User
}

func (e *Event) Leader() *User {
	return e.Participants[0]
}

func (e *Event) Cancel() {
	e.Status = EventStatusCancelled
}

func (e *Event) AddParticipant(user *User) error {
	switch e.Status {
	case EventStatusCancelled:
		return ErrEventIsCancelled
	case EventStatusFull:
		return ErrNoOpenSpots
	}

	for _, p := range e.Participants {
		if p.ID == user.ID {
			return ErrUserAlreadyJoined
		}
	}

	e.Participants = append(e.Participants, user)

	if len(e.Participants) >= e.Activity.MemberCount() {
		e.Status = EventStatusFull
	}

	return nil
}

func (e *Event) RemoveParticipant(user *User) error {
	if len(e.Participants) == 1 && e.Participants[0].ID == user.ID {
		return ErrUserIsLeader
	}

	for i, p := range e.Participants {
		if p.ID == user.ID {
			e.Participants = append(e.Participants[:i], e.Participants[i+1:]...)
			e.Status = EventStatusSearching
			return nil
		}
	}

	return ErrUserNotParticipant
}

func NewEvent(act Activity, leader *User, desc string) (*Event, error) {
	if desc == "" {
		return nil, errors.New("event may not have an empty description")
	}

	return &Event{
		ID:           NextSnowflake(),
		Activity:     act,
		Status:       EventStatusSearching,
		Description:  desc,
		Participants: append(make([]*User, 0, act.MemberCount()), leader),
	}, nil
}

type EventRepository interface {
	Store(event *Event) error
	FindAll() map[Snowflake]*Event
	Find(id Snowflake) (*Event, error)
	Remove(id Snowflake) error
}
