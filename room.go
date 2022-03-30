package nogfx

type Room struct {
	Exits map[string]Room
}
