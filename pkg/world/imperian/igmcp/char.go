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
	Name             *string `json:"name"`
	Fullname         *string `json:"fullname"`
	Age              *int    `json:"age,string"`
	Race             *string `json:"race"`
	Specialisation   *string `json:"specialisation"`
	Level            *int
	XP               *int    `json:"-"`
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
	Target           *int
	Gender           *int // ISO/IEC 5218
}

// Hydrate populates the message with data.
func (msg CharStatus) Hydrate(data []byte) (gmcp.ServerMessage, error) {
	type CharStatusAlias CharStatus
	var child struct {
		CharStatusAlias
		CLevel  *string `json:"level"`
		CCity   *string `json:"city"`
		CHouse  *string `json:"house"`
		COrder  *string `json:"order"`
		CTarget *string `json:"target"`
		CGender *string `json:"gender"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return nil, err
	}

	msg = (CharStatus)(child.CharStatusAlias)

	if child.CLevel != nil {
		level, rank := splitLevelRank(*child.CLevel)
		if rank == nil {
			return nil, fmt.Errorf(
				"failed parsing level '%s'", *child.CLevel,
			)
		}

		msg.Level = gox.NewInt(level)
		msg.XP = rank
	}

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

	if child.CTarget != nil && *child.CTarget != "None" {
		// Yes, sometimes it's a string, sometimes it's an int. Yay!
		target, err := strconv.Atoi(*child.CTarget)
		if err != nil {
			return nil, fmt.Errorf(
				"failed parsing target '%s'", *child.CTarget,
			)
		}

		msg.Target = gox.NewInt(target)
	}

	if child.CGender != nil {
		switch *child.CGender {
		case "male":
			msg.Gender = gox.NewInt(1)
		case "female":
			msg.Gender = gox.NewInt(2)
		default:
			msg.Gender = gox.NewInt(9)
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
