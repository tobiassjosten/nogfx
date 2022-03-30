package main

import (
	"bufio"
	"log"

	"github.com/tobiassjosten/nogfx"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

func main() {
	world := nogfx.NewWorld()

	ui := tui.NewTUI(world)
	go ui.Run()

	stream, err := telnet.Dial("tcp", "achaea.com:23")
	if err != nil {
		log.Fatal(err)
	}

	// err = stream.Do(telnet.GMCP)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	serverOutput := make(chan []byte)
	go func(serverOutput chan []byte) {
		scanner := bufio.NewScanner(stream)

		for scanner.Scan() {
			serverOutput <- scanner.Bytes()
		}

		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
	}(serverOutput)

	for {
		select {
		case output := <-serverOutput:
			ui.ServerOutput <- []rune(string(output))

		case input := <-ui.UserInput:
			_, err = stream.Write([]byte(string(input) + "\n"))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
