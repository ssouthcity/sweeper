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
	r.MessageComponent("plan-leave *", h.LeaveEvent)
	r.MessageComponent("plan-cancel *", h.CancelEvent)
}

func (h PlanningHandler) PlanEvent(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	data := i.ApplicationCommandData()

	act := sweeper.Activity(CommandOption(data.Options, "activity").IntValue())
	usrID := sweeper.Snowflake(i.Member.User.ID)
	desc := CommandOption(data.Options, "description").StringValue()

	id, err := h.planning.PlanEvent(act, usrID, desc)
	if err != nil {
		return nil, err
	}

	evt, err := h.planning.Event(id)
	if err != nil {
		return nil, err
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
							CustomID: fmt.Sprintf("plan-leave %s", id),
							Label:    "Leave",
							Style:    discordgo.SecondaryButton,
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
	}, nil
}

func (h PlanningHandler) JoinEvent(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	data := i.MessageComponentData()

	evtID := sweeper.Snowflake(ComponentParam(data.CustomID, 0))
	usrID := sweeper.Snowflake(i.Member.User.ID)

	if err := h.planning.JoinEvent(evtID, usrID); err != nil {
		return nil, err
	}

	evt, err := h.planning.Event(evtID)
	if err != nil {
		return nil, err
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{h.eventEmbed(evt)},
		},
	}, nil
}

func (h PlanningHandler) LeaveEvent(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	data := i.MessageComponentData()

	eventID := sweeper.Snowflake(ComponentParam(data.CustomID, 0))
	userID := sweeper.Snowflake(i.Member.User.ID)

	if err := h.planning.LeaveEvent(eventID, userID); err != nil {
		return nil, err
	}

	evt, err := h.planning.Event(eventID)
	if err != nil {
		return nil, err
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{h.eventEmbed(evt)},
		},
	}, nil
}

func (h PlanningHandler) CancelEvent(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	data := i.MessageComponentData()

	eventID := sweeper.Snowflake(ComponentParam(data.CustomID, 0))
	userID := sweeper.Snowflake(i.Member.User.ID)

	if err := h.planning.CancelEvent(eventID, userID); err != nil {
		return nil, err
	}

	evt, err := h.planning.Event(eventID)
	if err != nil {
		return nil, err
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{h.eventEmbed(evt)},
			Components: []discordgo.MessageComponent{},
		},
	}, nil
}

func (h PlanningHandler) eventEmbed(event *sweeper.Event) *discordgo.MessageEmbed {
	var participants []string
	for _, p := range event.Participants {
		participants = append(participants, p.Username)
	}

	if len(participants) == 0 {
		participants = append(participants, "no participants yet")
	} else {
		participants[0] = "**" + participants[0] + "**"
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
