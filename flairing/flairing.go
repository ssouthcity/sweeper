package flairing

import "github.com/ssouthcity/sweeper"

type FlairingService interface {
	ChangeClass(userID sweeper.Snowflake, class sweeper.Class) error
}

type flairingService struct {
	users sweeper.UserRepository
}

func (s *flairingService) ChangeClass(userID sweeper.Snowflake, class sweeper.Class) error {
	u, err := s.users.Find(userID)
	if err != nil {
		return err
	}

	u.SetClass(class)

	if err := s.users.Store(u); err != nil {
		return err
	}

	return nil
}

func NewFlairingService(users sweeper.UserRepository) FlairingService {
	return &flairingService{users}
}
