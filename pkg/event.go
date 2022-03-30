package pkg

import nogfx "github.com/tobiassjosten/nogfx"

type Event struct {
	name    string
	mutator func(*nogfx.World) error
}

func (event *Event) Name() string {
	return event.name
}

func (event *Event) Apply(world *nogfx.World) error {
	return event.mutator(world)
}

/*
CHAT_MESSAGE
TELL_MESSAGE
OUTPUT
INPUT
*/
