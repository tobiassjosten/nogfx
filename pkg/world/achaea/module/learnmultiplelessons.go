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
	world pkg.World

	total     int
	remaining int
	target    []byte
	start     time.Time
	timer     *time.Timer
}

// NewLearnMultipleLessons creates a new LearnMultipleLessons module.
func NewLearnMultipleLessons(world pkg.World) pkg.Module {
	return &LearnMultipleLessons{
		world: world,
	}
}

func (mod *LearnMultipleLessons) InputTriggers() []pkg.Trigger[pkg.Input] {
	return []pkg.Trigger[pkg.Input]{
		{
			Pattern:  []byte("learn {^} {^ from *}"),
			Callback: mod.onStart,
		},
	}
}

func (mod *LearnMultipleLessons) OutputTriggers() []pkg.Trigger[pkg.Output] {
	return []pkg.Trigger[pkg.Output]{
		{
			Pattern:  []byte("* begins the lesson in ^."),
			Callback: mod.onBegin,
		},
		{
			Pattern:  []byte("* bows to you and commences the lesson in ^."),
			Callback: mod.onBegin,
		},
		{
			Pattern:  []byte("* continues your training in ^."),
			Callback: mod.onUpdate,
		},
		{
			Pattern:  []byte("* finishes the lesson in ^."),
			Callback: mod.onFinish,
		},
		{
			Pattern:  []byte("Storing ^ remaining inks, * bows to you, the lesson in Tattoos complete."),
			Callback: mod.onFinish,
		},
		{
			Pattern:  []byte("* bows to you - the lesson in ^ is over."),
			Callback: mod.onFinish,
		},
	}

}

func (mod *LearnMultipleLessons) onStart(match pkg.TriggerMatch[pkg.Input]) pkg.Input {
	input := match.Content

	mod.reset()

	number, err := strconv.Atoi(string(match.Captures[0]))
	if err != nil || number <= maxLessons {
		return input
	}

	mod.total = number
	mod.remaining = number
	mod.target = match.Captures[1]

	mod.start = time.Now()
	mod.countdown()
	input = input.Replace(match.Index, mod.learn())

	return input
}

func (mod *LearnMultipleLessons) onBegin(match pkg.TriggerMatch[pkg.Output]) pkg.Output {
	// Don't show ongoing learning messages except the very first.
	if mod.total-mod.remaining == maxLessons {
		return match.Content
	}
	return mod.onUpdate(match)
}

func (mod *LearnMultipleLessons) onUpdate(match pkg.TriggerMatch[pkg.Output]) pkg.Output {
	output := match.Content

	if mod.timer == nil {
		return output
	}

	output = output.Remove(match.Index)

	mod.countdown()

	return output
}

func (mod *LearnMultipleLessons) onFinish(match pkg.TriggerMatch[pkg.Output]) pkg.Output {
	output := match.Content

	if mod.remaining <= 0 {
		output = output.AddAfter(match.Index, []byte(fmt.Sprintf(
			"%d of %d lessons learned.",
			mod.total-mod.remaining, mod.total,
		)))

		mod.reset()

		return output
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

	output = output.Replace(match.Index, []byte(fmt.Sprintf(
		"%d of %d lessons learned, %s remaining.",
		mod.total-mod.remaining, mod.total, timeleft,
	)))

	mod.start = time.Now()
	mod.countdown()
	mod.world.Send(mod.learn())

	return output
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
