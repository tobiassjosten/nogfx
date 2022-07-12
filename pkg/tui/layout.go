package tui

const (
	mainMinWidth = 80
	mainMaxWidth = 120

	borderWidth = 2

	minimapRoomWidth   = 4
	minimapRoomHeight  = 2
	minimapRoomsMargin = 3

	// We want to be able to show at least the current room and those
	// drectly adjacent to it, or else don't bother with the minimap.
	sideMinWidth = minimapRoomWidth*3 + minimapRoomsMargin
	mapMinHeight = minimapRoomHeight*3 + minimapRoomsMargin

	mainSideMinWidth = mainMinWidth + borderWidth + sideMinWidth
	mainSideMaxWidth = mainMinWidth + borderWidth + (mainMaxWidth-mainMinWidth)*2

	paneInput  = "input"
	paneMain   = "main"
	paneMap    = "map"
	paneOutput = "output"
	paneScreen = "screen"
	paneSide   = "side"
	paneTarget = "target"
	paneVitals = "vitals"
)

var paneNames = []string{
	paneInput,
	paneMain,
	paneMap,
	paneOutput,
	paneScreen,
	paneSide,
	paneTarget,
	paneVitals,
}

type pane struct {
	rows   Rows // null == virtual pane (dimensions holder)
	x      int
	y      int
	width  int
	height int
}

func newpane(rows Rows, x, y int) pane {
	if len(rows) == 0 {
		return pane{rows, x, y, 0, 0}
	}

	return pane{rows, x, y, len(rows[0]), len(rows)}
}

// Layout orchestrates all the panes, to determine which one goes where with
// what dimensions.
type Layout struct {
	tui *TUI
}

func (l *Layout) panes() []pane {
	panes := []pane{}
	for _, name := range paneNames {
		panes = append(panes, l.pane(name))
	}

	return panes
}

func (l *Layout) pane(name string) pane {
	width, height := l.tui.screen.Size()

	switch name {
	case paneInput:
		return l.inputPane()

	case paneMain:
		return l.mainPane()

	case paneMap:
		return l.mapPane()

	case paneOutput:
		return l.outputPane()

	case paneScreen:
		return pane{nil, 0, 0, width, height}

	case paneSide:
		return l.sidePane()

	case paneTarget:
		return l.targetPane()

	case paneVitals:
		return l.vitalsPane()
	}

	return pane{}
}

func (l *Layout) inputPane() pane {
	main := l.pane(paneMain)
	target := l.pane(paneTarget)

	vitalsHeight := 0
	if main.height >= 5 {
		vitalsHeight = 1
	}

	// Split half remaining with output, rounding down to yield presedence.
	maxHeight := (main.height - target.height - vitalsHeight) / 2

	rows, cx, cy := l.tui.RenderInput(main.width, maxHeight)

	x := main.x
	y := main.y + main.height - len(rows) - target.height

	l.tui.cursorpos = []int{x + cx, y + cy}

	return newpane(rows, x, y)
}

func (l *Layout) mainPane() pane {
	screen := l.pane(paneScreen)

	width := min(screen.width, mainMaxWidth)
	height := screen.height

	if screen.width > mainSideMinWidth && screen.width < mainSideMaxWidth {
		remainWidth := screen.width - mainMinWidth - borderWidth
		width = mainMinWidth + remainWidth - remainWidth/2
	}

	x := screen.x
	y := screen.y

	return pane{nil, x, y, width, height}
}

func (l *Layout) mapPane() pane {
	side := l.pane(paneSide)

	if side.width == 0 || side.height < mapMinHeight {
		return pane{}
	}

	rows := l.tui.RenderMap(side.width, side.height)

	x := side.x
	y := side.y + borderWidth

	return newpane(rows, x, y)
}

func (l *Layout) outputPane() pane {
	input := l.pane(paneInput)
	main := l.pane(paneMain)
	target := l.pane(paneTarget)
	vitals := l.pane(paneVitals)

	bottomMargin := input.height + target.height + vitals.height
	maxHeight := main.height - bottomMargin

	rows := l.tui.RenderOutput(main.width, maxHeight)

	x := main.x
	y := main.y + main.height - len(rows) - bottomMargin

	return newpane(rows, x, y)
}

func (l *Layout) sidePane() pane {
	screen := l.pane(paneScreen)
	main := l.pane(paneMain)

	if screen.width < mainSideMinWidth {
		return pane{}
	}

	width := screen.width - main.width - borderWidth
	height := screen.height

	x := screen.x + screen.width - width
	y := screen.y

	return pane{nil, x, y, width, height}
}

func (l *Layout) targetPane() pane {
	main := l.pane(paneMain)

	if main.height < 6 {
		return pane{}
	}

	rows := l.tui.RenderTarget(main.width)

	x := main.x
	y := main.y + main.height - len(rows)

	return newpane(rows, x, y)
}

func (l *Layout) vitalsPane() pane {
	input := l.pane(paneInput)
	main := l.pane(paneMain)
	target := l.pane(paneTarget)

	if main.height < 5 {
		return pane{}
	}

	rows := l.tui.RenderVitals(main.width)

	x := main.x
	y := main.y + main.height - len(rows) - input.height - target.height

	return newpane(rows, x, y)
}
