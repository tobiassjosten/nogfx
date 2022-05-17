package module

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/simpex"
)

var (
	modLMLIStart    = []byte("learn {^} {^ from *}")
	modLMLOBegin1   = []byte("* begins the lesson in ^.")
	modLMLOBegin2   = []byte("* bows to you and commences the lesson in ^.")
	modLMLOContinue = []byte("* continues your training in ^.")
	modLMLOFinish1  = []byte("* bows to you - the lesson in ^ is over.")
	modLMLOFinish2  = []byte("* finishes the lesson in ^.")
)

// LearnMultipleLessons lets players learn more than 15 lessons in one swoop by
// automatically chaining learning sessions together.
type LearnMultipleLessons struct {
	client pkg.Client
	ui     pkg.UI

	total     int
	remaining int
	target    []byte
	timer     *time.Timer
}

// NewLearnMultipleLessons creates a new LearnMultipleLessons module.
func NewLearnMultipleLessons(client pkg.Client, ui pkg.UI) pkg.Module {
	return &LearnMultipleLessons{
		client: client,
		ui:     ui,
	}
}

// ProcessInput processes player input.
func (mod *LearnMultipleLessons) ProcessInput(input []byte) [][]byte {
	matches := simpex.Match(modLMLIStart, input)
	if matches == nil {
		return [][]byte{}
	}

	number, err := strconv.Atoi(string(matches[0]))
	if err != nil || number <= 15 {
		return [][]byte{}
	}

	mod.total = number
	mod.remaining = number
	mod.target = matches[1]

	newinput := mod.learn()

	return [][]byte{newinput}
}

// ProcessOutput processes server output.
func (mod *LearnMultipleLessons) ProcessOutput(output []byte) [][]byte {
	if mod.total == 0 {
		return [][]byte{}
	}

	switch {
	case simpex.Match(modLMLOBegin1, output) != nil,
		simpex.Match(modLMLOBegin2, output) != nil:
		mod.time()

	case simpex.Match(modLMLOContinue, output) != nil:
		mod.time()

	case simpex.Match(modLMLOFinish1, output) != nil,
		simpex.Match(modLMLOFinish2, output) != nil:
		output = append(output, []byte(
			fmt.Sprintf(" [%d/%d]", mod.total-mod.remaining, mod.total),
		)...)

		if input := mod.learn(); len(input) > 0 {
			mod.client.Send(input)
		}

		return [][]byte{output}
	}

	return [][]byte{}
}

func (mod *LearnMultipleLessons) stop() {
	if mod.timer != nil {
		mod.timer.Stop()
	}

	mod.total = 0
	mod.remaining = 0
	mod.target = []byte{}
	mod.timer = nil
}

func (mod *LearnMultipleLessons) time() {
	if mod.timer != nil {
		mod.timer.Stop()
	}

	mod.timer = time.AfterFunc(10*time.Second, func() {
		mod.stop()
	})
}

func (mod *LearnMultipleLessons) learn() []byte {
	count := 15
	if mod.remaining < count {
		count = mod.remaining
	}

	if count == 0 {
		mod.stop()
		return []byte{}
	}

	mod.remaining -= count

	mod.time()

	return []byte(fmt.Sprintf("learn %d %s", count, mod.target))
}
