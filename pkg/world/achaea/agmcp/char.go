package agmcp

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
	Specialisation   *string `json:"specialisation"`
	XPRank           *int    `json:"xprank,string"`
	Class            *string `json:"class"`
	City             *string
	CityRank         *int
	House            *string
	HouseRank        *int
	Order            *string
	OrderRank        *int
	BoundCredits     *int    `json:"boundcredits,string"`
	UnboundCredits   *int    `json:"unboundcredits,string"`
	Lessons          *int    `json:"lessons,string"`
	ExplorerRank     *string `json:"explorerrank"`
	MayanCrowns      *int    `json:"mayancrowns,string"`
	BoundMayanCrowns *int    `json:"boundmayancrowns,string"`
	Gold             *int    `json:"gold,string"`
	Bank             *int    `json:"bank,string"`
	UnreadNews       *int    `json:"unread_news,string"`
	UnreadMessages   *int    `json:"unread_msgs,string"`
	Target           *string `json:"target"`
}

// Hydrate populates the message with data.
func (msg CharStatus) Hydrate(data []byte) (gmcp.ServerMessage, error) {
	type CharStatusAlias CharStatus
	var child struct {
		CharStatusAlias
		CCity  *string `json:"city"`
		CHouse *string `json:"house"`
		COrder *string `json:"order"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return nil, err
	}

	msg = (CharStatus)(child.CharStatusAlias)

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

	// @todo The CharStatus fields are only set when they change (or are
	// first initiated). So rather than treating "(None)" as a non-value,
	// it should be considered as "", to communicate that a change happened
	// and that the new value is empty. Like .Target does it below.

	if child.CCity != nil && *child.CCity != "(None)" {
		city, rank := splitRank(*child.CCity)
		if rank == nil {
			return nil, fmt.Errorf(
				"failed parsing city '%s'", *child.CCity,
			)
		}

		msg.City = gox.NewString(city)
		msg.CityRank = rank
	}

	if child.CHouse != nil && *child.CHouse != "(None)" {
		house, rank := splitRank(*child.CHouse)
		if rank == nil {
			return nil, fmt.Errorf(
				"failed parsing house '%s'", *child.CHouse,
			)
		}

		msg.House = gox.NewString(house)
		msg.HouseRank = rank
	}

	if child.COrder != nil && *child.COrder != "(None)" {
		order, rank := splitRank(*child.COrder)
		if rank == nil {
			return nil, fmt.Errorf(
				"failed parsing order '%s'", *child.COrder,
			)
		}

		msg.Order = gox.NewString(order)
		msg.OrderRank = rank
	}

	if msg.Target != nil {
		target := strings.TrimSuffix(*msg.Target, " (player)")
		if target == "None" {
			msg.Target = gox.NewString("")
		} else {
			msg.Target = &target
		}
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

	Ferocity *int // Infernal.
	Kai      *int // Monk.
	Karma    *int
	Spec     *string // Infernal, Paladin, Runewarden.
	Stance   *string // Bard, Blademaster, Monk.

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
