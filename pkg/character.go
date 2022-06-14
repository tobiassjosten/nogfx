package pkg

type CharacterVital struct {
	Value int
	Max   int
}

type Character struct {
	Vitals map[string]CharacterVital
}
