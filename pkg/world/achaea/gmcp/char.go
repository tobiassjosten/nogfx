package gmcp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/icza/gox/gox"
)

var (
	_ Message = &CharItemsInv{}
	_ Message = &CharName{}
	_ Message = &CharStatus{}
	_ Message = &CharVitals{}
)

// CharItemsInv is a client-sent GMCP message to request a list of items in the
// player's inventory.
type CharItemsInv struct{}

// Hydrate populates the message with data.
func (msg CharItemsInv) Hydrate(_ []byte) (Message, error) {
	return msg, nil
}

// String is the message's string representation.
func (msg CharItemsInv) String() string {
	return "Char.Items.Inv"
}

// CharName is a server-sent GMCP message containing basic information about
// the player's character. Only sent on login.
type CharName struct {
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
}

// Hydrate populates the message with data.
func (msg CharName) Hydrate(data []byte) (Message, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

// String is the message's string representation.
func (msg CharName) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Char.Name %s", data)
}

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
	Target           *string
	Gender           *int // ISO/IEC 5218
}

// Hydrate populates the message with data.
func (msg CharStatus) Hydrate(data []byte) (Message, error) {
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
		return msg, err
	}

	msg = (CharStatus)(child.CharStatusAlias)

	if child.CLevel != nil {
		if level, rank := splitLevelRank(*child.CLevel); level == 0 && rank == 0 {
			return msg, fmt.Errorf("failed parsing level: %w", err)
		} else {
			msg.Level = gox.NewInt(level)
			msg.XP = gox.NewInt(rank)
		}
	}

	if child.CCity != nil && *child.CCity != "(None)" {
		if city, rank := splitRank(*child.CCity); city == "" && rank == 0 {
			return msg, fmt.Errorf("failed parsing city: %w", err)
		} else {
			msg.City = gox.NewString(city)
			msg.CityRank = gox.NewInt(rank)
		}
	}

	if child.CHouse != nil && *child.CHouse != "(None)" {
		if house, rank := splitRank(*child.CHouse); house == "" && rank == 0 {
			return msg, fmt.Errorf("failed parsing city: %w", err)
		} else {
			msg.House = gox.NewString(house)
			msg.HouseRank = gox.NewInt(rank)
		}
	}

	if child.COrder != nil && *child.COrder != "(None)" {
		msg.Order = gox.NewString(*child.COrder)
	}

	if child.CTarget != nil && *child.CTarget != "None" {
		msg.Target = gox.NewString(*child.CTarget)
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

// String is the message's string representation.
func (msg CharStatus) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Char.Status %s", data)
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
func (msg CharVitals) Hydrate(data []byte) (Message, error) {
	type CharVitalsAlias CharVitals
	var child struct {
		CharVitalsAlias
		CBal  string `json:"bal"`
		CEq   string `json:"eq"`
		CVote string `json:"vote"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return msg, err
	}

	msg = (CharVitals)(child.CharVitalsAlias)
	msg.Bal = child.CBal == "1"
	msg.Eq = child.CEq == "1"
	msg.Vote = child.CVote == "1"

	return msg, nil
}

// String is the message's string representation.
func (msg CharVitals) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Char.Vitals %s", data)
}

// CharVitalsStats is structured data extending CharVitals.
type CharVitalsStats struct {
	Bleed int
	Rage  int

	Ferocity *int    // Infernal.
	Kai      *int    // Monk.
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

	err := json.Unmarshal(data, &list)
	if err != nil {
		return err
	}

	for _, item := range list {
		parts := strings.SplitN(item, ": ", 2)
		if len(parts) != 2 {
			return fmt.Errorf("misformed Char.Vitals.charstats '%s'", item)
		}

		switch parts[0] {
		case "Bleed":
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf(
					"failed getting 'Bleed' value for Char.Vitals.charstats: %w",
					err,
				)
			}
			stats.Bleed = value

		case "Ferocity":
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf(
					"failed getting 'Ferocity' value for Char.Vitals.charstats: %w",
					err,
				)
			}
			stats.Ferocity = gox.NewInt(value)

		case "Kai":
			value, err := strconv.Atoi(parts[1][:len(parts[1])-1])
			if err != nil {
				return fmt.Errorf(
					"failed getting 'Kai' value for Char.Vitals.charstats: %w",
					err,
				)
			}
			stats.Kai = gox.NewInt(value)

		case "Rage":
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf(
					"failed getting 'Rage' value for Char.Vitals.charstats: %w",
					err,
				)
			}
			stats.Rage = value

		case "Spec":
			stats.Spec = gox.NewString(parts[1])

		case "Stance":
			stats.Stance = gox.NewString(parts[1])

		default:
			return fmt.Errorf("invalid Char.Vitals.charstats '%s'", item)
		}
	}

	return nil
}