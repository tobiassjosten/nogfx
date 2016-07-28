package tui

import (
	"github.com/nsf/termbox-go"
)

type Screen struct {
	outputBox OutputBox
	inputBox  InputBox
	userInput chan string
	Width     int
	Height    int
}

func NewScreen(userInput chan string) *Screen {
	screen := new(Screen)
	screen.outputBox = OutputBox{screen: screen}
	screen.inputBox = InputBox{screen: screen}
	screen.userInput = userInput
	return screen
}

func (screen *Screen) Add(line string) {
	screen.outputBox.Add(line)
	termbox.Interrupt()
}

func (screen *Screen) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	screen.outputBox.Draw(0, 0, screen.Width-1, screen.Height-2)
	screen.inputBox.Draw(0, screen.Height-1, screen.Width-1, screen.Height-1)
	termbox.Flush()
}

func (screen *Screen) Main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	screen.Width, screen.Height = termbox.Size()

	sentInput := false

mainloop:
	for {
		screen.Draw()
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventResize:
			screen.Width = ev.Width
			screen.Height = ev.Height
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlD:
				close(screen.userInput)
				break mainloop
			case termbox.KeyBackspace:
				screen.inputBox.Remove()
			case termbox.KeyEnter:
				screen.userInput <- screen.inputBox.Get()
				sentInput = true
			default:
				if sentInput {
					screen.inputBox.Clear()
					sentInput = false
				}
				screen.inputBox.Add(ev.Ch)
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
