package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/icza/gox/gox"
)

// @todo fixa doc comments "server-sent GMCP message…"

// CharLogin is a client-sent GMCP message to log a character in.
type CharLogin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// ID is the prefix before the message's data.
func (msg *CharLogin) ID() string {
	return "Char.Login"
}

// Marshal converts the message to a string.
func (msg *CharLogin) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharLogin) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharName is a server-sent GMCP message containing basic information about
// the player's character. Only sent on login.
type CharName struct {
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
}

// ID is the prefix before the message's data.
func (msg *CharName) ID() string {
	return "Char.Name"
}

// Marshal converts the message to a string.
func (msg *CharName) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharName) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharStatus is a server-sent GMCP message containing character values. The
// initial message sent contains all values but subsequent messages only carry
// changes, with omitted properties assumed unchanged.
type CharStatus struct {
	Fullname *string  `json:"fullname,omitempty"`
	Gender   *string  `json:"gender,omitempty"`
	Level    *float64 `json:"-"`
	Name     *string  `json:"name,omitempty"`
	Race     *string  `json:"race,omitempty"`
}

// ID is the prefix before the message's data.
func (msg *CharStatus) ID() string {
	return "Char.Status"
}

func (msg *CharStatus) MarshalLevel() *string {
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

// Marshal converts the message to a string.
func (msg *CharStatus) Marshal() string {
	proxy := struct {
		*CharStatus
		PLevel *string `json:"level,omitempty"`
	}{
		CharStatus: msg,
	}

	proxy.PLevel = msg.MarshalLevel()

	data, _ := json.Marshal(proxy)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharStatus) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	var proxy struct {
		*CharStatus
		PLevel *string `json:"level"`
	}

	err := json.Unmarshal(data, &proxy)
	if err != nil {
		return err
	}

	*msg = CharStatus{}
	if proxy.CharStatus != nil {
		*msg = (CharStatus)(*proxy.CharStatus)
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

	return nil
}

// CharStatusVars is a server-sent GMCP message listing character variables.
type CharStatusVars map[string]string

// ID is the prefix before the message's data.
func (msg *CharStatusVars) ID() string {
	return "Char.StatusVars"
}

// Marshal converts the message to a string.
func (msg *CharStatusVars) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharStatusVars) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharVitals is a server-sent GMCP message containing character attributes.
type CharVitals struct {
	HP    int      `json:"hp,string"`
	MaxHP int      `json:"maxhp,string"`
	MP    int      `json:"mp,string"`
	MaxMP int      `json:"maxmp,string"`
	Stats []string `json:"charstats"`
}

// ID is the prefix before the message's data.
func (msg *CharVitals) ID() string {
	return "Char.Vitals"
}

// Marshal converts the message to a string.
func (msg *CharVitals) Marshal() string {
	proxy := struct {
		*CharVitals
		Stats []string `json:"charstats"`
	}{
		CharVitals: msg,
	}

	proxy.Stats = msg.Stats
	if msg.Stats == nil {
		proxy.Stats = []string{}
	}

	data, _ := json.Marshal(proxy)
	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharVitals) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}

	return nil
}
