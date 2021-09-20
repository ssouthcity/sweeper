package sweeper

type User struct {
	ID       Snowflake
	Username string
	Class    Class
}

func (u *User) SetClass(c Class) {
	u.Class = c
}

type UserRepository interface {
	Find(id Snowflake) (*User, error)
	Store(u *User) error
}
