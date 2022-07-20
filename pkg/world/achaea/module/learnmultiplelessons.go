package module

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/tobiassjosten/nogfx/pkg"
)

// @todo Make this use 20 lessons at a time with the myrrh/bisemutum defense:
// "Your mind is racing with enhanced speed."
var maxLessons int = 15

// LearnMultipleLessons lets players learn an unlimited amount of lessons in
// one swoop by automatically chaining learning sessions together.
type LearnMultipleLessons struct {
	total     int
	remaining int
	target    []byte
	start     time.Time
	timer     *time.Timer
}

// NewLearnMultipleLessons creates a new LearnMultipleLessons module.
func NewLearnMultipleLessons() pkg.Module {
	return &LearnMultipleLessons{}
}

// Triggers returns a list of triggers.
func (mod *LearnMultipleLessons) Triggers() []pkg.Trigger {
	whenActive := func(callback pkg.Callback) pkg.Callback {
		if mod.timer == nil {
			return pkg.NoopCallback
		}
		return callback
	}

	return []pkg.Trigger{
		{
			Kind:     pkg.Input,
			Pattern:  []byte("learn {^} {^ from *}"),
			Callback: mod.onStart,
		},
		{
			Kind:     pkg.Output,
			Pattern:  []byte("* begins the lesson in ^."),
			Callback: whenActive(mod.onBegin),
		},
		{
			Kind:     pkg.Output,
			Pattern:  []byte("* bows to you and commences the lesson in ^."),
			Callback: whenActive(mod.onBegin),
		},
		{
			Kind:     pkg.Output,
			Pattern:  []byte("* continues your training in ^."),
			Callback: whenActive(mod.onUpdate),
		},
		{
			Kind:     pkg.Output,
			Pattern:  []byte("* finishes the lesson in ^."),
			Callback: whenActive(mod.onFinish),
		},
		{
			Kind:     pkg.Output,
			Pattern:  []byte("Storing ^ remaining inks, * bows to you, the lesson in Tattoos complete."),
			Callback: whenActive(mod.onFinish),
		},
		{
			Kind:     pkg.Output,
			Pattern:  []byte("* bows to you - the lesson in ^ is over."),
			Callback: whenActive(mod.onFinish),
		},
	}

}

func (mod *LearnMultipleLessons) onStart(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		mod.reset()

		number, err := strconv.Atoi(string(match.Captures[0]))
		if err != nil || number <= maxLessons {
			continue
		}

		mod.total = number
		mod.remaining = number
		mod.target = match.Captures[1]

		mod.start = time.Now()
		mod.countdown()
		inout.Input = inout.Input.Replace(match.Index, mod.learn())
	}

	return inout
}

func (mod *LearnMultipleLessons) onBegin(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		// Don't show ongoing learning messages except the very first.
		if mod.total-mod.remaining == maxLessons {
			continue
		}
		inout = mod.onUpdate([]pkg.Match{match}, inout)
	}

	return inout
}

func (mod *LearnMultipleLessons) onUpdate(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		inout.Output = inout.Output.Omit(match.Index)

		mod.countdown()
	}

	return inout
}

func (mod *LearnMultipleLessons) onFinish(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		if mod.remaining <= 0 {
			inout.Output = inout.Output.AddAfter(match.Index, []byte(fmt.Sprintf(
				"%d of %d lessons learned.",
				mod.total-mod.remaining, mod.total,
			)))

			mod.reset()

			continue
		}

		timeleft := ""

		duration := time.Since(mod.start)
		remaining := math.Ceil(float64(mod.remaining) / float64(maxLessons))
		estimate := duration * time.Duration(remaining)

		if mins := estimate.Minutes(); mins >= 1 {
			timeleft += fmt.Sprintf("%.0f minutes ", mins)
			estimate -= time.Duration(mins) * time.Minute
		}
		timeleft += fmt.Sprintf("%.0f seconds", estimate.Seconds())

		inout.Output = inout.Output.Replace(match.Index, []byte(fmt.Sprintf(
			"%d of %d lessons learned, %s remaining.",
			mod.total-mod.remaining, mod.total, timeleft,
		)))

		mod.start = time.Now()
		mod.countdown()
		inout.Input = inout.Input.Add(mod.learn())
	}

	return inout
}

func (mod *LearnMultipleLessons) reset() {
	if mod.timer != nil {
		mod.timer.Stop()
	}

	mod.total = 0
	mod.remaining = 0
	mod.target = []byte{}
	mod.start = time.Time{}
	mod.timer = nil
}

func (mod *LearnMultipleLessons) countdown() {
	if mod.timer != nil {
		mod.timer.Stop()
	}

	mod.timer = time.AfterFunc(15*time.Second, func() {
		mod.reset()
	})
}

func (mod *LearnMultipleLessons) learn() []byte {
	count := maxLessons
	if mod.remaining < count {
		count = mod.remaining
	}
	mod.remaining -= count

	return []byte(fmt.Sprintf("learn %d %s", count, mod.target))
}
