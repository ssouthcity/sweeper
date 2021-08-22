package interaction

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/sweeper/planning"
)

type InteractionHandler struct {
	router   *Router
	planning planning.PlanningService
}

func (h InteractionHandler) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p := getPattern(i)

	handler, err := h.router.Get(p)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This slash command is currently not supported",
				Flags:   1 << 6,
			},
		})

		return
	}

	res := handler(s, i)

	if err := s.InteractionRespond(i.Interaction, res); err != nil {
		return
	}
}

func NewHandler(p planning.PlanningService) *InteractionHandler {
	r := NewRouter()

	ph := &PlanningHandler{p}
	ph.interactions(r)

	return &InteractionHandler{
		router:   r,
		planning: p,
	}
}

func getPattern(i *discordgo.InteractionCreate) string {
	var p string

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		p = i.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		p = strings.Split(i.MessageComponentData().CustomID, " ")[0]
	}

	return p
}

func getOption(ops []*discordgo.ApplicationCommandInteractionDataOption, name string) *discordgo.ApplicationCommandInteractionDataOption {
	for _, o := range ops {
		if o.Name == name {
			return o
		}
	}

	return nil
}
