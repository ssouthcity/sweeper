package planning

// import (
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/ssouthcity/sweeper"
// 	"github.com/ssouthcity/sweeper/mock"
// )

// func TestPlanEvent(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	events := mock.NewMockEventRepository(ctrl)

// 	events.EXPECT().
// 		Store(gomock.AssignableToTypeOf(&sweeper.Event{})).
// 		Return(nil)

// 	planning := NewPlanningService(events)

// 	_, err := planning.PlanEvent(sweeper.Raid, "DSC")
// 	if err != nil {
// 		t.Errorf("expected planning to succeed, got err %s", err)
// 	}
// }

// func TestJoinEvent(t *testing.T) {
// 	var (
// 		id   sweeper.Snowflake = sweeper.NextSnowflake()
// 		user *sweeper.User     = &sweeper.User{
// 			ID:       "user-mock",
// 			Username: "Carl",
// 		}
// 		event *sweeper.Event = &sweeper.Event{
// 			ID:           id,
// 			Activity:     sweeper.Raid,
// 			Description:  "DSC",
// 			Participants: make([]*sweeper.User, 0, 6),
// 		}
// 	)

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	events := mock.NewMockEventRepository(ctrl)

// 	events.EXPECT().
// 		Find(gomock.Eq(id)).
// 		Return(event, nil)

// 	events.EXPECT().
// 		Store(gomock.Eq(event)).
// 		Return(nil)

// 	planning := NewPlanningService(events)

// 	if err := planning.JoinEvent(id, user); err != nil {
// 		t.Errorf("expected join event to succeed, got error %s", err)
// 	}
// }
