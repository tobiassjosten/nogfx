package achaea

import (
	"math"

	"github.com/tobiassjosten/nogfx/pkg"
)

// type Message struct {
// 	Pattern string
// 	Format  string
// 	Tags    []pkg.Tag
// }

// type Messages []Message

// var (
// 	TekuraFPSidekick = Message{
// 		Pattern: "You kick {*}.",
// 		Format:  "You kick %[1]s",
// 		Tags:    []pkg.Tag{},
// 	}
// 	TekuraFPAttacks = Messages{TekuraFPSidekick}

// 	TelepathyFPMindCrush = Message{
// 		Pattern: "You kick {*}.",
// 		Format:  "You kick %[1]s",
// 		Tags:    []pkg.Tag{},
// 	}
// 	TelepathyFPAttacks = Messages{TelepathyFPMindCrush}

// 	FPAttacks = (Messages{}).Merge(TekuraFPAttacks, TelepathyFPAttacks)

// 	TekuraSPSidekick = Message{
// 		Pattern: "{*} kicks you.",
// 		Format:  "%[1]s kicks you",
// 		Tags:    []pkg.Tag{},
// 	}

// 	TekuraTPSidekick = Message{
// 		Pattern: "{*} kicks {*}.",
// 		Format:  "%[1]s kicks %[2]s",
// 		Tags:    []pkg.Tag{},
// 	}
// )

var (
	attack     = []byte("queue addclear eqbal combo sdk ucp ucp")
	cleareqbal = []byte("clearqueue eqbal")
)

type Bashing struct {
	world *world

	active    bool
	attacking bool
	killed    int
}

// Processor enhances bashing-related tasks in Achaea.
func (bsh *Bashing) Processor() pkg.Processor {
	return pkg.ChainProcessor(
		pkg.MatchInput("kill", bsh.onInit),

		pkg.MatchOutput(
			"You have slain *, retrieving the corpse.",
			bsh.onSlain,
		),

		pkg.MatchOutputs([]string{
			// @todo How can we reuse what we've got in TunnelVision?
			"You pump out at * with a powerful side kick.",
			"You launch a powerful uppercut at *.",
			"A dizzying beam of energy strikes you as your attack rebounds off of *'s shield.",
		}, bsh.onAttack),

		pkg.MatchOutput(
			"A ^ pile of sovereigns spills from the corpse.",
			bsh.onGold,
		),

		pkg.MatchOutput(
			"A nearly invisible magical shield forms around {*}.",
			bsh.onShield,
		),

		pkg.ProcessorFunc(bsh.postprocess),
	)
}

func (bsh *Bashing) postprocess(inout pkg.Inoutput) pkg.Inoutput {
	bsh.attacking = false
	bsh.killed = 0

	return inout
}

func (bsh *Bashing) onInit(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	if bsh.world.Target.isPlayer {
		return inout
	}

	bsh.active = true

	for _, match := range matches {
		inout.ReplaceInput(match.Index, attack)
	}

	return inout
}

func (bsh *Bashing) onSlain(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	if !bsh.active {
		return inout
	}

	for _, match := range matches {
		bsh.killed = int(math.Max(float64(bsh.killed), float64(match.Index)))
	}

	if bsh.world.Target.Queue() > 0 {
		return inout
	}

	bsh.active = false

	inout = inout.RemoveInputMatching(attack)
	bsh.attacking = false

	inout = inout.AppendInput(cleareqbal)

	// @todo Modify output to show percentage/level gain, once we're able
	// to trigger on GMCP messages (Char.Status specifically).

	return inout
}

func (bsh *Bashing) onAttack(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	// With some attacks we might kill the target and start attacking a new
	// one within the same paragraph, so here we make sure we stay active
	// if there are attacks after the kill.
	if bsh.killed > 0 && !bsh.active {
		for _, match := range matches {
			if match.Index > bsh.killed {
				bsh.active = true
			}
		}

		inout = inout.RemoveInputMatching(cleareqbal)
	}

	if !bsh.active {
		return inout
	}

	if bsh.attacking {
		inout = inout.RemoveInputMatching(attack)
	}

	inout = inout.AppendInput(attack)
	bsh.attacking = true

	return inout
}

func (bsh *Bashing) onGold(_ []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	if bsh.killed == 0 {
		return inout
	}

	inout = inout.AppendInput([]byte("get sovereigns"))
	inout = inout.AppendInput([]byte("put sovereigns in pack"))

	return inout
}

func (bsh *Bashing) onShield(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	_ = bsh.world.Target
	return inout
}
