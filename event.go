package nogfx

type Event interface {
	Name() string
	Apply(World)
}
