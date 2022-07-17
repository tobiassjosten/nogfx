package pkg

type Trigger[T any] struct {
	Pattern  []byte
	Callback func(TriggerMatch[T]) T
}

// @todo Kan jag constraina T till input/output?
type TriggerMatch[T any] struct {
	Content  T
	Captures [][]byte
	Index    int
}

// Module is an atomic extension of game logic.
type Module interface {
	InputTriggers() []Trigger[Input]
	OutputTriggers() []Trigger[Output]
}

// ModuleConstructor is a function that constructs a Module.
type ModuleConstructor func(World) Module
