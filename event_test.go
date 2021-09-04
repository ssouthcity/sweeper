package sweeper

import (
	"fmt"
	"testing"
)

func TestNewEvent(t *testing.T) {
	if _, err := NewEvent(Raid, ""); err == nil {
		t.Errorf("expected empty description to raise an error")
	}

	evt, err := NewEvent(Raid, "DSC clan night")
	if err != nil {
		t.Errorf("expected new event to succeed, got error %s", err)
	}

	if a := evt.Activity; a != Raid {
		t.Errorf("expected activity to be 'Raid', got %s", a)
	}

	if s := evt.Description; s != "DSC clan night" {
		t.Errorf("expected description to be 'DSC clan night', got %s", s)
	}

	if l := len(evt.Participants); l != 0 {
		t.Errorf("expected amount of participants to be 0, got %d", l)
	}

	if c := cap(evt.Participants); c != 6 {
		t.Errorf("expected capacity of participants to be 6, got %d", c)
	}
}

func TestParticipants(t *testing.T) {
	evt, err := NewEvent(Trials, "flawless run on reset")
	if err != nil {
		t.Errorf("expected new event to succeed, got error %s", err)
	}

	if evt.IsFull() {
		t.Errorf("expected event to be empty")
	}

	if err := evt.AddParticipant(&User{ID: "0"}); err != nil {
		t.Errorf("expected add participant to succeed, got error %s", err)
	}

	if err := evt.AddParticipant(&User{ID: "0"}); err != ErrAlreadyJoined {
		t.Errorf("expected error to be already joined, got error %s", err)
	}

	for i := 1; i < Trials.MemberCount(); i++ {
		if err := evt.AddParticipant(&User{ID: Snowflake(fmt.Sprintf("%d", i))}); err != nil {
			t.Errorf("expected add participant to succeed, got error %s", err)
		}
	}

	if !evt.IsFull() {
		t.Errorf("expected event to be full")
	}

	if err := evt.AddParticipant(&User{ID: "test"}); err != ErrNoOpenSpots {
		t.Errorf("expected add participant to raise ErrNoOpenSpots, got error %s", err)
	}
}
