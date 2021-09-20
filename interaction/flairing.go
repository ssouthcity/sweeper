package interaction

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/sweeper"
	"github.com/ssouthcity/sweeper/flairing"
)

type FlairingHandler struct {
	flairing flairing.FlairingService
}

func (h *FlairingHandler) interactions(r *InteractionRouter) {
	r.ApplicationCommand("class", h.classMenu)
	r.MessageComponent("class-pick *", h.classSelect)
}

func (h *FlairingHandler) classMenu(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "choose your class",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:  fmt.Sprintf("class-pick %s", i.Member.User.ID),
							MinValues: 1,
							MaxValues: 1,
							Options: []discordgo.SelectMenuOption{
								{
									Label: "Titan",
									Value: fmt.Sprint(sweeper.Titan),
									Emoji: discordgo.ComponentEmoji{
										Name: "titan",
										ID:   "862064884593328199",
									},
								},
								{
									Label: "Hunter",
									Value: fmt.Sprint(sweeper.Hunter),
									Emoji: discordgo.ComponentEmoji{
										Name: "hunter",
										ID:   "862064884619542538",
									},
								},
								{
									Label: "Warlock",
									Value: fmt.Sprint(sweeper.Warlock),
									Emoji: discordgo.ComponentEmoji{
										Name: "warlock",
										ID:   "862064884702773268",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (h *FlairingHandler) classSelect(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := i.MessageComponentData()

	classID, err := strconv.Atoi(data.Values[0])
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "invalid class option",
				Flags:   1 << 6,
			},
		}
	}

	userID := sweeper.Snowflake(ComponentParam(data.CustomID, 0))
	class := sweeper.Class(classID)

	if i.Member.User.ID != string(userID) {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "you are not the caller of this command",
				Flags:   1 << 6,
			},
		}
	}

	if err := h.flairing.ChangeClass(userID, class); err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "unable to update your class",
				Flags:   1 << 6,
			},
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf("you are now a %s", class.String()),
			Components: []discordgo.MessageComponent{},
		},
	}
}
