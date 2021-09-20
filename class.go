package sweeper

type Class int

const (
	Titan = iota + 1
	Hunter
	Warlock
)

func (c Class) String() string {
	var s string

	switch c {
	case Titan:
		s = "Titan"
	case Hunter:
		s = "Hunter"
	case Warlock:
		s = "Warlock"
	}

	return s
}
