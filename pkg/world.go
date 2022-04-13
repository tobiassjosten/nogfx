package pkg

type World interface {
	Input([]byte) []byte
	Output([]byte) []byte
	Command([]byte) error
}

type GenericWorld struct {
	ui     UI
	client Client
}

func NewGenericWorld(ui UI, client Client) *GenericWorld {
	return &GenericWorld{
		ui:     ui,
		client: client,
	}
}

func (world *GenericWorld) Input(input []byte) []byte {
	return input
}

func (world *GenericWorld) Output(output []byte) []byte {
	return output
}

func (world *GenericWorld) Command(command []byte) error {
	return nil
}
