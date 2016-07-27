package main

import (
	// "bufio"
	// "fmt"
	"github.com/tobiassjosten/nogfx-cli/tui"
	// "net"
	// "os"
	"strconv"
	"time"
)

func main() {
	userInput := make(chan string)
	screen := tui.NewScreen(userInput)
	go screen.Main()

	// conn, err := net.Dial("tcp", "achaea.com:23")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer conn.Close()

	serverOutput := make(chan string)
	// go func() {
	// 	connbuf := bufio.NewReader(conn)
	// 	for {
	// 		str, err := connbuf.ReadString('\n')
	// 		if err != nil {
	// 			break
	// 		}
	// 		serverOutput <- str
	// 	}
	// }()
	go func() {
		x := 1
		for {
			time.Sleep(1500 * time.Millisecond)
			x++
			serverOutput <- strconv.Itoa(x)
		}
	}()

mainloop:
	for {
		select {
		case input, ok := <-userInput:
			if !ok || "q" == input {
				break mainloop
			}
			// conn.Write([]byte(input))
		case output := <-serverOutput:
			screen.Add(output)
		}
	}
}
