package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/icza/gox/gox"
)

type Message interface {
	String() string
}

func Parse(command []byte) (Message, error) {
	parts := bytes.SplitN(command, []byte{' '}, 2)

	var message Message

	switch string(parts[0]) {
	case "Char.Items.Inv":
		message = CharItemsInv{}

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

func (msg CharName) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		data = []byte("{}")
	}

	return fmt.Sprintf("Char.Name %s", data)
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

	fmt.Printf("unmarshaled '%+v' (%T)\n", msg, msg)

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
