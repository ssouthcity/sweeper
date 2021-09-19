package interaction

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/sweeper"
	"github.com/ssouthcity/sweeper/planning"
)

type PlanningHandler struct {
	planning planning.PlanningService
}

func (h PlanningHandler) interactions(r *InteractionRouter) {
	r.ApplicationCommand("plan", h.PlanEvent)
	r.MessageComponent("plan-join *", h.JoinEvent)
	r.MessageComponent("plan-cancel *", h.CancelEvent)
}

func (h PlanningHandler) PlanEvent(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := i.ApplicationCommandData()

	act := sweeper.Activity(CommandOption(data.Options, "activity").IntValue())
	usrID := sweeper.Snowflake(i.Member.User.ID)
	desc := CommandOption(data.Options, "description").StringValue()

	id, err := h.planning.PlanEvent(act, usrID, desc)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "you can not plan an event at the moment",
				Flags:   1 << 6,
			},
		}
	}

	evt, err := h.planning.Event(id)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "the event you planned no longer exists",
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

	evtID := sweeper.Snowflake(ComponentParam(data.CustomID, 0))
	usrID := sweeper.Snowflake(i.Member.User.ID)

	if err := h.planning.JoinEvent(evtID, usrID); err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "you can not join this event",
				Flags:   1 << 6,
			},
		}
	}

	evt, err := h.planning.Event(evtID)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "the event you are trying to join no longer exists",
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

	eventID := sweeper.Snowflake(ComponentParam(data.CustomID, 0))
	userID := sweeper.Snowflake(i.Member.User.ID)

	if err := h.planning.CancelEvent(eventID, userID); err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "you are not the fireteam leader of this event",
				Flags:   1 << 6,
			},
		}
	}

	evt, err := h.planning.Event(eventID)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "the event you are trying to cancel no longer exists",
				Flags:   1 << 6,
			},
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{h.eventEmbed(evt)},
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

	var color int
	switch event.Activity {
	case sweeper.Raid:
		color = 00000000
	case sweeper.Trials:
		color = 16760576
	}

	var status string
	switch event.Status {
	case sweeper.EventStatusCancelled:
		status = "cancelled"
	case sweeper.EventStatusFull:
		status = "full"
	case sweeper.EventStatusSearching:
		status = "searching"
	}

	return &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Color:       color,
		Title:       event.Activity.String(),
		Description: event.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Participants",
				Value:  strings.Join(participants, "\n"),
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  status,
				Inline: true,
			},
		},
	}
}
