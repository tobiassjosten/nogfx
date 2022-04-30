package achaea

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/simpex"
)

var (
	modLMLIStart    = []byte("learn {^} {^ from *}")
	modLMLOBegin    = []byte("* bows to you and commences the lesson in ^.")
	modLMLOContinue = []byte("* continues your training in ^.")
	modLMLOFinish   = []byte("* bows to you - the lesson in ^ is over.")
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
func (mod *LearnMultipleLessons) ProcessInput(input []byte) []byte {
	matches := simpex.Match(modLMLIStart, input)
	if matches == nil {
		return input
	}

	number, err := strconv.Atoi(string(matches[0]))
	if err != nil || number <= 15 {
		return input
	}

	mod.total = number
	mod.remaining = number
	mod.target = matches[1]

	mod.learn()

	return nil
}

// ProcessOutput processes server output.
func (mod *LearnMultipleLessons) ProcessOutput(output []byte) []byte {
	if mod.total == 0 {
		return output
	}

	switch {
	case simpex.Match(modLMLOBegin, output) != nil:
		mod.time()

	case simpex.Match(modLMLOContinue, output) != nil:
		mod.time()

	case simpex.Match(modLMLOFinish, output) != nil:
		output = append(output, []byte(fmt.Sprintf(" [%d/%d]", mod.total-mod.remaining, mod.total))...)
		mod.learn()
	}

	return output
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

func (mod *LearnMultipleLessons) learn() {
	count := 15
	if mod.remaining < count {
		count = mod.remaining
	}

	if count == 0 {
		mod.stop()
		return
	}

	// @todo Replace with Client.Send() so as to let Client handle errors.
	_, err := mod.client.Write([]byte(fmt.Sprintf(
		"learn %d %s\n", count, mod.target,
	)))
	if err != nil {
		mod.ui.Print([]byte(fmt.Sprintf("[Failed learning! %s]", err)))
	}

	mod.remaining -= count

	mod.time()
}
