package pkg

type World interface {
	Input([]byte) []byte
	Output([]byte) []byte
	Command([]byte)
}

type GenericWorld struct {
	client Client
}

func NewGenericWorld(client Client) *GenericWorld {
	return &GenericWorld{
		client: client,
	}
}

func (world *GenericWorld) Input(input []byte) []byte {
	return input
}

func (world *GenericWorld) Output(output []byte) []byte {
	return output
}

func (world *GenericWorld) Command(command []byte) {
}
