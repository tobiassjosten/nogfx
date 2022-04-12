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
// - https://github.com/keneanung/GMCPAdditions

type Message interface {
	String() string
}

type Hydrator interface {
	Hydrate([]byte) (Message, error)
}

func Parse(command []byte) (Message, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	var hydrator Hydrator

	switch string(parts[0]) {
	case "Char.Items.Inv":
		return CharItemsInv{}, nil

	case "Char.Name":
		hydrator = CharName{}

	case "Char.Status":
		hydrator = CharStatus{}

	case "Char.Vitals":
		hydrator = CharVitals{}

	default:
		return nil, fmt.Errorf("unknown message '%s'", parts[0])
	}

	if len(parts) == 1 {
		return nil, fmt.Errorf("missing '%T' data", hydrator)
	}

	return hydrator.Hydrate(parts[1])
}

type CharItemsInv struct{}

func (msg CharItemsInv) String() string {
	return "Char.Items.Inv"
}

type CharName struct {
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
}

func (msg CharName) Hydrate(data []byte) (Message, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func (msg CharName) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Char.Name %s", data)
}

type CharStatus struct {
	Name             string `json:"name"`
	Fullname         string `json:"fullname"`
	Age              int    `json:"age,string"`
	Race             string `json:"race"`
	Specialisation   string `json:"specialisation"`
	Level            int
	XP               int    `json:"-"`
	XPRank           int    `json:"xprank,string"`
	Class            string `json:"class"`
	City             string
	CityRank         int
	House            string
	HouseRank        int
	Order            *string
	BoundCredits     int    `json:"boundcredits,string"`
	UnboundCredits   int    `json:"unboundcredits,string"`
	Lessons          int    `json:"lessons,string"`
	ExplorerRank     string `json:"explorerrank"`
	MayanCrowns      int    `json:"mayancrowns,string"`
	BoundMayanCrowns int    `json:"boundmayancrowns,string"`
	Gold             int    `json:"gold,string"`
	Bank             int    `json:"bank,string"`
	UnreadNews       int    `json:"unread_news,string"`
	UnreadMessages   int    `json:"unread_msgs,string"`
	Target           *string
	Gender           int // ISO/IEC 5218
}

func (msg CharStatus) Hydrate(data []byte) (Message, error) {
	type CharStatusAlias CharStatus
	var child struct {
		CharStatusAlias
		CLevel  string `json:"level"`
		CCity   string `json:"city"`
		CHouse  string `json:"house"`
		COrder  string `json:"order"`
		CTarget string `json:"target"`
		CGender string `json:"gender"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return msg, err
	}

	msg = (CharStatus)(child.CharStatusAlias)

	if level, rank := splitLevelRank(child.CLevel); level == 0 && rank == 0 {
		return msg, fmt.Errorf("failed parsing level: %w", err)
	} else {
		msg.Level = level
		msg.XP = rank
	}

	if city, rank := splitRank(child.CCity); city == "" && rank == 0 {
		return msg, fmt.Errorf("failed parsing city: %w", err)
	} else {
		msg.City = city
		msg.CityRank = rank
	}

	if house, rank := splitRank(child.CHouse); house == "" && rank == 0 {
		return msg, fmt.Errorf("failed parsing city: %w", err)
	} else {
		msg.House = house
		msg.HouseRank = rank
	}

	if child.COrder != "(None)" {
		msg.Order = gox.NewString(child.COrder)
	}

	if child.CTarget != "None" {
		msg.Target = gox.NewString(child.CTarget)
	}

	if child.CGender == "male" {
		msg.Gender = 1
	} else if child.CGender == "female" {
		msg.Gender = 2
	} else {
		msg.Gender = 9
	}

	return msg, nil
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

func (msg CharVitals) Hydrate(data []byte) (Message, error) {
	type CharVitalsAlias CharVitals
	var child struct {
		CharVitalsAlias
		CBal  *string `json:"bal"`
		CEq   *string `json:"eq"`
		CVote *string `json:"vote"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return msg, err
	}

	msg = (CharVitals)(child.CharVitalsAlias)
	if child.CBal != nil {
		msg.Bal = gox.NewBool(*child.CBal == "1")
	}
	if child.CEq != nil {
		msg.Eq = gox.NewBool(*child.CEq == "1")
	}
	if child.CVote != nil {
		msg.Vote = gox.NewBool(*child.CVote == "1")
	}

	return msg, nil
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
