package achaea

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/process"
)

// @todo Make this use 20 lessons at a time with the myrrh/bisemutum defense:
// "Your mind is racing with enhanced speed."
var maxLessons = 15

// Learning lets players learn an unlimited amount of lessons in one swoop by
// automatically chaining learning sessions together.
type Learning struct {
	total     int
	remaining int
	target    []byte
	start     time.Time
	timer     *time.Timer
}

// Processor enhances learning-related tasks in Achaea.
func (lrn *Learning) Processor() pkg.Processor {
	whenActive := func(callback pkg.Callback) pkg.Callback {
		return func(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
			if lrn.timer == nil {
				return inout
			}
			return callback(matches, inout)
		}
	}

	process.MatchInput(
		"learn {^} {^ from *}",
		func(m process.Match, ins, outs [][]byte) ([][]byte, [][]byte, error) {
			return nil, nil, nil
		},
	)

	return pkg.ChainProcessor(
		pkg.MatchInput("learn {^} {^ from *}", lrn.onStart),
		pkg.MatchOutputs([]string{
			"* begins the lesson in ^.",
			"* bows to you and commences the lesson in ^.",
		}, whenActive(lrn.onBegin)),
		pkg.MatchOutput(
			"* continues your training in ^.",
			whenActive(lrn.onUpdate),
		),
		pkg.MatchOutputs([]string{
			"* finishes the lesson in ^.",
			"Storing ^ remaining inks, * bows to you, the lesson in Tattoos complete.",
			"* bows to you - the lesson in ^ is over.",
		}, whenActive(lrn.onFinish)),
	)
}

func (lrn *Learning) onStart(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		lrn.reset()

		number, err := strconv.Atoi(string(match.Captures[0]))
		if err != nil || number <= maxLessons {
			continue
		}

		lrn.total = number
		lrn.remaining = number
		lrn.target = match.Captures[1]

		lrn.start = time.Now()
		lrn.countdown()
		inout.Input = inout.Input.Replace(match.Index, lrn.learn())
	}

	return inout
}

func (lrn *Learning) onBegin(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		// Don't show ongoing learning messages except the very first.
		if lrn.total-lrn.remaining == maxLessons {
			continue
		}
		inout = lrn.onUpdate([]pkg.Match{match}, inout)
	}

	return inout
}

func (lrn *Learning) onUpdate(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		inout.Output = inout.Output.Remove(match.Index)

		lrn.countdown()
	}

	return inout
}

func (lrn *Learning) onFinish(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		if lrn.remaining <= 0 {
			inout.Output = inout.Output.AddAfter(match.Index, []byte(fmt.Sprintf(
				"%d of %d lessons learned.",
				lrn.total-lrn.remaining, lrn.total,
			)))

			lrn.reset()

			continue
		}

		timeleft := ""

		duration := time.Since(lrn.start)
		remaining := math.Ceil(float64(lrn.remaining) / float64(maxLessons))
		estimate := duration * time.Duration(remaining)

		if mins := estimate.Minutes(); mins >= 1 {
			timeleft += fmt.Sprintf("%.0f minutes ", mins)
			estimate -= time.Duration(mins) * time.Minute
		}
		timeleft += fmt.Sprintf("%.0f seconds", estimate.Seconds())

		inout.Output = inout.Output.Replace(match.Index, []byte(fmt.Sprintf(
			"%d of %d lessons learned, %s remaining.",
			lrn.total-lrn.remaining, lrn.total, timeleft,
		)))

		lrn.start = time.Now()
		lrn.countdown()
		inout = inout.AppendInput(lrn.learn())
	}

	return inout
}

func (lrn *Learning) reset() {
	if lrn.timer != nil {
		lrn.timer.Stop()
	}

	lrn.total = 0
	lrn.remaining = 0
	lrn.target = []byte{}
	lrn.start = time.Time{}
	lrn.timer = nil
}

func (lrn *Learning) countdown() {
	if lrn.timer != nil {
		lrn.timer.Stop()
	}

	lrn.timer = time.AfterFunc(15*time.Second, func() {
		lrn.reset()
	})
}

func (lrn *Learning) learn() []byte {
	count := maxLessons
	if lrn.remaining < count {
		count = lrn.remaining
	}
	lrn.remaining -= count

	return []byte(fmt.Sprintf("learn %d %s", count, lrn.target))
}
