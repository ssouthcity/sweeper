package sweeper

type User struct {
	ID       Snowflake
	Username string
}

type UserRepository interface {
	Find(id Snowflake) (*User, error)
}
