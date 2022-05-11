package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
)

// TUI orchestrates different panes to make up the primary user interface.
type TUI struct {
	screen tcell.Screen

	inputs chan []byte
	input  *Input

	outputs chan []byte
	output  *Output

	vitals map[string]*Vital
	vorder []string

	room    *navigation.Room
	running bool
}

// NewTUI creates a new TUI.
func NewTUI(screen tcell.Screen) *TUI {
	var (
		outputStyle = tcell.Style{}
	)

	tui := &TUI{
		screen: screen,

		inputs: make(chan []byte),
		input:  &Input{},

		outputs: make(chan []byte),
		output:  &Output{},

		vitals: map[string]*Vital{},
	}

	screen.SetStyle(outputStyle)
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)

	return tui
}

// Inputs exposes the outgoing channel for player input.
func (tui *TUI) Inputs() <-chan []byte {
	return tui.inputs
}

// Outputs exposes the incoming channel for server output.
func (tui *TUI) Outputs() chan<- []byte {
	return tui.outputs
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

				if ok := tui.HandleEvent(ev); ok {
					tui.Draw()
				}
			}
		}
	}()

	for {
		select {
		case output := <-tui.outputs:
			tui.output.Append(output)
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

	width, height := tui.screen.Size()

	mainWidth, borderWidth := width, 2
	mainMinWidth := outputMinWidth + borderWidth

	minimapWidth, minimapHeight := 0, height

	// If we can fit a minimap, let's calculate how much space we can
	// afford it. Main panes get at least 80 and at most 120 but in between
	// share the excess with the minimap, before giving it all the rest.
	if width >= mainMinWidth+minimapMinWidth && height >= minimapMinHeight {
		mainWidth = min(
			outputMaxWidth,
			outputMinWidth+(width-mainMinWidth-minimapMinWidth)/2,
		)
		minimapWidth = width - mainWidth - borderWidth
	}

	tui.screen.Clear()

	input := tui.RenderInput(mainWidth)
	tui.paint(0, height-len(input), input)

	// Only give vitals space if there's leftovers from the input.
	vitals := tui.RenderVitals(mainWidth)
	vitals = vitals[0:max(0, min(len(vitals), height-len(input)-1))]
	tui.paint(0, height-len(input)-len(vitals), vitals)

	output := tui.RenderOutput(mainWidth, height-len(input)-len(vitals))
	tui.paint(0, height-len(output)-len(input)-len(vitals), output)

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
