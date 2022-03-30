package tui

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/tobiassjosten/nogfx"
)

type TUI struct {
	screen tcell.Screen
	world  *nogfx.World

	input struct {
		buffer []rune
		sent   bool
	}
	// output []struct {
	// 	lines int
	// 	raw   []rune
	// }

	width  int
	height int

	UserInput    chan []rune
	ServerOutput chan []rune
	Quit         chan bool
}

func NewTUI(world *nogfx.World) *TUI {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	screen.SetStyle(style)

	return &TUI{
		screen:       screen,
		world:        world,
		UserInput:    make(chan []rune),
		ServerOutput: make(chan []rune),
	}
}

func (t *TUI) Run() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	err := t.screen.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer t.screen.Fini()

	t.width, t.height = t.screen.Size()

	t.screen.Clear()

	interactions := make(chan tcell.Event)
	go func() {
		for {
			interactions <- t.screen.PollEvent()
		}
	}()

Loop:
	for {
		select {
		case <-t.Quit:
			break Loop

		case <-t.ServerOutput:
			// @todo calculate lines
			// commented out because of type mismatch with t.output
			// case output := <-t.ServerOutput:
			// t.output = append(t.output, output)
			t.draw()

		case interaction := <-interactions:
			switch interaction := interaction.(type) {
			case *tcell.EventKey:
				_ = interaction.Modifiers()

				switch interaction.Key() {
				case tcell.KeyESC, tcell.KeyCtrlC, tcell.KeyCtrlD:
					log.Fatal("ESC")

				case tcell.KeyEnter:
					t.UserInput <- t.input.buffer
					t.input.sent = true
					t.draw()

				case tcell.KeyRune:
					if t.input.sent {
						t.input.buffer = []rune{}
						t.input.sent = false
					}
					t.input.buffer = append(t.input.buffer, interaction.Rune())
					t.draw()
				}

			case *tcell.EventResize:
				t.width, t.height = t.screen.Size()
				t.draw()

			case *tcell.EventInterrupt:
				t.screen.Sync()
			}
		}
	}
}

func (t *TUI) UpdateWorld(world *nogfx.World) {
	t.world = world
	t.draw()
}

func (t *TUI) draw() {
	padding := 1

	minimapWidth, minimapHeight := 13, 7
	minimapx, minimapy := 0, t.height-minimapHeight
	minimapX, minimapY := minimapx+minimapWidth, minimapy+minimapHeight
	t.drawMinimap(minimapx, minimapy, minimapX, minimapY)

	// @todo calculate height
	inputboxWidth, inputboxHeight := t.width-minimapWidth-padding, 1
	inputboxx, inputboxy := minimapX+padding+1, t.height-inputboxHeight
	inputboxX, inputboxY := inputboxx+inputboxWidth, inputboxy+inputboxHeight
	t.drawInputbox(t.input.buffer, inputboxx, inputboxy, inputboxX, inputboxY)

	// @todo calculate measurements
	// commented out because of type mismatch with t.output
	// if len(t.output) > 0 {
	// 	outputboxWidth, outputboxHeight := t.width-minimapWidth-padding, t.height-inputboxHeight-padding
	// 	outputboxx, outputboxy := minimapX+padding+1, 0
	// 	outputboxX, outputboxY := outputboxx+outputboxWidth, outputboxy+outputboxHeight
	// 	t.drawOutputbox(t.output, outputboxx, outputboxy, outputboxX, outputboxY)
	// }

	t.screen.Show()
}
