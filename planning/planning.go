package planning

import (
	"github.com/ssouthcity/sweeper"
)

type PlanningService interface {
	PlanEvent(a sweeper.Activity, d string) (sweeper.Snowflake, error)
	JoinEvent(id sweeper.Snowflake, userID sweeper.Snowflake) error
	Event(id sweeper.Snowflake) (*sweeper.Event, error)
	CancelEvent(id sweeper.Snowflake) error
}

type planningService struct {
	events sweeper.EventRepository
	users  sweeper.UserRepository
}

func NewPlanningService(er sweeper.EventRepository, ur sweeper.UserRepository) PlanningService {
	return &planningService{
		events: er,
		users:  ur,
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

func (s *planningService) JoinEvent(id sweeper.Snowflake, userID sweeper.Snowflake) error {
	evt, err := s.events.Find(id)
	if err != nil {
		return err
	}

	usr, err := s.users.Find(userID)
	if err != nil {
		return err
	}

	if err := evt.AddParticipant(usr); err != nil {
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

func (s *planningService) CancelEvent(id sweeper.Snowflake) error {
	return s.events.Remove(id)
}
