package achaea

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg"
	"golang.org/x/exp/slices"
)

const (
	OmitTag pkg.Tag = "omit"

	MyAttackTag           pkg.Tag = "myattack"
	MyAttackVisionTag     pkg.Tag = "myattackvision"
	MyAttackTargetTag     pkg.Tag = "myattacktarget"
	MyAttackHitTag        pkg.Tag = "myattackhit"
	MyAttackMissTag       pkg.Tag = "myattackmiss"
	MyAttackDodgeTag      pkg.Tag = "myattackdodge"
	MyAttackReflectionTag pkg.Tag = "myattackreflection"
	MyAttackUnshieldedTag pkg.Tag = "myattackunshielded"
	MyAttackLimbDmgTag    pkg.Tag = "myattacklimbdmg"
	MyAttackAwokeTag      pkg.Tag = "myattackawoke"
	MyAttackOmitTag       pkg.Tag = "myattackomit"
	MyAttackCrit1Tag      pkg.Tag = "myattackcrit1"
	MyAttackCrit2Tag      pkg.Tag = "myattackcrit2"
	MyAttackCrit3Tag      pkg.Tag = "myattackcrit3"
	MyAttackCrit4Tag      pkg.Tag = "myattackcrit4"
	MyAttackCrit5Tag      pkg.Tag = "myattackcrit5"
	MyAttackNPCDefence    pkg.Tag = "myattacknpcdefence"

	MyCuringTag pkg.Tag = "curing"
	MyCuredTag  pkg.Tag = "cured"
)

// TunnelVision filters and rewrites output to make it easier and faster to
// perceive in large quantities.
type TunnelVision struct {
	character *Character
}

