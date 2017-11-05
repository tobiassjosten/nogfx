package nogfx

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * STRUCTS AND CONSTRUCTORS
 */

type Message struct {
	text  string
	lines int
}

type Screen struct {
	userInput  chan string
	userBuffer string
	userSent   bool
	Width      int
	Height     int
	events     []*Event
	messages   []*Message
}

func NewScreen() (*Screen, <-chan string) {
	userInput := make(chan string)
	return &Screen{userInput: userInput}, userInput
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * PUBLIC METHODS
 */

func (screen *Screen) SetEvents(events []*Event) {
	screen.events = events
	screen.ProcessEvents()
}

func (screen *Screen) ProcessEvents() {
	minX, maxX := 0, screen.Width-1

	screen.messages = []*Message{}

	color := false
	for _, event := range screen.events {
		x := minX

		message := &Message{text: event.text, lines: 1}

		for _, character := range message.text {
			if '\x1b' == character {
				color = true
			} else if color && 'm' == character {
				color = false
			} else if !color {
				if '\r' == character {
					continue
				} else if '\n' == character || x > maxX {
					message.lines++
					x = minX
				} else {
					x++
				}
			}
		}

		screen.messages = append(screen.messages, message)
	}
}

func (screen *Screen) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	screen.drawOutput()
	screen.drawInput()
	termbox.Flush()
}

func (screen *Screen) Main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	screen.Width, screen.Height = termbox.Size()

mainloop:
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventResize:
			screen.Width = ev.Width
			screen.Height = ev.Height
			screen.ProcessEvents()
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlD:
				close(screen.userInput)
				break mainloop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if screen.userSent {
					screen.userBuffer = ""
					screen.userSent = false
				} else if 0 < len(screen.userBuffer) {
					screen.userBuffer = screen.userBuffer[0 : len(screen.userBuffer)-1]
				}
			case termbox.KeyEnter:
				screen.userSent = true
				screen.userInput <- screen.userBuffer
			case termbox.KeySpace:
				ev.Ch = ' '
				fallthrough
			default:
				if screen.userSent {
					screen.userBuffer = ""
					screen.userSent = false
				}
				screen.userBuffer += string(ev.Ch)
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		screen.Draw()
	}
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * OUTPUT BOX
 */

// ColorDefault ColorBlack ColorRed ColorGreen ColorYellow ColorBlue ColorMagenta ColorCyan ColorWhite
// AttrBold AttrUnderline AttrReverse

// ansiColors := map[string]termbox.Attribute{
// 	"30": termbox.ColorDefault,
// }

// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
// "\x1b[38;5;007m" === "\x1b[38;5;7m"
// 38;5;x = foreground x 0-255 color, 48;5;x = background x 0-255 color
// 0 = reset all, 22 = reset color and intensity, 24 = reset underline
// 4 = underline
// 3x = foreground, 4x = background: 0 = black, 1 = red, 2 = green, 3 = yellow, 4 = blue, 5 = magenta, 6 = cyan, 7 = white

func (screen *Screen) drawOutput() {
	minX, minY, maxX, maxY := 0, 0, screen.Width-1, screen.Height-2

	var messages []*Message
	var totalLines int

	// Determine which events' messages will fit and extract them to `messages` slice.
	for i := len(screen.messages) - 1; i >= 0 && totalLines <= maxY-minY; i-- {
		messages = append([]*Message{screen.messages[i]}, messages...)
		totalLines += screen.messages[i].lines
	}

	// Calculate starting points to (potentially, if minY permits) begin printing lines.
	y := minY - (totalLines - maxY - minY + 1)

	color := false
	for _, message := range messages {
		x := minX

		fgcolor, bgcolor := termbox.ColorDefault, termbox.ColorDefault

		for _, character := range message.text {
			if '\x1b' == character {
				color = true
			} else if color && 'm' == character {
				color = false
			} else if !color {
				if '\r' == character {
					continue
				}

				if '\n' == character {
					// @todo Gör detsamma för prompt-avslutaren
					y++
					x = 0
					continue
				}

				if x > maxX {
					y++
					x = 0
				}

				if y >= minY {
					termbox.SetCell(x, y, character, fgcolor, bgcolor)
				}
				x++
			}
		}

		y++
	}
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 * INPUT BOX
 */

func (screen *Screen) drawInput() {
	x, y, X := 0, screen.Height-1, screen.Width-1

	fg := termbox.ColorWhite
	if screen.userSent {
		fg = termbox.ColorCyan
	}

	for _, c := range []rune(screen.userBuffer) {
		termbox.SetCell(x, y, c, fg, termbox.ColorDefault)
		x += runewidth.RuneWidth(c)
	}

	termbox.SetCursor(x, y)

	for ; x < X; x++ {
		termbox.SetCell(x, y, '_', termbox.ColorWhite, termbox.ColorDefault)
	}
}
