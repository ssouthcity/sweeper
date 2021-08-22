package sweeper

type Activity int

const (
	Raid Activity = iota + 1
	Trials
)

func (a Activity) String() string {
	var name string

	switch a {
	case Raid:
		name = "Raid"
	case Trials:
		name = "Trials of Osiris"
	}

	return name
}

func (a Activity) MemberCount() int {
	var c int

	switch a {
	case Raid:
		c = 6
	case Trials:
		c = 3
	}

	return c
}
