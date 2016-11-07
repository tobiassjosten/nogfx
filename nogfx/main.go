package main

import (
	"github.com/tobiassjosten/nogfx-cli/tui"
	"os"
)

func address() (address string) {
	host := "achaea.com"
	port := "23"

	switch {
	case len(os.Args) >= 3:
		port = os.Args[2]
		fallthrough
	case len(os.Args) >= 2:
		host = os.Args[1]
	}

	return host + ":" + port
}

func main() {
	userInput := make(chan string)
	screen := tui.NewScreen(userInput)
	go screen.Main()

	telnet, serverOutput := NewTelnet()
	go telnet.Main(address())

	for {
		select {
		case input := <-userInput:
			telnet.Send(input)
			screen.Add(" > " + input)
		case output := <-serverOutput:
			screen.Add(output)
		}
	}
}
