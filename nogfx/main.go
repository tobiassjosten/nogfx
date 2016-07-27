package main

import (
	"bufio"
	"github.com/tobiassjosten/nogfx-cli/tui"
	"net"
)

func main() {
	userInput := make(chan string)
	screen := tui.NewScreen(userInput)
	go screen.Main()

	conn, err := net.Dial("tcp", "localhost:4000")
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
				break
			}
			serverOutput <- str
		}
	}()

	for {
		select {
		case input := <-userInput:
			conn.Write([]byte(input))
			screen.Add(input)
		case output := <-serverOutput:
			screen.Add(output)
		}
	}
}
