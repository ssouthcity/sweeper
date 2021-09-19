package sweeper

import (
	"fmt"
	"testing"
)

func TestNewEvent(t *testing.T) {
	user := &User{ID: "test", Username: "test"}

	if _, err := NewEvent(Raid, user, ""); err == nil {
		t.Errorf("expected empty description to raise an error")
	}

	evt, err := NewEvent(Raid, user, "DSC clan night")
	if err != nil {
		t.Errorf("expected new event to succeed, got error %s", err)
	}

	if a := evt.Activity; a != Raid {
		t.Errorf("expected activity to be 'Raid', got %s", a)
	}

	if s := evt.Description; s != "DSC clan night" {
		t.Errorf("expected description to be 'DSC clan night', got %s", s)
	}

	if l := len(evt.Participants); l != 1 {
		t.Errorf("expected amount of participants to be 0, got %d", l)
	}

	if c := cap(evt.Participants); c != 6 {
		t.Errorf("expected capacity of participants to be 6, got %d", c)
	}
}

func TestEventStatus(t *testing.T) {
	user := &User{ID: "0", Username: "Test"}

	evt, err := NewEvent(Trials, user, "flawless run on reset")
	if err != nil {
		t.Errorf("expected new event to succeed, got error %s", err)
	}

	if evt.Status != EventStatusSearching {
		t.Errorf("expected event status to be 'searching' after creation")
	}

	for i := 1; i < Trials.MemberCount(); i++ {
		evt.AddParticipant(&User{ID: Snowflake(fmt.Sprintf("%d", i))})
	}

	if evt.Status != EventStatusFull {
		t.Errorf("expected event status to be 'full' after adding participants")
	}

	evt.Cancel()

	if evt.Status != EventStatusCancelled {
		t.Errorf("expected event status to be 'cancelled' after cancelling event")
	}
}

func TestParticipants(t *testing.T) {
	user := &User{ID: "0", Username: "Test"}

	evt, err := NewEvent(Trials, user, "flawless run on reset")
	if err != nil {
		t.Errorf("expected new event to succeed, got error %s", err)
	}

	if evt.Leader() != user {
		t.Errorf("expected leader to be the creator of the event")
	}

	if err := evt.AddParticipant(&User{ID: "0"}); err != ErrUserAlreadyJoined {
		t.Errorf("expected error to be already joined, got error %s", err)
	}

	if err := evt.RemoveParticipant(&User{ID: "0"}); err != ErrUserIsLeader {
		t.Errorf("expected leader to be unable to leave, got error %s", err)
	}

	for i := 1; i < Trials.MemberCount(); i++ {
		if err := evt.AddParticipant(&User{ID: Snowflake(fmt.Sprintf("%d", i))}); err != nil {
			t.Errorf("expected add participant to succeed, got error %s", err)
		}
	}

	if err := evt.AddParticipant(&User{ID: "fail"}); err != ErrNoOpenSpots {
		t.Errorf("expected add participant to raise ErrNoOpenSpots, got error %s", err)
	}

	if err := evt.RemoveParticipant(&User{ID: "fail"}); err != ErrUserNotParticipant {
		t.Errorf("expected remove participant to fail, got error %s", err)
	}

	if err := evt.RemoveParticipant(&User{ID: "0"}); err != nil {
		t.Errorf("expected remove participant to succeed, got error %s", err)
	}
}

func TestCancel(t *testing.T) {
	user := &User{ID: "0", Username: "Test"}

	evt, err := NewEvent(Raid, user, "vault of glass master")
	if err != nil {
		t.Errorf("expected new event to succeed, got error %s", err)
	}

	evt.Cancel()

	if err := evt.AddParticipant(&User{ID: "1"}); err != ErrEventIsCancelled {
		t.Errorf("expected add participant to fail, got error %s", err)
	}
}
