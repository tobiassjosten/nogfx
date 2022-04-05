package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// "github.com/tobiassjosten/nogfx"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

type MockData struct {
	reader io.Reader
	writer io.Writer
}

func (mock *MockData) Read(p []byte) (int, error) {
	return mock.reader.Read(p)
}

func (mock MockData) Write(p []byte) (int, error) {
	return mock.writer.Write(p)
}

func main() {
	log.SetOutput(ioutil.Discard)

	fileFlags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile("nogfx.log", fileFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	ui, playerInput, err := tui.NewTUI()
	if err != nil {
		log.Fatal(err)
	}

	uiOutput := make(chan []byte)
	go ui.Run(uiOutput)

	// connection, err := net.Dial("tcp", "achaea.com:23")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	connection := &MockData{
		strings.NewReader("trololol\n"),
		&strings.Builder{},
	}

	stream, serverCommands := telnet.NewStream(connection)

	serverOutput := make(chan []byte)
	serverDone := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(stream)

		for scanner.Scan() {
			serverOutput <- scanner.Bytes()
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		serverDone <- struct{}{}
	}()

	// world := nogfx.NewWorld()

	// Kanske ska det här ligga i world.Run() som main loop? Så att world
	// också lättare kan hooka in, applicera sin egen logik, och exekvera
	// sina reaktioner direkt till serverInChan och tuiInChan (som world
	// alltså måste feedas med).
	// main:
	for {
		select {
		case _ = <-serverDone:
			fmt.Println("server done")
			// break main

		case input := <-playerInput:
			stream.Write(input)

		case output := <-serverOutput:
			fmt.Println(`>`, string(output))
			uiOutput <- output

			// processa med typ world.Process()
			// skicka vidare (eventuellt förändrat) till tuiInChan
			// case input := <-tuiOutChan:
			// 	// processa med typ world.Process()
			// 	// skicka vidare (eventuellt förändrat) till serverInChan

		case command, ok := <-serverCommands:
			if !ok {
				continue
			}

			fmt.Println("{", command, "}")
			fmt.Println("{", string(command), "}")
		}
	}
}
