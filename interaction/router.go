package interaction

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type HandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponse

type Router struct {
	interactions map[string]HandlerFunc
}

func (r *Router) Get(key string) (HandlerFunc, error) {
	if h, ok := r.interactions[key]; ok {
		return h, nil
	}

	return nil, fmt.Errorf("interaction %s has no handler", key)
}

func (r *Router) Handle(key string, handler HandlerFunc) error {
	if _, ok := r.interactions[key]; ok {
		return fmt.Errorf("interaction %s is already handled", key)
	}

	r.interactions[key] = handler

	return nil
}

func NewRouter() *Router {
	return &Router{
		interactions: make(map[string]HandlerFunc),
	}
}
