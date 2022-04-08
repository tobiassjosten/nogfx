package pkg

import (
	"bufio"
	"io"
	"log"
)

type World struct {
	ui     UI
	stream io.ReadWriter
}

type UI interface {
	Run(<-chan []byte)
}

func NewWorld(ui UI, stream io.ReadWriter) *World {
	return &World{
		ui:     ui,
		stream: stream,
	}
}

func (world *World) Run(inputs <-chan []byte, commands <-chan []byte) {
	uiOutput := make(chan []byte)
	go world.ui.Run(uiOutput)

	serverOutput := make(chan []byte)
	serverDone := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(world.stream)

		for scanner.Scan() {
			serverOutput <- scanner.Bytes()
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		serverDone <- struct{}{}
	}()

	// main:
	for {
		select {
		case _ = <-serverDone:
			uiOutput <- []byte("server disconnected")
			// break main

		case input := <-inputs:
			// @todo Process input.
			world.stream.Write(append(input, '\n'))

		case output := <-serverOutput:
			// @todo Process output.
			uiOutput <- output

		case command, ok := <-commands:
			// @todo Process command.
			if !ok {
				continue
			}
			log.Println(string(command))
		}
	}
}