func (tv TunnelVision) preprocessor() pkg.Processor {
	ps := []pkg.Processor{}

	// GMCP for defs/affs doesn't always come before the output. So we need
	// to first parse all those, then parse the IO.

	for _, group := range myAttackPatterns {
		for attack, pattern := range group {
			ps = append(ps, pkg.MatchOutput(
				pattern,
				tv.cbMyAttackTag(attack),
			))
		}
	}

	for tag, pattern := range myCriticalHitPatterns {
		ps = append(ps, pkg.MatchOutput(
			pattern,
			tv.cbTag(tag),
		))
	}

	return pkg.ChainProcessor(append(ps,
		pkg.MatchOutputs([]string{
			"You may drink another health or mana elixir.",
			"You may eat another bit of irid moss or potash.",
			"You may apply another salve to yourself.",
		}, tv.cbTag(OmitTag)),

		pkg.MatchOutputs([]string{
			// Afflictions.
			"You exert superior mental control and your wounds clot before your eyes.",

			// Balance.
			"You may drink another health or mana elixir.",
			"You may eat another bit of irid moss or potash.",
			"You may apply another salve to yourself.",

			// Defences.
			"You shut your eyes and concentrate on the Soulrealms. A moment later, you feel inextricably linked with the realm of Death.", // deathsight
			"Your vision sharpens with light as you gain night sight.",                                                                    // nightsight
			"Your body begins to feel lighter and you feel that you are floating slightly.",                                               // levitating
			"A brief shiver runs through your body.",                                                                                      // weathering
			"Flexing your muscles, you concentrate on forcing unnatural toughness over the surface of your skin.",                         // toughness

			// Queue.
			"[System]: Queued ^ commands cleared.",
			"[System]: Added * to your ^ queue.",
			"[System]: Running queued ^ command: *",

			// Weather.
			"Occasional drops of rain fall to the ground from a sky grey with pregnant clouds.",
			"Occasional raindrops fall on your head as the drizzle continues.",
		}, tv.cbTag(OmitTag)),

		pkg.MatchOutputs([]string{
			"You connect!",
			"You connect to the ^!",
			"You connect to the ^ ^!",

			// Monk attacks.
			"With deadly precision, you quickly jab the nerve cluster in the ^ shoulder of *.",
			"The palmstrike smashes into the temple of *.",
			"As your kick strikes the magical shield surrounding *, it bursts into shimmering fragments.",
			"You knock the legs out from under * and send ^ sprawling.",
			"You pick up ^ prone form, turn, and smash ^ into the ground.",
			"You lift * triumphantly into the air, then yank ^ down into your raised knee with back breaking force.",
			"* is hurled to the {^} by the strength of the blow.",
		}, tv.cbTag(MyAttackHitTag)),

		pkg.MatchOutputs([]string{
			"* twists ^ body out of harm's way.",
			"* backs away and out of your reach.",
		}, tv.cbTag(MyAttackDodgeTag)),

		pkg.MatchOutput("You miss.", tv.cbTag(MyAttackMissTag)),

		pkg.MatchOutput(
			"Your kick scythes through nothing, hitting only empty air.",
			tv.cbTag(MyAttackUnshieldedTag),
		),

		pkg.MatchOutput(
			"A reflection of ^ blinks out of existence.",
			tv.cbTag(MyAttackReflectionTag),
		),

		pkg.MatchOutput(
			"* stands firm and does not budge against the thrust kick.",
			tv.cbTag(MyAttackDodgeTag, "sturdiness"),
		),

		pkg.MatchOutput(
			"* ceases tending to ^ wounds.",
			tv.cbTag(MyAttackAwokeTag),
		),

		pkg.MatchOutputs([]string{
			"As your blow lands with a crunch, you perceive that you have dealt {*%} damage to ^'s ^ ^.",
			"As your blow lands with a crunch, you perceive that you have dealt {*%} damage to ^'s ^.",
		}, tv.cbTag(MyAttackLimbDmgTag)),

		pkg.MatchOutputs([]string{
			"Ahh, I am truly sorry, but I do not see anyone by that name here.",
			"I do not recognise anything called that here.",
			"Nothing can be seen here by that name.",
			"You cannot see that being here.",
			"You detect nothing here by that name.",
		}, tv.cbTag(MyAttackOmitTag)),

		pkg.MatchOutput(
			"* springs to *'s defence.",
			tv.cbTag(MyAttackNPCDefence),
		),

		pkg.MatchOutput(
			"You bleed {^} health.",
			tv.cbTag(OmitTag, func(match pkg.Match) bool {
				bleed, err := strconv.Atoi(string(match.Captures[0]))
				if err != nil {
					log.Printf("failed parsing bleed: %s", err)
					return false
				}
				return bleed < tv.character.MaxHealth/10
			}),
		),

		// @todo Don't omit these if they are manually initiated!
		pkg.MatchOutputs([]string{
			"You take a drink from *.",
			"You down the last drop from *.",
			"You eat a potash crystal.",
			"You take out some salve and quickly rub it on your skin.",
		}, tv.cbTag(MyCuringTag)),

		pkg.MatchOutputs([]string{
			"The elixir heals and soothes you.",
			"Your mind feels stronger and more alert.",
			"You feel your health and mana replenished.",
			"A feeling of comfortable warmth spreads over you.",                             // insulation
			"Your body begins to feel lighter and you feel that you are floating slightly.", // levitating
		}, tv.cbTag(MyCuredTag)),

		// You remove 1 potash, bringing the total in the rift to 63.

		pkg.MatchOutputs([]string{
			// Genji, Atavian.
			"{An atavian traveller} opens a gash in your arm with a skilful blow.",
			"{An atavian traveller} slashes you viciously.",
			"Drawing forth an ancient obsidian scimitar, {Gaharas, daughter of Gehan} slashes into you with a brutal strike.",
			"His impish eyes shining, {Tinja, the atavian} swiftly throws a series of kicks at you.",
			"Squinting at you, {the atavian shaman} brings down his staff heavily.",
			"With a grim look on ^ face, {*} slashes you with a scimitar.",
			"With a powerful swing, {*} hacks into you with ^ massive halberd.",
			"{Gaharas, daughter of Gehan} rushes toward you, unleashing a flurry of talons.",
			"{Gaharas, daughter of Gehan} rises slightly off the ground, buffeting her wings in furious defense of her people.",

			// Genji, Manticore.
			"{*} leaps at you, claws blazing.",
			"Leaping at you, {*} tears at your flesh with its fangs.",

			// Genji, Atavian + Manticore.
			"{*} flaps ^ wings and buffets you with a gust of wind.",
		}, tv.onNPCAttack),
	)...)
}

func (tv TunnelVision) rewriteProcessor() pkg.Processor {
	// @todo Color texts (e.g. MyAttackNPCDefence).
	return pkg.ChainProcessor(
		pkg.ProcessorFunc(tv.rewriteOmits),
		pkg.ProcessorFunc(tv.rewriteMyAttacks),
	)
}

var myCriticalHitPatterns = map[pkg.Tag]string{
	MyAttackCrit1Tag: "You have scored a CRITICAL hit!",
	MyAttackCrit2Tag: "You have scored a CRUSHING CRITICAL hit!",
	MyAttackCrit3Tag: "You have scored an OBLITERATING CRITICAL hit!",
	MyAttackCrit4Tag: "You have scored an ANNIHILATINGLY POWERFUL CRITICAL hit!",
	MyAttackCrit5Tag: "You have scored a WORLD-SHATTERING CRITICAL HIT!!!",
}

func (TunnelVision) cbMyAttackTag(attack string) pkg.Callback {
	return func(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
		for _, match := range matches {
			inout = inout.AddOutputTag(match.Index, MyAttackTag, attack)

			target := strings.ToLower(string(match.Captures[0]))
			target = strings.ToUpper(target[:1]) + target[1:]
			inout = inout.AddOutputTag(match.Index, MyAttackTargetTag, target)
		}

		return inout
	}
}

