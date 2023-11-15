package process

const (
	PromptTag Tag = "prompt"
)

// @todo Figure out how to use Tag and implement it for that use.

type Tag string

type Tags map[Tag]any

// "omit"

// "mevs"
// "vsme"
// "allyvs"
// "vsfoe"

// "subject" or "subjects"
// "verb" or "verbs"
// "object" or "objects"
// "adjective" or "adjectives"

// T.ex: "mevs" "vsfoe" "object:someone" "verb:axekick" "effect:
// hur hantera dodge, critx2, afflictions, cures, defenses, etc?

// You take a drink from a black vial. -> möjliggör curing a, b, och c
// Your asthma clears up. -> x inte möjliggjord, illusion

func (tags Tags) Add(name Tag, values ...any) Tags {
	if tags == nil {
		tags = Tags{}
	}

	var value any = struct{}{}
	if len(values) > 0 {
		value = values[0]
	}

	tags[name] = value

	return tags
}
