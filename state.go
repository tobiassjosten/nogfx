package nogfx

type State struct {
	health int
}

func NewState() *State {
	return &State{}
}
