package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/sweeper"
)

type userRepository struct {
	session *discordgo.Session
}

func (r *userRepository) Find(id sweeper.Snowflake) (*sweeper.User, error) {
	u, err := r.session.User(string(id))
	if err != nil {
		return nil, err
	}

	return &sweeper.User{
		ID:       sweeper.Snowflake(u.ID),
		Username: u.Username,
	}, nil
}

func NewUserRepository(session *discordgo.Session) sweeper.UserRepository {
	return &userRepository{session: session}
}
