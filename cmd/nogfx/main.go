package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/telnet"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

func main() {
	log.SetOutput(ioutil.Discard)

	fileFlags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile("nogfx.log", fileFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	ui, inputs, err := tui.NewTUI()
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.Dial("tcp", "achaea.com:23")
	if err != nil {
		log.Fatal(err)
	}

	client, commands := telnet.NewClient(connection)
	client.AcceptWill(telnet.GMCP)

	world := pkg.NewWorld(ui, client)
	world.Run(inputs, commands)
}
