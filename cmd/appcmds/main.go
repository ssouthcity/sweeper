package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	token = flag.String("token", os.Getenv("DISCORD_TOKEN"), "token for discord authentication")
	guild = flag.String("guild", "", "guild for command syncing, global if omitted")
	spec  = flag.String("cmds", "commands.json", "path to json file containing the app commands spec")
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.Parse()

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Err(err).Msg("unable to create session")
	}

	f, err := os.Open(*spec)
	if err != nil {
		log.Err(err).Msg("unable to open command spec")
	}

	var cmds []*discordgo.ApplicationCommand

	if err := json.NewDecoder(f).Decode(&cmds); err != nil {
		log.Err(err).Msg("unable to unmarshal json")
	}

	if err := s.Open(); err != nil {
		log.Err(err).Msg("unable to establish connection to discord")
	}
	defer s.Close()

	c, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, *guild, cmds)
	if err != nil {
		log.Err(err).Msg("unable to write to application commands")
	}

	names := make([]string, 0, len(c))
	for _, c := range c {
		names = append(names, c.Name)
	}

	log.Info().Strs("cmds", names).Msg("synchronized application commands")
}
