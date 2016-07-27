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
	screen.Width, screen.Height = termbox.Size()
	return screen
}

func (screen *Screen) Add(line string) {
	screen.outputBox.Add(line)
	termbox.Interrupt()
}

func (screen *Screen) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	screen.outputBox.Draw()
	screen.inputBox.Draw()
	termbox.Flush()
}

func (screen *Screen) Main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

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
			default:
				screen.inputBox.Add(ev.Ch)
				if ev.Key == termbox.KeyEnter {
					screen.userInput <- screen.inputBox.Get()
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
