package pkg

// Module is an atomic extension of game logic.
type Module interface {
	Triggers() []Trigger
}
