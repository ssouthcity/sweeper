package planning

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/sweeper"
)

type PlanningService interface {
	PlanEvent(a sweeper.Activity, d string) (sweeper.Snowflake, error)
	JoinEvent(id sweeper.Snowflake, user *discordgo.User) error
	Event(id sweeper.Snowflake) (*sweeper.Event, error)
}

type planningService struct {
	events sweeper.EventRepository
}

func NewPlanningService(er sweeper.EventRepository) PlanningService {
	return &planningService{
		events: er,
	}
}

func (s *planningService) PlanEvent(activity sweeper.Activity, description string) (sweeper.Snowflake, error) {
	evt, err := sweeper.NewEvent(activity, description)
	if err != nil {
		return "", err
	}

	if err := s.events.Store(evt); err != nil {
		return "", err
	}

	return evt.ID, nil
}

func (s *planningService) JoinEvent(id sweeper.Snowflake, user *discordgo.User) error {
	evt, err := s.events.Find(id)
	if err != nil {
		return err
	}

	if err := evt.AddParticipant(user); err != nil {
		return err
	}

	if err := s.events.Store(evt); err != nil {
		return err
	}

	return nil
}

func (s *planningService) Event(id sweeper.Snowflake) (*sweeper.Event, error) {
	evt, err := s.events.Find(id)
	if err != nil {
		return nil, err
	}

	return evt, nil
}
