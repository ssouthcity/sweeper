package interaction

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/sweeper/flairing"
	"github.com/ssouthcity/sweeper/planning"
)

var ErrAlreadyRegistered = errors.New("interaction is already registered")

type InteractionHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse

type InteractionRouter struct {
	applicationCommands map[string]InteractionHandler
	messageComponents   map[string]InteractionHandler
	planning            planning.PlanningService
}

func (r *InteractionRouter) ApplicationCommand(name string, handler InteractionHandler) error {
	if _, ok := r.applicationCommands[name]; ok {
		return ErrAlreadyRegistered
	}
	r.applicationCommands[name] = handler
	return nil
}

func (r *InteractionRouter) MessageComponent(pattern string, handler InteractionHandler) error {
	if _, ok := r.messageComponents[pattern]; ok {
		return ErrAlreadyRegistered
	}
	r.messageComponents[pattern] = handler
	return nil
}

func (r *InteractionRouter) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var res *discordgo.InteractionResponse
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		res = r.handleApplicationCommand(s, i)
	case discordgo.InteractionMessageComponent:
		res = r.handleMessageComponent(s, i)
	}

	s.InteractionRespond(i.Interaction, res)
}

func (r *InteractionRouter) handleApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := i.ApplicationCommandData()

	if h, ok := r.applicationCommands[data.Name]; ok {
		return h(s, i)
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "this command seems to be unavailable at the moment",
			Flags:   1 << 6,
		},
	}
}

func (r *InteractionRouter) handleMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := i.MessageComponentData()

	for key, val := range r.messageComponents {
		if matched, _ := regexp.MatchString(key, data.CustomID); matched {
			return val(s, i)
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "this action can not be made at the moment",
			Flags:   1 << 6,
		},
	}
}

func NewHandler(p planning.PlanningService, f flairing.FlairingService) *InteractionRouter {
	r := &InteractionRouter{
		applicationCommands: make(map[string]InteractionHandler),
		messageComponents:   make(map[string]InteractionHandler),
		planning:            p,
	}

	ph := &PlanningHandler{p}
	ph.interactions(r)

	fh := &FlairingHandler{f}
	fh.interactions(r)

	return r
}

func CommandOption(ops []*discordgo.ApplicationCommandInteractionDataOption, name string) *discordgo.ApplicationCommandInteractionDataOption {
	for _, o := range ops {
		if o.Name == name {
			return o
		}
	}

	return nil
}

func ComponentParam(customID string, param int) string {
	return strings.Split(customID, " ")[param+1]
}
