package interaction

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/sweeper"
	"github.com/ssouthcity/sweeper/planning"
)

type PlanningHandler struct {
	planning planning.PlanningService
}

func (h PlanningHandler) interactions(r *Router) {
	r.Handle("plan", h.PlanEvent)
	r.Handle("plan-join", h.JoinEvent)
	r.Handle("plan-cancel", h.CancelEvent)
}

func (h PlanningHandler) PlanEvent(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := i.ApplicationCommandData()

	act := sweeper.Activity(getOption(data.Options, "activity").IntValue())
	desc := getOption(data.Options, "description").StringValue()

	id, err := h.planning.PlanEvent(act, desc)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("error invalid params %s", err),
				Flags:   1 << 6,
			},
		}
	}

	evt, err := h.planning.Event(id)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "event was not successfully created",
				Flags:   1 << 6,
			},
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{h.eventEmbed(evt)},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: fmt.Sprintf("plan-join %s", id),
							Label:    "Join",
							Style:    discordgo.PrimaryButton,
						},
						discordgo.Button{
							CustomID: fmt.Sprintf("plan-cancel %s", id),
							Label:    "Cancel",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		},
	}
}

func (h PlanningHandler) JoinEvent(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := i.MessageComponentData()

	evtID := sweeper.Snowflake(strings.Split(data.CustomID, " ")[1])
	usrID := sweeper.Snowflake(i.Member.User.ID)

	if err := h.planning.JoinEvent(evtID, usrID); err != nil {
		var c string
		if errors.Is(err, sweeper.ErrNoOpenSpots) {
			c = "this event has no more spots"
		} else if errors.Is(err, sweeper.ErrAlreadyJoined) {
			c = "you are already in this event"
		} else {
			c = "could not join event"
		}

		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: c,
				Flags:   1 << 6,
			},
		}
	}

	evt, err := h.planning.Event(evtID)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "event was not successfully created",
				Flags:   1 << 6,
			},
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{h.eventEmbed(evt)},
		},
	}
}

func (h PlanningHandler) CancelEvent(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := i.MessageComponentData()

	id := sweeper.Snowflake(strings.Split(data.CustomID, " ")[1])

	if err := h.planning.CancelEvent(id); err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "unable to cancel event",
				Flags:   1 << 6,
			},
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "event was cancelled",
			Embeds:     []*discordgo.MessageEmbed{},
			Components: []discordgo.MessageComponent{},
		},
	}
}

func (h PlanningHandler) eventEmbed(event *sweeper.Event) *discordgo.MessageEmbed {
	var participants []string
	for _, p := range event.Participants {
		participants = append(participants, p.Username)
	}

	if len(participants) == 0 {
		participants = append(participants, "no participants yet")
	}

	return &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       event.Activity.String(),
		Description: event.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Participants",
				Value: strings.Join(participants, "\n"),
			},
		},
	}
}
