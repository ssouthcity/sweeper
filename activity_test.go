package sweeper

import "testing"

func TestActivityStringer(t *testing.T) {
	activities := map[Activity]string{
		Raid:   "Raid",
		Trials: "Trials of Osiris",
	}

	for activity, expected := range activities {
		t.Run(expected, func(t *testing.T) {
			if s := activity.String(); s != expected {
				t.Errorf("expected activity to read %s, got %s", expected, s)
			}
		})
	}
}

func TestActivityMemberCount(t *testing.T) {
	activities := map[Activity]int{
		Raid:   6,
		Trials: 3,
	}

	for activity, expected := range activities {
		t.Run(activity.String(), func(t *testing.T) {
			if c := activity.MemberCount(); c != expected {
				t.Errorf("expected member count to be %d, got %d", expected, c)
			}
		})
	}
}
