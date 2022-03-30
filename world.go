package nogfx

type World struct {
	Character  *Character
	Characters []*Character
	Rooms      []*Room
}

func NewWorld() *World {
	character := &Character{}

	return &World{
		Character:  character,
		Characters: []*Character{character},
	}
}
