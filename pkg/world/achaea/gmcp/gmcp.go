package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/icza/gox/gox"
)

// @todo Implement the full set:
// - https://nexus.ironrealms.com/GMCP
// - https://nexus.ironrealms.com/GMCP_Data

type Message interface {
	String() string
}

func Parse(command []byte) (Message, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	var message Message

	switch string(parts[0]) {
	case "Char.Items.Inv":
		message = CharItemsInv{}

	case "Char.Name":
		if len(parts) == 1 {
			return nil, fmt.Errorf("missing 'Char.Name' data")
		}

		msg, err := (&CharName{}).Hydrate(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed hydrating 'Char.Name': %w", err)
		}

		message = msg

	case "Char.Status":
		if len(parts) == 1 {
			return nil, fmt.Errorf("missing 'Char.Status' data")
		}

		msg, err := (&CharStatus{}).Hydrate(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed hydrating 'Char.Status': %w", err)
		}

		message = msg

	case "Char.Vitals":
		if len(parts) == 1 {
			return nil, fmt.Errorf("missing 'Char.Vitals' data")
		}

		msg, err := (&CharVitals{}).Hydrate(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed hydrating 'Char.Vitals': %w", err)
		}

		message = msg

	default:
		return nil, fmt.Errorf("invalid message '%s'", parts[0])
	}

	return message, nil
}

type CharItemsInv struct{}

func (msg CharItemsInv) String() string {
	return "Char.Items.Inv"
}

type CharName struct {
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
}

func (msg *CharName) Hydrate(data []byte) (CharName, error) {
	err := json.Unmarshal(data, msg)
	if err != nil {
		return *msg, err
	}

	return *msg, nil
}

func (msg CharName) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Char.Name %s", data)
}

// Hur separerar vi Achaea-specifika saker från generella GMCP-saker?

type CharStatus struct {
	Name             string `json:"name"`
	Fullname         string `json:"fullname"`
	Age              int    `json:"age"`
	Race             string `json:"race"`
	Specialisation   string `json:"specialisation"`
	Level            int    `json:"level"`
	XP               int    `json:"xp"`
	XPRank           int    `json:"xprank"`
	Class            string `json:"class"`
	City             string `json:"city"`
	CityRank         int
	House            string  `json:"house"`
	Order            *string `json:"order"`
	BoundCredits     int     `json:"boundcredits"`
	UnboundCredits   int     `json:"unboundcredits"`
	Lessons          int     `json:"lessons"`
	ExplorerRank     string  `json:"explorerrank"`
	MayanCrowns      int     `json:"mayancrowns"`
	BoundMayanCrowns int     `json:"boundmayancrowns"`
	Gold             int     `json:"gold"`
	Bank             int     `json:"bank"`
	UnreadNews       int     `json:"unread_news"`
	UnreadMessages   int     `json:"unread_msgs"`
	Target           string  `json:"target"`
	Gender           string  `json:"gender"`
}

// Char.Status { "name": "Durak", "fullname": "Mason Durak", "age": "184", "race": "Dwarf", "specialisation": "Brawler", "level": "68 (19%)", "xp": "19%", "xprank": "999", "class": "Monk", "city": "Hashan (1)", "house": "The Somatikos(1)", "order": "(None)", "boundcredits": "20", "unboundcredits": "1", "lessons": "4073", "explorerrank": "an Itinerant", "mayancrowns": "0", "boundmayancrowns": "0", "gold": "35", "bank": "159060", "unread_news": "3751", "unread_msgs": "1", "target": "None", "gender": "male" }

func (msg *CharStatus) Hydrate(data []byte) (CharStatus, error) {
	type CharStatusAlias CharStatus
	var child struct {
		CharStatusAlias
		CBal  *string `json:"bal"`
		CEq   *string `json:"eq"`
		CVote *string `json:"vote"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return *msg, err
	}

	msg = (*CharStatus)(&child.CharStatusAlias)
	// if child.CBal != nil {
	// 	msg.Bal = gox.NewBool(*child.CBal == "1")
	// }

	return *msg, nil
}

func (msg CharStatus) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Char.Status %s", data)
}

type CharVitalsStats struct {
	Bleed    *int
	Ferocity *int
	Kai      *int
	Rage     *int
	Spec     *string
	Stance   *string
}

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
			stats.Bleed = gox.NewInt(value)

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
			stats.Rage = gox.NewInt(value)

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

type CharVitals struct {
	HP     *int `json:"hp,string"`
	MaxHP  *int `json:"maxhp,string"`
	MP     *int `json:"mp,string"`
	MaxMP  *int `json:"maxmp,string"`
	EP     *int `json:"ep,string"`
	MaxEP  *int `json:"maxep,string"`
	WP     *int `json:"wp,string"`
	MaxWP  *int `json:"maxwp,string"`
	NL     *int `json:"nl,string"`
	Bal    *bool
	Eq     *bool
	Vote   *bool
	Prompt *string `json:"string"`

	Stats CharVitalsStats `json:"charstats"`
}

func (msg *CharVitals) Hydrate(data []byte) (CharVitals, error) {
	type CharVitalsAlias CharVitals
	var child struct {
		CharVitalsAlias
		CBal  *string `json:"bal"`
		CEq   *string `json:"eq"`
		CVote *string `json:"vote"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return *msg, err
	}

	msg = (*CharVitals)(&child.CharVitalsAlias)
	if child.CBal != nil {
		msg.Bal = gox.NewBool(*child.CBal == "1")
	}
	if child.CEq != nil {
		msg.Eq = gox.NewBool(*child.CEq == "1")
	}
	if child.CVote != nil {
		msg.Vote = gox.NewBool(*child.CVote == "1")
	}

	return *msg, nil
}

func (msg CharVitals) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Char.Vitals %s", data)
}

type CommChannelPlayers struct{}

func (msg CommChannelPlayers) String() string {
	return "Comm.Channel.Players"
}

type CoreHello struct {
	Client  string `json:"client"`
	Version string `json:"version"`
}

func (msg CoreHello) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Core.Hello %s", data)
}

type CoreSupportsSet struct {
	Char        bool
	CharSkills  bool
	CharItems   bool
	CommChannel bool
	Room        bool
	IRERift     bool
}

func (msg CoreSupportsSet) String() string {
	list := []string{}
	if msg.Char {
		list = append(list, "Char 1")
	}
	if msg.CharSkills {
		list = append(list, "Char.Skills 1")
	}
	if msg.CharItems {
		list = append(list, "Char.Items 1")
	}
	if msg.CommChannel {
		list = append(list, "Comm.Channel 1")
	}
	if msg.Room {
		list = append(list, "Room 1")
	}
	if msg.IRERift {
		list = append(list, "IRE.Rift 1")
	}

	data, err := json.Marshal(list)
	if err != nil {
		data = []byte("[]")
	}

	return fmt.Sprintf("Core.Supports.Set %s", data)
}

type IRERiftRequest struct{}

func (msg IRERiftRequest) String() string {
	return "IRE.Rift.Request"
}
