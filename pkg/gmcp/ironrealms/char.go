package ironrealms

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"

	"github.com/icza/gox/gox"
)

// CharStatus is a server-sent GMCP message containing character values. The
// initial message sent contains all values but subsequent messages only carry
// changes, with omitted properties assumed unchanged.
type CharStatus struct {
	*gmcp.CharStatus
	Bank       *int    `json:"bank,string,omitempty"`
	City       *string `json:"-"`
	CityRank   *int    `json:"-"`
	Class      *string `json:"class,omitempty"`
	Gold       *int    `json:"gold,string,omitempty"`
	UnreadMsgs *int    `json:"unread_msgs,string,omitempty"`
	UnreadNews *int    `json:"unread_news,string,omitempty"`
}

// ID is the prefix before the message's data.
func (msg *CharStatus) ID() string {
	return "Char.Status"
}

func (msg *CharStatus) MarshalCity() *string {
	if msg.City == nil {
		return nil
	}

	if *msg.City == "" {
		// I know, it's stupid, but at least it's consistent.
		return gox.NewString("(None)")
	}

	city := *msg.City
	if msg.CityRank != nil {
		city = fmt.Sprintf("%s (%d)", city, *msg.CityRank)
	}

	return &city
}

// Marshal converts the message to a string.
func (msg *CharStatus) Marshal() string {
	proxy := struct {
		*CharStatus
		PCity  *string `json:"city,omitempty"`
		PLevel *string `json:"level,omitempty"`
	}{
		CharStatus: msg,
	}

	proxy.PCity = msg.MarshalCity()
	proxy.PLevel = msg.MarshalLevel()

	data, _ := json.Marshal(proxy)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharStatus) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	proxy := struct {
		*CharStatus
		PCity *string `json:"city"`
	}{
		CharStatus: &CharStatus{
			CharStatus: &gmcp.CharStatus{},
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

	if proxy.PCity != nil {
		if *proxy.PCity == "(None)" {
			msg.City = gox.NewString("")
		} else {
			city, rank := gmcp.SplitRankInt(*proxy.PCity)
			msg.City = gox.NewString(city)
			msg.CityRank = gox.NewInt(rank)
		}
	}

	return nil
}

// CharVitals is a server-sent GMCP message containing character attributes.
type CharVitals struct {
	Bal bool `json:"-"`
	Eq  bool `json:"-"`
	NL  int  `json:"nl,string"`
}

// ID is the prefix before the message's data.
func (msg *CharVitals) ID() string {
	return "Char.Vitals"
}

// Marshal converts the message to a string.
func (msg *CharVitals) Marshal() string {
	proxy := struct {
		*CharVitals
		PBal string `json:"bal"`
		PEq  string `json:"eq"`
	}{
		CharVitals: msg,
		PBal:       "0",
		PEq:        "0",
	}

	if msg.Bal {
		proxy.PBal = "1"
	}
	if msg.Eq {
		proxy.PEq = "1"
	}

	data, _ := json.Marshal(proxy)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharVitals) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	var proxy struct {
		*CharVitals
		PBal string `json:"bal"`
		PEq  string `json:"eq"`
	}

	err := json.Unmarshal(data, &proxy)
	if err != nil {
		return err
	}

	*msg = CharVitals{}
	if proxy.CharVitals != nil {
		*msg = (CharVitals)(*proxy.CharVitals)
	}

	msg.Bal = proxy.PBal == "1"
	msg.Eq = proxy.PEq == "1"

	return nil
}
