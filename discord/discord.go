package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/sweeper"
)

type userRepository struct {
	session      *discordgo.Session
	guildID      string
	classRoleIDs map[sweeper.Class]string
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

func (r *userRepository) Store(u *sweeper.User) error {
	for _, rid := range r.classRoleIDs {
		if err := r.session.GuildMemberRoleRemove(r.guildID, string(u.ID), rid); err != nil {
			return err
		}
	}

	if err := r.session.GuildMemberRoleAdd(r.guildID, string(u.ID), r.classRoleIDs[u.Class]); err != nil {
		return err
	}

	return nil
}

func NewUserRepository(session *discordgo.Session, guildID string, classRoleIDs map[sweeper.Class]string) sweeper.UserRepository {
	return &userRepository{session: session, guildID: guildID, classRoleIDs: classRoleIDs}
}
