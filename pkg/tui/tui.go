package tui

import (
	"context"
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
)

// TUI orchestrates different panes to make up the primary user interface.
type TUI struct {
	screen tcell.Screen

	layout *Layout

	cacheMutex sync.Mutex
	panesCache map[string]Rows

	inputs    chan []byte
	input     *Input
	cursorpos []int

	outputs chan []byte
	output  *Output

	character pkg.Character
	room      *navigation.Room
	target    *pkg.Target

	running bool
}

// NewTUI creates a new TUI.
func NewTUI(screen tcell.Screen) *TUI {
	var (
		outputStyle = tcell.Style{}
	)

	tui := &TUI{
		screen: screen,

		panesCache: map[string]Rows{},

		inputs: make(chan []byte),
		input:  &Input{},

		outputs: make(chan []byte),
		output:  &Output{},
	}
	tui.layout = &Layout{tui}

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

func (tui *TUI) setCache(name string, rows Rows) {
	tui.cacheMutex.Lock()
	defer tui.cacheMutex.Unlock()

	if rows == nil {
		delete(tui.panesCache, name)
		return
	}

	tui.panesCache[name] = rows
}

func (tui *TUI) clearCache() {
	tui.cacheMutex.Lock()
	defer tui.cacheMutex.Unlock()

	tui.panesCache = map[string]Rows{}
}

func (tui *TUI) getCache(name string) (Rows, bool) {
	tui.cacheMutex.Lock()
	defer tui.cacheMutex.Unlock()

	rows, ok := tui.panesCache[name]

	return rows, ok
}

// SetCharacter updates the current character and causes a repaint.
func (tui *TUI) SetCharacter(character pkg.Character) {
	tui.character = character
	tui.setCache(paneVitals, nil)
	tui.Draw()
}

// SetRoom updates the current room and causes a repaint.
func (tui *TUI) SetRoom(room *navigation.Room) {
	tui.room = room
	tui.setCache(paneMap, nil)
	tui.Draw()
}

// SetTarget updates the current target and causes a repaint.
func (tui *TUI) SetTarget(target *pkg.Target) {
	tui.target = target
	tui.setCache(paneTarget, nil)
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
		numpad := false

		for {
			event := tui.screen.PollEvent()
			if event == nil {
				return
			}

			switch ev := event.(type) {
			case *tcell.EventResize:
				tui.clearCache()
				tui.Draw()
				tui.screen.Sync()

			case *tcell.EventKey:
				// Numpad keys are handled differently. First
				// there's one event general to all numpad keys
				// and then, immediately following, is another
				// for the specific numpad key.
				if isNumpad(ev) {
					numpad = true
					continue
				} else if numpad {
					numpad = false
					ev = makeNumpad(ev)
				}

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

	tui.Draw()

	for {
		select {
		case output := <-tui.outputs:
			tui.output.Append(output)
			tui.setCache(paneOutput, nil)
			tui.Draw()

		case <-ctx.Done():
			tui.screen.Fini()
			return nil
		}
	}
}

// Draw updates the terminal and prints the contents of the panes.
func (tui *TUI) Draw() {
	if !tui.running {
		return
	}

	for _, p := range tui.layout.panes() {
		tui.paint(p.x, p.y, p.rows)
	}

	if pos := tui.cursorpos; pos != nil {
		tui.screen.ShowCursor(pos[0], pos[1])
	} else {
		tui.screen.HideCursor()
	}

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

const (
	keyNumEnter tcell.Key = iota + 1024
	keyNumEqual
	keyNumMulti
	keyNumPlus
	keyNumMinus
	keyNumDot
	keyNumDiv
	keyNum0
	keyNum1
	keyNum2
	keyNum3
	keyNum4
	keyNum5
	keyNum6
	keyNum7
	keyNum8
	keyNum9
)

var numpadKeys = map[int]tcell.Key{
	77:  keyNumEnter,
	88:  keyNumEqual,
	106: keyNumMulti,
	107: keyNumPlus,
	109: keyNumMinus,
	110: keyNumDot,
	111: keyNumDiv,
	112: keyNum0,
	113: keyNum1,
	114: keyNum2,
	115: keyNum3,
	116: keyNum4,
	117: keyNum5,
	118: keyNum6,
	119: keyNum7,
	120: keyNum8,
	121: keyNum9,
}

func isNumpad(ev *tcell.EventKey) bool {
	return ev.Key() == tcell.KeyRune &&
		ev.Rune() == 'O' &&
		ev.Modifiers() == tcell.ModAlt
}

func makeNumpad(ev *tcell.EventKey) *tcell.EventKey {
	if key, ok := numpadKeys[int(ev.Rune())]; ok {
		return tcell.NewEventKey(key, 0, 0)
	}

	return ev
}
