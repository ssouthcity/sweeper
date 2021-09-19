package planning

import (
	"errors"

	"github.com/ssouthcity/sweeper"
)

type PlanningService interface {
	PlanEvent(a sweeper.Activity, l sweeper.Snowflake, d string) (sweeper.Snowflake, error)
	JoinEvent(id sweeper.Snowflake, userID sweeper.Snowflake) error
	Event(id sweeper.Snowflake) (*sweeper.Event, error)
	Events() map[sweeper.Snowflake]*sweeper.Event
	CancelEvent(id sweeper.Snowflake, userID sweeper.Snowflake) error
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

func (s *planningService) PlanEvent(activity sweeper.Activity, leaderID sweeper.Snowflake, description string) (sweeper.Snowflake, error) {
	leader, err := s.users.Find(leaderID)
	if err != nil {
		return "", err
	}

	evt, err := sweeper.NewEvent(activity, leader, description)
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

func (s *planningService) Events() map[sweeper.Snowflake]*sweeper.Event {
	evts := s.events.FindAll()

	return evts
}

func (s *planningService) CancelEvent(id sweeper.Snowflake, userID sweeper.Snowflake) error {
	evt, err := s.events.Find(id)
	if err != nil {
		return err
	}

	if evt.Leader().ID != userID {
		return errors.New("user is not permitted to do this action")
	}

	evt.Cancel()

	return s.events.Store(evt)
}
