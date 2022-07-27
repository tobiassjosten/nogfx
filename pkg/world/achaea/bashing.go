package achaea

import (
	"math"

	"github.com/tobiassjosten/nogfx/pkg"
)

var (
	attack     = []byte("queue addclear eqbal combo sdk ucp ucp")
	cleareqbal = []byte("clearqueue eqbal")
)

type BashingMod struct {
	world *World

	active    bool
	attacking bool
	killed    int
}

func NewBashingMod(world *World) pkg.Module {
	return &BashingMod{world: world}
}

func (mod *BashingMod) PostInoutput() {
	mod.attacking = false
	mod.killed = 0
}

func (mod *BashingMod) Triggers() []pkg.Trigger {
	return []pkg.Trigger{
		{
			Kind:     pkg.Input,
			Pattern:  []byte("kill"),
			Callback: mod.onInit,
		},

		{
			Kind:     pkg.Output,
			Pattern:  []byte("You have slain *, retrieving the corpse."),
			Callback: mod.onSlain,
		},

		{
			Kind: pkg.Output,
			Patterns: [][]byte{
				// Monk Sidekick.
				[]byte("You pump out at * with a powerful side kick."),
				// Monk Uppercut.
				[]byte("You launch a powerful uppercut at *."),
				// Shield.
				[]byte("A dizzying beam of energy strikes you as your attack rebounds off of *'s shield."),
			},
			Callback: mod.onAttack,
		},

		{
			Kind:     pkg.Output,
			Pattern:  []byte("A ^ pile of sovereigns spills from the corpse."),
			Callback: mod.onGold,
		},
	}
}

func (mod *BashingMod) onInit(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	mod.active = true

	for _, match := range matches {
		inout.ReplaceInput(match.Index, attack)
	}

	return inout
}

func (mod *BashingMod) onSlain(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	if !mod.active {
		return inout
	}

	for _, match := range matches {
		mod.killed = int(math.Max(float64(mod.killed), float64(match.Index)))
	}

	if mod.world.Target.Queue() > 0 {
		return inout
	}

	mod.active = false

	inout = inout.RemoveInputMatching(attack)
	mod.attacking = false

	inout = inout.AppendInput(cleareqbal)

	// @todo Modify output to show percentage/level gain, once we're able
	// to trigger on GMCP messages (Char.Status specifically).

	return inout
}

func (mod *BashingMod) onAttack(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	// With some attacks we might kill the target and start attacking a new
	// one within the same paragraph, so here we make sure we stay active
	// if there are attacks after the kill.
	if mod.killed > 0 && !mod.active {
		for _, match := range matches {
			if match.Index > mod.killed {
				mod.active = true
			}
		}

		inout = inout.RemoveInputMatching(cleareqbal)
	}

	if !mod.active {
		return inout
	}

	if mod.attacking {
		inout = inout.RemoveInputMatching(attack)
	}

	inout = inout.AppendInput(attack)
	mod.attacking = true

	return inout
}

func (mod *BashingMod) onGold(_ []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	if mod.killed == 0 {
		return inout
	}

	inout = inout.AppendInput([]byte("get gold"))
	inout = inout.AppendInput([]byte("put gold in pack"))

	return inout
}
