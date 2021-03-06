package main

import (
	"context"
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ssouthcity/sweeper"
	"github.com/ssouthcity/sweeper/discord"
	"github.com/ssouthcity/sweeper/flairing"
	"github.com/ssouthcity/sweeper/interaction"
	mgo "github.com/ssouthcity/sweeper/mongo"
	"github.com/ssouthcity/sweeper/planning"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	token = flag.String("token", os.Getenv("DISCORD_TOKEN"), "bot token for discord authentication")
	dburi = flag.String("mongo", os.Getenv("MONGO_HOST"), "uri for the mongodb database")
)

func main() {
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	session, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Err(err).Msg("could not create session")
	}

	mclient, err := mongo.NewClient(options.Client().ApplyURI(*dburi))
	if err != nil {
		log.Err(err).Msg("could not build db client")
	}

	if err := mclient.Connect(context.Background()); err != nil {
		log.Err(err).Msg("could not connect to mongo atlas")
	}
	defer mclient.Disconnect(context.Background())

	db := mclient.Database("sweeper")

	userRepo := discord.NewUserRepository(session, "672544296913600533", map[sweeper.Class]string{sweeper.Titan: "779149217565638708", sweeper.Hunter: "779149035486838795", sweeper.Warlock: "779149290554654730"})
	eventRepo := mgo.NewEventRepository(db.Collection("events"), userRepo)

	planningSrv := planning.NewPlanningService(eventRepo, userRepo)
	flairingSrv := flairing.NewFlairingService(userRepo)

	iHandler := interaction.NewHandler(planningSrv, flairingSrv)

	session.AddHandler(iHandler.HandleInteraction)

	if err := session.Open(); err != nil {
		log.Err(err).Msg("could not connect to discord")
	}
	defer session.Close()

	select {}
}
