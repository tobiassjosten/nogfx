package achaea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"

	"github.com/icza/gox/gox"
)

// CharStatus is a server-sent GMCP message containing character values. The
// initial message sent contains all values but subsequent messages only carry
// changes, with omitted properties assumed unchanged.
type CharStatus struct {
	*ironrealms.CharStatus
	Age              *int    `json:"age,string,omitempty"`
	BoundCredits     *int    `json:"boundcredits,string,omitempty"`
	BoundMayanCrowns *int    `json:"boundmayancrowns,string,omitempty"`
	ExplorerRank     *string `json:"explorerrank,omitempty"`
	House            *string `json:"-"`
	HouseRank        *int    `json:"-"`
	Lessons          *int    `json:"lessons,string,omitempty"`
	MayanCrowns      *int    `json:"mayancrowns,string,omitempty"`
	Order            *string `json:"-"`
	OrderRank        *int    `json:"-"`
	Specialisation   *string `json:"specialisation,omitempty"`
	Target           *string `json:"target,omitempty"`
	UnboundCredits   *int    `json:"unboundcredits,string,omitempty"`
	XPRank           *int    `json:"xprank,string,omitempty"`
}

// ID is the prefix before the message's data.
func (msg *CharStatus) ID() string {
	return "Char.Status"
}

func (msg *CharStatus) MarshalValue(value *string, rank *int) *string {
	if value == nil {
		return nil
	}

	if *value == "" {
		// Stupid but at least it's consistent…
		return gox.NewString("(None)")
	}

	newvalue := *value
	if rank != nil {
		newvalue = fmt.Sprintf("%s (%d)", newvalue, *rank)
	}

	return &newvalue
}

// Marshal converts the message to a string.
func (msg *CharStatus) Marshal() string {
	proxy := struct {
		*CharStatus
		PCity  *string `json:"city,omitempty"`
		PHouse *string `json:"house,omitempty"`
		PLevel *string `json:"level,omitempty"`
		POrder *string `json:"order,omitempty"`
	}{
		CharStatus: msg,
	}

	proxy.PCity = msg.MarshalCity()
	proxy.PHouse = msg.MarshalValue(msg.House, msg.HouseRank)
	proxy.PLevel = msg.MarshalLevel()
	proxy.POrder = msg.MarshalValue(msg.Order, msg.OrderRank)

	data, _ := json.Marshal(proxy)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharStatus) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	proxy := struct {
		*CharStatus
		PHouse *string `json:"house"`
		POrder *string `json:"order"`
	}{
		CharStatus: &CharStatus{
			CharStatus: &ironrealms.CharStatus{
				CharStatus: &gmcp.CharStatus{},
			},
		},
	}

	err := json.Unmarshal(data, &proxy)
	if err != nil {
		return err
	}

	*msg = CharStatus{}
	if proxy.CharStatus != nil {
		*msg = (CharStatus)(*proxy.CharStatus)
	}

	err = msg.CharStatus.Unmarshal(data)
	if err != nil {
		return err
	}

	if proxy.PHouse != nil {
		if *proxy.PHouse == "(None)" {
			msg.House = gox.NewString("")
		} else {
			house, rank := gmcp.SplitRankInt(*proxy.PHouse)
			msg.House = gox.NewString(house)
			msg.HouseRank = gox.NewInt(rank)
		}
	}

	if proxy.POrder != nil {
		if *proxy.POrder == "(None)" {
			msg.Order = gox.NewString("")
		} else {
			order, rank := gmcp.SplitRankInt(*proxy.POrder)
			msg.Order = gox.NewString(order)
			msg.OrderRank = gox.NewInt(rank)
		}
	}

	if msg.Target != nil {
		target := strings.TrimSuffix(*msg.Target, " (player)")
		if target == "None" {
			target = ""
		}
		msg.Target = &target
	}

	return nil
}

// CharVitals is a server-sent GMCP message containing character attributes.
type CharVitals struct {
	EP     int `json:"ep,string"`
	MaxEP  int `json:"maxep,string"`
	WP     int `json:"wp,string"`
	MaxWP  int `json:"maxwp,string"`
	Vote   bool
	Prompt string `json:"string"`

	Stats CharVitalsStats `json:"charstats"`
}

// ID is the prefix before the message's data.
func (msg *CharVitals) ID() string {
	return "Char.Vitals"
}

// Marshal converts the message to a string.
func (msg *CharVitals) Marshal() string {
	return msg.ID()
}

// Unmarshal populates the message with data.
func (msg *CharVitals) Unmarshal(data []byte) error {
	return nil
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
