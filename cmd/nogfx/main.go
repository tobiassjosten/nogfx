package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	// "github.com/tobiassjosten/nogfx"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	// "github.com/tobiassjosten/nogfx/pkg/tui"
)

func main() {
	// world := nogfx.NewWorld()

	// ui := tui.NewTUI(world)
	// go ui.Run()

	connection, err := net.Dial("tcp", "achaea.com:23")
	if err != nil {
		log.Fatal(err)
	}

	stream, serverCommands := telnet.NewStream(connection)
	serverOutput := make(chan []byte)
	serverDone := make(chan struct{})

	// err = stream.Do(telnet.GMCP)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go func(stream *telnet.Stream, serverOutput chan<- []byte) {
		scanner := bufio.NewScanner(stream)

		for scanner.Scan() {
			serverOutput <- scanner.Bytes()
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		serverDone <- struct{}{}
	}(stream, serverOutput)

	quit := false

	// Kanske ska det här ligga i world.Run() som main loop? Så att world
	// också lättare kan hooka in, applicera sin egen logik, och exekvera
	// sina reaktioner direkt till serverInChan och tuiInChan (som world
	// alltså måste feedas med).
main:
	for {
		select {
		case command := <-serverCommands:
			fmt.Println("{", command, "}")
		case output := <-serverOutput:
			fmt.Println(`>`, string(output))

			// processa med typ world.Process()
			// skicka vidare (eventuellt förändrat) till tuiInChan
			// case input := <-tuiOutChan:
			// 	// processa med typ world.Process()
			// 	// skicka vidare (eventuellt förändrat) till serverInChan

			if !quit {
				_, err := stream.Write([]byte("3\n"))
				if err != nil {
					log.Fatal(err)
				}
				quit = true
			}
		case _ = <-serverDone:
			fmt.Println("DONE LOOPING")
			break main
		}
	}
}
