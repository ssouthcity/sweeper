package sweeper

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

var ErrNoOpenSpots = errors.New("no spots are open for this event")
var ErrAlreadyJoined = errors.New("you are already in this event")

type Event struct {
	ID           Snowflake
	Activity     Activity
	Description  string
	Participants []*discordgo.User
}

func (e *Event) IsFull() bool {
	return len(e.Participants) >= e.Activity.MemberCount()
}

func (e *Event) AddParticipant(user *discordgo.User) error {
	if e.IsFull() {
		return ErrNoOpenSpots
	}

	for _, p := range e.Participants {
		if p.ID == user.ID {
			return ErrAlreadyJoined
		}
	}

	e.Participants = append(e.Participants, user)
	return nil
}

func NewEvent(act Activity, desc string) (*Event, error) {
	if desc == "" {
		return nil, errors.New("event may not have an empty description")
	}

	return &Event{
		ID:           NextSnowflake(),
		Activity:     act,
		Description:  desc,
		Participants: make([]*discordgo.User, 0, act.MemberCount()),
	}, nil
}

type EventRepository interface {
	Store(event *Event) error
	FindAll() []*Event
	Find(id Snowflake) (*Event, error)
	Remove(id Snowflake) error
}