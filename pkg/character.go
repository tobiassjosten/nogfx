package pkg

// CharacterVital is a measurement of some vital aspect, like health or mana.
type CharacterVital struct {
	Value int
	Max   int
}

// Character represents the persona being played.
type Character struct {
	Vitals map[string]CharacterVital
}
