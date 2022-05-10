package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
)

// Panes is a collection of various panes used throughout the user interface.
type Panes struct {
	Input  *InputPane
	Output *OutputPane
	Vitals *VitalsPane
}

// NewPanes creates a new Panes with the standard collection of panes.
func NewPanes() Panes {
	return Panes{
		Input:  NewInputPane(),
		Output: NewOutputPane(),
		Vitals: NewVitalsPane(),
	}
}

// TUI orchestrates different panes to make up the primary user interface.
type TUI struct {
	screen  tcell.Screen
	panes   Panes
	inputs  chan []byte
	room    *navigation.Room
	running bool
}

// NewTUI creates a new TUI.
func NewTUI(screen tcell.Screen, panes Panes) *TUI {
	var (
		outputStyle = tcell.Style{}
	)

	tui := &TUI{
		screen: screen,
		inputs: make(chan []byte),
		panes:  panes,
	}

	screen.SetStyle(outputStyle)
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)

	return tui
}

// Inputs exposes the outgoing channel for player input.
func (tui *TUI) Inputs() <-chan []byte {
	return tui.inputs
}

// SetRoom updates the current room and causes a repaint.
func (tui *TUI) SetRoom(room *navigation.Room) {
	tui.room = room
	tui.Draw()
}

// Run is the main loop of the user interface, where everything is orchestrated.
func (tui *TUI) Run(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)

	tui.running = true
	defer func() { tui.running = false }()

	if err := tui.screen.Init(); err != nil {
		cancel()
		return fmt.Errorf("failed initializing screen: %w", err)
	}

	go func() {
		for {
			event := tui.screen.PollEvent()
			if event == nil {
				return
			}

			switch ev := event.(type) {
			case *tcell.EventResize:
				tui.Draw()
				tui.screen.Sync()

			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlD {
					cancel()
					return
				}

				if ok := tui.panes.Output.HandleEvent(ev); ok {
					tui.Draw()
				}

				if ok, input := tui.panes.Input.HandleEvent(ev); ok {
					if input != nil {
						tui.inputs <- []byte(string(input))
					}
					tui.Draw()
				}
			}
		}
	}()

	for {
		select {
		case output := <-tui.panes.Output.outputs:
			tui.panes.Output.Add(output)
			tui.Draw()

		case <-ctx.Done():
			tui.screen.Fini()
			return nil
		}
	}
}

// Resize calculates the layout of the various panes.
func (tui *TUI) Resize(width, height int) {
}

// Draw updates the terminal and prints the contents of the panes.
func (tui *TUI) Draw() {
	if !tui.running {
		return
	}

	// @todo Cache renditions so as to not redraw everything every time.

	tui.screen.Clear()

	width, height := tui.screen.Size()

	mainMinWidth := 80
	mainWidth := width

	borderWidth := 2

	roomWidth, _, roomsMargin := 4, 2, 3
	minimapMinWidth := roomWidth*3 + roomsMargin
	minimapWidth, minimapHeight := 0, 0

	// If we can fit a minimap, let's calculate how much space we can
	// afford it. Main panes get at least 80 and at most 120 but in between
	// share the excess with the minimap, before giving it all the rest.
	if width >= mainMinWidth+borderWidth+minimapMinWidth {
		mainWidth = min(120, mainMinWidth+(width-mainMinWidth-borderWidth-minimapMinWidth)/2)

		minimapWidth = width - mainWidth - borderWidth
		minimapHeight = height
	}

	inputWidth := mainWidth
	inputHeight := min(height, tui.panes.Input.Height())
	inputX, inputY := 0, height-inputHeight
	tui.panes.Input.Position(inputX, inputY, inputWidth, inputHeight)

	// Draw VitalsPane if OutputPane has at least one row.
	vitalsHeight := min(tui.panes.Vitals.Height(), max(0, height-inputHeight-1))
	tui.panes.Vitals.Position(inputX, inputY-1, mainWidth, vitalsHeight)

	tui.panes.Input.Draw(tui.screen)
	tui.panes.Vitals.Draw(tui.screen)

	outputWidth := mainWidth
	outputHeight := height - inputHeight - vitalsHeight

	tui.paint(0, 0, tui.panes.Output.Render(outputWidth, outputHeight))

	tui.paint(
		mainWidth+borderWidth, height-minimapHeight,
		RenderMap(tui.room, minimapWidth, minimapHeight),
	)

	tui.screen.Show()
}

func (tui *TUI) paint(x, y int, rows Rows) {
	for yy, row := range rows {
		for xx, cell := range row {
			tui.screen.SetContent(
				x+xx, y+yy,
				cell.Content, nil, cell.Style,
			)
		}
	}
}
