package igmcp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/icza/gox/gox"
)

var (
	_ gmcp.ServerMessage = &CharStatus{}
	_ gmcp.ServerMessage = &CharVitals{}
)

// CharStatus is a server-sent GMCP message containing character values. The
// initial message sent contains all values but subsequent messages only carry
// changes, with omitted properties assumed unchanged.
type CharStatus struct {
	gmcp.CharStatus
	Class *string `json:"class"`
}

// Hydrate populates the message with data.
func (msg CharStatus) Hydrate(data []byte) (gmcp.ServerMessage, error) {
	parentMessage, err := msg.CharStatus.Hydrate(data)
	if err != nil {
		return nil, err
	}

	parentMsg, ok := parentMessage.(gmcp.CharStatus)
	if !ok {
		return nil, fmt.Errorf(
			"expected gmcp.CharStatus, got '%+v'", parentMessage,
		)
	}
	msg.CharStatus = parentMsg

	err = json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharVitals is a server-sent GMCP message containing character attributes.
type CharVitals struct {
	HP     int `json:"hp,string"`
	MaxHP  int `json:"maxhp,string"`
	MP     int `json:"mp,string"`
	MaxMP  int `json:"maxmp,string"`
	EP     int `json:"ep,string"`
	MaxEP  int `json:"maxep,string"`
	WP     int `json:"wp,string"`
	MaxWP  int `json:"maxwp,string"`
	NL     int `json:"nl,string"`
	Bal    bool
	Eq     bool
	Vote   bool
	Prompt string `json:"string"`

	Stats CharVitalsStats `json:"charstats"`
}

// Hydrate populates the message with data.
func (msg CharVitals) Hydrate(data []byte) (gmcp.ServerMessage, error) {
	type CharVitalsAlias CharVitals
	var child struct {
		CharVitalsAlias
		CBal  string `json:"bal"`
		CEq   string `json:"eq"`
		CVote string `json:"vote"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return nil, err
	}

	msg = (CharVitals)(child.CharVitalsAlias)
	msg.Bal = child.CBal == "1"
	msg.Eq = child.CEq == "1"
	msg.Vote = child.CVote == "1"

	return msg, nil
}

// CharVitalsStats is structured data extending CharVitals.
type CharVitalsStats struct {
	Bleed int
	Rage  int

	Ferocity *int    // Infernal.
	Kai      *int    // Monk.
	Spec     *string // Infernal, Paladin, Runewarden.
	Stance   *string // Bard, Blademaster, Monk.
	Karma    *int

	// @todo Implement the one following (first checking keys in game).
	// Channels // Magi.
	// CurrentMorph // Druid, Sentinel.
	// Devotion // Paladin, Priest.
	// ElementalChannels // Sylvan.
	// EntityBalance // Occultist.
	// Essence // Apostate.
	// Karma // Occultist.
	// NumberOfSpiritsBound // Shaman.
	// SecretedVenom // Serpent.
	// SunlightEnergy // Druid, Sylvan.
	// VoiceBalance // Bard.

}

// UnmarshalJSON hydrates CharVitalsStats from a list of unstructured strings.
func (stats *CharVitalsStats) UnmarshalJSON(data []byte) error {
	var list []string

	// This should only be invoked from CharVitals.UnmarshalJSON(), so any
	// formatting errors will be caught there.
	_ = json.Unmarshal(data, &list)

	for _, item := range list {
		parts := strings.SplitN(item, ": ", 2)
		if len(parts) != 2 {
			return fmt.Errorf("misformed charstat '%s'", item)
		}

		switch parts[0] {
		case "Bleed":
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid charstat '%s'", item)
			}
			stats.Bleed = value

		case "Rage":
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid charstat '%s'", item)
			}
			stats.Rage = value

		case "Ferocity":
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid charstat '%s'", item)
			}
			stats.Ferocity = gox.NewInt(value)

		case "Kai":
			value, err := strconv.Atoi(parts[1][:len(parts[1])-1])
			if err != nil {
				return fmt.Errorf("invalid charstat '%s'", item)
			}
			stats.Kai = gox.NewInt(value)

		case "Karma":
			value, err := strconv.Atoi(parts[1][:len(parts[1])-1])
			if err != nil {
				return fmt.Errorf("invalid charstat '%s'", item)
			}
			stats.Karma = gox.NewInt(value)

		case "Spec":
			stats.Spec = gox.NewString(parts[1])

		case "Stance":
			if parts[1] != "None" {
				stats.Stance = gox.NewString(parts[1])
			}

		default:
			return fmt.Errorf("invalid charstat '%s'", item)
		}
	}

	return nil
}
