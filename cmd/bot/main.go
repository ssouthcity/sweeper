package main

import (
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ssouthcity/sweeper"
	"github.com/ssouthcity/sweeper/inmem"
	"github.com/ssouthcity/sweeper/interaction"
	"github.com/ssouthcity/sweeper/planning"
)

var (
	token = flag.String("token", os.Getenv("DISCORD_TOKEN"), "bot token for discord authentication")
)

func main() {
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var er sweeper.EventRepository = inmem.NewEventRepository()

	var ps planning.PlanningService = planning.NewPlanningService(er)

	session, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Err(err).Msg("could not create session")
	}

	ih := interaction.NewHandler(ps)

	session.AddHandler(ih.OnInteractionCreate)

	if err := session.Open(); err != nil {
		log.Err(err).Msg("could not connect to discord")
	}
	defer session.Close()

	select {}
}
