package pkg

type CharacterVital struct {
	Value int
	Max   int
}

type Character struct {
	Vitals map[string]CharacterVital
}

type Target struct {
	Name       string
	Health     int
	Candidates []string
}
