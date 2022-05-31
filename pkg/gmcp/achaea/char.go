package achaea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/icza/gox/gox"
)

// CharStatus is a server-sent GMCP message containing character values. The
// initial message sent contains all values but subsequent messages only carry
// changes, with omitted properties assumed unchanged.
type CharStatus struct {
	Age              *int     `json:"age,string,omitempty"`
	Bank             *int     `json:"bank,string,omitempty"`
	BoundCredits     *int     `json:"boundcredits,string,omitempty"`
	BoundMayanCrowns *int     `json:"boundmayancrowns,string,omitempty"`
	City             *string  `json:"-"`
	CityRank         *int     `json:"-"`
	Class            *string  `json:"class,omitempty"`
	ExplorerRank     *string  `json:"explorerrank,omitempty"`
	Fullname         *string  `json:"fullname,omitempty"`
	Gender           *string  `json:"gender,omitempty"`
	Gold             *int     `json:"gold,string,omitempty"`
	House            *string  `json:"-"`
	HouseRank        *int     `json:"-"`
	Lessons          *int     `json:"lessons,string,omitempty"`
	Level            *float64 `json:"-"`
	MayanCrowns      *int     `json:"mayancrowns,string,omitempty"`
	Name             *string  `json:"name,omitempty"`
	Order            *string  `json:"-"`
	OrderRank        *int     `json:"-"`
	Race             *string  `json:"race,omitempty"`
	Specialisation   *string  `json:"specialisation,omitempty"`
	Target           *string  `json:"target,omitempty"`
	UnboundCredits   *int     `json:"unboundcredits,string,omitempty"`
	UnreadMsgs       *int     `json:"unread_msgs,string,omitempty"`
	UnreadNews       *int     `json:"unread_news,string,omitempty"`
	XPRank           *int     `json:"xprank,string,omitempty"`
}

// ID is the prefix before the message's data.
func (msg *CharStatus) ID() string {
	return "Char.Status"
}

func (msg *CharStatus) marshalLevel() *string {
	if msg.Level == nil {
		return nil
	}

	progress := fmt.Sprintf("%.4f", math.Mod(*msg.Level, 1)*100)
	progress = strings.TrimRight(progress, ".0")
	if progress == "" {
		progress = "0"
	}

	return gox.NewString(
		fmt.Sprintf("%d (%s%%)", int(*msg.Level), progress),
	)
}

func (msg *CharStatus) marshalValue(value *string, rank *int) *string {
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

	proxy.PCity = msg.marshalValue(msg.City, msg.CityRank)
	proxy.PHouse = msg.marshalValue(msg.House, msg.HouseRank)
	proxy.PLevel = msg.marshalLevel()
	proxy.POrder = msg.marshalValue(msg.Order, msg.OrderRank)

	data, _ := json.Marshal(proxy)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharStatus) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	if msg == nil {
		*msg = CharStatus{}
	}

	proxy := struct {
		*CharStatus
		PCity  *string `json:"city"`
		PHouse *string `json:"house"`
		PLevel *string `json:"level"`
		POrder *string `json:"order"`
	}{
		CharStatus: msg,
	}

	err := json.Unmarshal(data, &proxy)
	if err != nil {
		return err
	}

	*msg = (CharStatus)(*proxy.CharStatus)

	if proxy.PCity != nil {
		if *proxy.PCity == "(None)" {
			msg.City = gox.NewString("")
		} else {
			city, rank := gmcp.SplitRankInt(*proxy.PCity)
			msg.City = gox.NewString(city)
			msg.CityRank = gox.NewInt(rank)
		}
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

	if proxy.PLevel != nil {
		parts := strings.SplitN(*proxy.PLevel, " ", 2)

		level, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return fmt.Errorf("failed parsing level: %w", err)
		}

		if len(parts) == 2 {
			progressStr := strings.Trim(parts[1], "(%)")
			progress, err := strconv.ParseFloat(progressStr, 64)
			if err != nil {
				return fmt.Errorf("failed parsing level progress: %w", err)
			}

			level += progress / 100
		}

		msg.Level = gox.NewFloat64(level)
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
	HP    int `json:"hp,string"`
	MaxHP int `json:"maxhp,string"`
	MP    int `json:"mp,string"`
	MaxMP int `json:"maxmp,string"`
	EP    int `json:"ep,string"`
	MaxEP int `json:"maxep,string"`
	WP    int `json:"wp,string"`
	MaxWP int `json:"maxwp,string"`
	NL    int `json:"nl,string"`

	Bal bool `json:"-"`
	Eq  bool `json:"-"`

	Vote bool `json:"-"`

	Prompt string `json:"string"`

	Stats CharVitalsStats `json:"charstats"`
}

// ID is the prefix before the message's data.
func (msg *CharVitals) ID() string {
	return "Char.Vitals"
}

func (msg *CharVitals) marshalBoolStringInt(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// Marshal converts the message to a string.
func (msg *CharVitals) Marshal() string {
	proxy := struct {
		*CharVitals
		PBal string `json:"bal"`
		PEq  string `json:"eq"`
	}{
		CharVitals: msg,
	}

	proxy.PBal = msg.marshalBoolStringInt(msg.Bal)
	proxy.PEq = msg.marshalBoolStringInt(msg.Eq)

	data, _ := json.Marshal(proxy)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharVitals) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	if msg == nil {
		*msg = CharVitals{}
	}

	proxy := struct {
		*CharVitals
		PBal  string `json:"bal"`
		PEq   string `json:"eq"`
		PVote string `json:"vote"`
	}{
		CharVitals: msg,
	}

	err := json.Unmarshal(data, &proxy)
	if err != nil {
		return err
	}

	*msg = (CharVitals)(*proxy.CharVitals)

	msg.Bal = proxy.PBal == "1"
	msg.Eq = proxy.PEq == "1"
	msg.Vote = proxy.PVote == "1"

	return nil
}

// CharVitalsStats is structured data extending CharVitals.
type CharVitalsStats struct {
	Bleed int
	Rage  int

	Ferocity *int
	Kai      *int
	Karma    *int
	Spec     *string
	Stance   *string

	// @todo Implement the following ones (first checking keys in game).
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

// MarshalJSON transforms the struct to JSON.
func (as *CharVitalsStats) MarshalJSON() ([]byte, error) {
	list := []string{
		fmt.Sprintf("Bleed: %d", as.Bleed),
		fmt.Sprintf("Rage: %d", as.Rage),
	}

	if as.Ferocity != nil {
		list = append(list, fmt.Sprintf("Ferocity: %d", *as.Ferocity))
	}
	if as.Kai != nil {
		list = append(list, fmt.Sprintf("Kai: %d%%", *as.Kai))
	}
	if as.Karma != nil {
		list = append(list, fmt.Sprintf("Karma: %d%%", *as.Karma))
	}
	if as.Spec != nil {
		list = append(list, fmt.Sprintf("Spec: %s", *as.Spec))
	}
	if as.Stance != nil {
		list = append(list, fmt.Sprintf("Stance: %s", *as.Stance))
	}

	return json.Marshal(list)
}

// UnmarshalJSON hydrates the struct from JSON.
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
