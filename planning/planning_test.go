package planning

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ssouthcity/sweeper"
	"github.com/ssouthcity/sweeper/mock"
)

func TestPlanEvent(t *testing.T) {
	var (
		activity = sweeper.Raid
		userID   = sweeper.NextSnowflake()
		desc     = "Last Wish Flawless"
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := mock.NewMockEventRepository(ctrl)

	events.EXPECT().
		Store(gomock.AssignableToTypeOf(&sweeper.Event{})).
		Return(nil)

	users := mock.NewMockUserRepository(ctrl)

	users.EXPECT().
		Find(userID).
		Return(&sweeper.User{ID: userID, Username: "Test"}, nil)

	planning := NewPlanningService(events, users)

	if _, err := planning.PlanEvent(activity, userID, desc); err != nil {
		t.Errorf("expected planning to succeed, got err %s", err)
	}
}

func TestJoinEvent(t *testing.T) {
	var (
		eventID = sweeper.NextSnowflake()
		userID  = sweeper.NextSnowflake()
		user    = &sweeper.User{
			ID:       userID,
			Username: "Carl",
		}
		event = &sweeper.Event{
			ID:           eventID,
			Activity:     sweeper.Raid,
			Description:  "DSC",
			Participants: make([]*sweeper.User, 0, 6),
		}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := mock.NewMockEventRepository(ctrl)

	events.EXPECT().
		Find(gomock.Eq(eventID)).
		Return(event, nil)

	events.EXPECT().
		Store(gomock.Eq(event)).
		Return(nil)

	users := mock.NewMockUserRepository(ctrl)

	users.EXPECT().
		Find(userID).
		Return(user, nil)

	planning := NewPlanningService(events, users)

	if err := planning.JoinEvent(eventID, userID); err != nil {
		t.Errorf("expected join event to succeed, got error %s", err)
	}
}

func TestEvent(t *testing.T) {
	var (
		eventID = sweeper.NextSnowflake()
		event   = &sweeper.Event{
			ID: eventID,
		}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := mock.NewMockEventRepository(ctrl)
	users := mock.NewMockUserRepository(ctrl)

	events.EXPECT().
		Find(eventID).
		Return(event, nil)

	planning := NewPlanningService(events, users)

	if _, err := planning.Event(eventID); err != nil {
		t.Errorf("expected event to be found, got error %s", err)
	}
}

func TestEvents(t *testing.T) {
	var (
		eventID  = sweeper.NextSnowflake()
		event    = &sweeper.Event{ID: eventID}
		eventMap = map[sweeper.Snowflake]*sweeper.Event{eventID: event}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := mock.NewMockEventRepository(ctrl)
	users := mock.NewMockUserRepository(ctrl)

	events.EXPECT().
		FindAll().
		Return(eventMap)

	planning := NewPlanningService(events, users)

	evts := planning.Events()

	for i, evt := range evts {
		if eventMap[i].ID != evt.ID {
			t.Errorf("expected events to be equal")
		}
	}
}

func TestCancelEvent(t *testing.T) {
	var (
		eventID = sweeper.NextSnowflake()
		userID  = sweeper.NextSnowflake()

		event = &sweeper.Event{
			ID: eventID,
			Participants: []*sweeper.User{
				{ID: userID, Username: "The Stranger"},
			},
		}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := mock.NewMockEventRepository(ctrl)
	users := mock.NewMockUserRepository(ctrl)

	events.EXPECT().
		Find(eventID).
		Return(event, nil)

	events.EXPECT().
		Store(event).
		Return(nil)

	planning := NewPlanningService(events, users)

	if err := planning.CancelEvent(eventID, userID); err != nil {
		t.Errorf("expected event to be cancelled, got err %s", err)
	}
}
