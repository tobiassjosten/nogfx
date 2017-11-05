package nogfx

import (
	"github.com/tobiassjosten/nogfx-cli/network"
	"strings"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * STRUCTS AND CONSTRUCTORS
 */

type Event struct {
	text      string
	userInput bool
}

type Engine struct {
	screen       *Screen
	telnet       *network.Telnet
	state        *State
	userInput    <-chan string
	serverOutput <-chan string
	events       []*Event
}

func NewEngine() *Engine {
	screen, userInput := NewScreen()
	telnet, serverOutput := network.NewTelnet()

	return &Engine{
		screen:       screen,
		telnet:       telnet,
		state:        NewState(),
		userInput:    userInput,
		serverOutput: serverOutput,
	}
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * PUBLIC METHODS
 */

func (engine *Engine) Run(address string) {
	go engine.screen.Main()
	go engine.telnet.Main(address)

	for {
		select {
		case input := <-engine.userInput:
			engine.UserInput(input)
		case output := <-engine.serverOutput:
			engine.ServerOutput(output)
		}
	}
}

func (engine *Engine) UserInput(input string) {
	engine.events = append(engine.events, &Event{text: input, userInput: true})
	engine.telnet.Send(input)
}

func (engine *Engine) ServerOutput(output string) {
	// @todo Load up a buffer and only dispatch at prompts. Remember: this can
	// be multiple messages. Break up prompts into their own messages.
	engine.events = append(engine.events, &Event{text: strings.Trim(output, "\r\n"), userInput: false})
	engine.screen.SetEvents(engine.events)
	engine.screen.Draw()
}

// events: prompt (+blackout?)
// fallbacks: user input, server output
