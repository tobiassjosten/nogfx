package pkg

// World represents a game and hooks into all their various specific logic.
type World interface {
	Input([]byte) []byte
	Output([]byte) []byte
	Command([]byte) error
}

// GenericWorld is used for unknown games and has no specific logic of its own.
type GenericWorld struct {
	ui     UI
	client Client
}

// NewGenericWorld creates a GenericWorld.
func NewGenericWorld(ui UI, client Client) *GenericWorld {
	return &GenericWorld{
		ui:     ui,
		client: client,
	}
}

// Input receives a player input and does nothing.
func (world *GenericWorld) Input(input []byte) []byte {
	return input
}

// Output receives server output and does nothing.
func (world *GenericWorld) Output(output []byte) []byte {
	return output
}

// Command receives a telnet command sequence from the server  and does nothing.
func (world *GenericWorld) Command(command []byte) error {
	return nil
}