func (TunnelVision) cbTag(tag pkg.Tag, values ...any) pkg.Callback {
	return func(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
		for _, match := range matches {
			if len(values) == 0 && len(match.Captures) > 0 {
				values = append(values, match.Captures[0])
			}
			if len(values) > 0 {
				if f, ok := values[0].(func(pkg.Match) bool); ok {
					if !f(match) {
						continue
					}
				} else if f, ok := values[0].(func(pkg.Match) any); ok {
					values[0] = f(match)
				}
			}

			inout = inout.AddOutputTag(match.Index, tag, values)
		}

		return inout
	}
}

func (TunnelVision) onNPCAttack(matches []pkg.Match, inout pkg.Inoutput) pkg.Inoutput {
	for _, match := range matches {
		npc := strings.ToLower(string(match.Captures[0]))
		npc = strings.ToUpper(npc[:1]) + npc[1:]

		inout = inout.ReplaceOutput(match.Index, []byte(fmt.Sprintf(
			"%s \x1a[31mattacks\x1a] you.", npc,
		)))
	}

	return inout
}

func (tv TunnelVision) rewriteMyAttacks(inout pkg.Inoutput) pkg.Inoutput {
	previous := []string{}

	for i := 0; i < len(inout.Output); i++ {
		tags := inout.Output[i].Tags

		if _, ok := tags[MyAttackTag]; !ok {
			previous = []string{}
			continue
		}

		if len(previous) > 0 {
			inout = inout.RemoveOutput(i)
			i--
		}

		vision := tv.myAttackVision(tags, previous)

		data := []byte(strings.Join(vision, " ") + ".")
		inout = inout.ReplaceOutput(i, data)

		previous = vision
	}

	return inout
}

var tagTexts = map[pkg.Tag]string{
	MyAttackMissTag:       "\x1a[33mmiss\x1a]",
	MyAttackDodgeTag:      "\x1a[33mdodge\x1a]",
	MyAttackReflectionTag: "\x1a[33mreflection\x1a]",
	MyAttackUnshieldedTag: "\x1a[33munshielded\x1a]",
	MyAttackAwokeTag:      "awoke",
	MyAttackCrit1Tag:      "x2",
	MyAttackCrit2Tag:      "x4",
	MyAttackCrit3Tag:      "x8",
	MyAttackCrit4Tag:      "\x1a[32mx16\x1a]",
	MyAttackCrit5Tag:      "\x1a[1;32mx32\x1a]",
}

func (TunnelVision) myAttackVision(tags pkg.Tags, vision []string) []string {
	if len(vision) == 0 {
		vision = []string{"You \x1a[32;1mattack\x1a]"}

		if target, ok := tags[MyAttackTargetTag]; ok {
			if starget, ok := target.(string); ok {
				vision = append(vision, starget)
			}
		}
	}

	if attack, ok := tags[MyAttackTag]; ok {
		if sattack, ok := attack.(string); ok {
			sattack = "\x1a[32;1m" + sattack + "\x1a]"
			vision = append(vision, "/", sattack)
		}
	}

	for tag, value := range tags {
		if text, ok := tagTexts[tag]; ok {
			vision = append(vision, text)
		} else if tag == MyAttackLimbDmgTag {
			vision = append(vision, value.(string))
		}
	}

	return vision
}

func (TunnelVision) rewriteOmits(inout pkg.Inoutput) pkg.Inoutput {
	var ptags pkg.Tags

	myAttackTags := []pkg.Tag{
		MyAttackHitTag,
		MyAttackMissTag,
		MyAttackDodgeTag,
		MyAttackReflectionTag,
		MyAttackUnshieldedTag,
		MyAttackAwokeTag,
		MyAttackOmitTag,
		MyAttackCrit1Tag,
		MyAttackCrit2Tag,
		MyAttackCrit3Tag,
		MyAttackCrit4Tag,
		MyAttackCrit5Tag,
	}

	shouldOmit := func(tags, ptags, ntags, iotags pkg.Tags) bool {
		for tag := range tags {
			if tag == OmitTag {
				return true
			}

			if tag == MyAttackTag {
				inout = inout.AddTag(MyAttackTag)
			}
			if tag == MyAttackOmitTag {
				_, ok := inout.Tags[MyAttackTag]
				return ok
			}

			if slices.Contains(myAttackTags, tag) {
				_, ok := ptags[MyAttackTag]
				return ok
			}

			if tag == MyCuringTag {
				_, ok := ntags[MyCuredTag]
				return ok
			}

			if tag == MyCuredTag {
				_, ok := ptags[MyCuringTag]
				return ok
			}
		}

		return false
	}

	for i := 0; i < len(inout.Output); i++ {
		ntags := pkg.Tags{}
		if len(inout.Output) > i+1 {
			ntags = inout.Output[i+1].Tags
		}

		omitting := shouldOmit(inout.Output[i].Tags, ptags, ntags, inout.Tags)
		ptags = inout.Output[i].Tags

		if !omitting {
			continue
		}

		if i > 0 {
			for t, v := range inout.Output[i].Tags {
				inout = inout.AddOutputTag(i-1, t, v)
			}
		}

		inout = inout.RemoveOutput(i)
		i--

		continue
	}

	return inout
}
