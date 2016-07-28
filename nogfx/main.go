package main

import (
	"bufio"
	"github.com/tobiassjosten/nogfx-cli/tui"
	"net"
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

	conn, err := net.Dial("tcp", address())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	serverOutput := make(chan string)
	go func() {
		connbuf := bufio.NewReader(conn)
		for {
			str, err := connbuf.ReadString('\n')
			if err != nil {
				screen.Add(err.Error())
				break
			}
			serverOutput <- str
		}
	}()

	for {
		select {
		case input := <-userInput:
			conn.Write(append(append([]byte(input), '\r'), '\n'))
			screen.Add(" > " + input)
		case output := <-serverOutput:
			screen.Add(output)
		}
	}
}
