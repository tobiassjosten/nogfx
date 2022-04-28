package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/icza/gox/gox"
)

var (
	_ ClientMessage = &CharItemsContents{}
	_ ClientMessage = &CharItemsInv{}
	_ ClientMessage = &CharItemsRoom{}
	_ ClientMessage = &CharLogin{}
	_ ClientMessage = &CharSkillsGet{}

	_ ServerMessage = &CharAfflictionsAdd{}
	_ ServerMessage = &CharAfflictionsList{}
	_ ServerMessage = &CharAfflictionsRemove{}
	_ ServerMessage = &CharDefencesAdd{}
	_ ServerMessage = &CharDefencesList{}
	_ ServerMessage = &CharDefencesRemove{}
	_ ServerMessage = &CharItemsAdd{}
	_ ServerMessage = &CharItemsList{}
	_ ServerMessage = &CharItemsRemove{}
	_ ServerMessage = &CharItemsUpdate{}
	_ ServerMessage = &CharItemsUpdate{}
	_ ServerMessage = &CharName{}
	_ ServerMessage = &CharSkillsGroups{}
	_ ServerMessage = &CharSkillsInfo{}
	_ ServerMessage = &CharSkillsList{}
	_ ServerMessage = &CharStatusVars{}
	_ ServerMessage = &CharStatus{}
	_ ServerMessage = &CharVitals{}
)

// CharAffliction is an affliction ailing a character.
type CharAffliction struct {
	Name        string `json:"name"`
	Cure        string `json:"cure"`
	Description string `json:"desc"`
}

// CharAfflictionsList is a server-sent GMCP message listing current character
// afflictions
type CharAfflictionsList []CharAffliction

// Hydrate populates the message with data.
func (msg CharAfflictionsList) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharAfflictionsAdd is a server-sent GMCP message listing current character
// afflictions
type CharAfflictionsAdd CharAffliction

// Hydrate populates the message with data.
func (msg CharAfflictionsAdd) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharAfflictionsRemove is a server-sent GMCP message listing current character
// afflictions
type CharAfflictionsRemove []CharAffliction

// Hydrate populates the message with data.
func (msg CharAfflictionsRemove) Hydrate(data []byte) (ServerMessage, error) {
	list := []string{}

	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		msg = append(msg, CharAffliction{Name: item})
	}

	return msg, nil
}

// CharDefence is a defence protecting a character.
type CharDefence struct {
	Name        string `json:"name"`
	Cure        string `json:"cure"`
	Description string `json:"desc"`
}

// CharDefencesList is a server-sent GMCP message listing current character
// afflictions
type CharDefencesList []CharDefence

// Hydrate populates the message with data.
func (msg CharDefencesList) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharDefencesAdd is a server-sent GMCP message listing current character
// afflictions
type CharDefencesAdd CharDefence

// Hydrate populates the message with data.
func (msg CharDefencesAdd) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharDefencesRemove is a server-sent GMCP message listing current character
// afflictions
type CharDefencesRemove []CharDefence

// Hydrate populates the message with data.
func (msg CharDefencesRemove) Hydrate(data []byte) (ServerMessage, error) {
	list := []string{}

	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		msg = append(msg, CharDefence{Name: item})
	}

	return msg, nil
}

// CharItemsContents is a client-sent GMCP message to request a list of items
// located inside another item.
type CharItemsContents struct {
	ID string
}

// String is the message's string representation.
func (msg CharItemsContents) String() string {
	return strings.TrimSpace(fmt.Sprintf("Char.Items.Contents %s", msg.ID))
}

// CharItemsInv is a client-sent GMCP message to request a list of items in the
// player's inventory.
type CharItemsInv struct{}

// String is the message's string representation.
func (msg CharItemsInv) String() string {
	return "Char.Items.Inv"
}

// CharItemsRoom is a client-sent GMCP message to request an updated list of
// items in the current room.
type CharItemsRoom struct {
	ID string
}

// String is the message's string representation.
func (msg CharItemsRoom) String() string {
	return fmt.Sprintf("Char.Items.Room")
}

// CharItem is an item within a player's inventory, the current room, or any
// other container.
type CharItem struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Attributes CharItemAttributes `json:"attrib"`
	Icon       string             `json:"icon"`
}

// CharItemAttributes is a set of flags denoting how to interact with an item.
type CharItemAttributes struct {
	Container    bool
	Dangerous    bool
	Dead         bool
	Edible       bool
	Fluid        bool
	Groupable    bool
	Monster      bool
	Riftable     bool
	Takeable     bool
	Wearable     bool
	WieldedLeft  bool
	WieldedRight bool
	Worn         bool
}

// UnmarshalJSON hydrates CharItemAttributes from a string.
func (as *CharItemAttributes) UnmarshalJSON(data []byte) error {
	for _, char := range bytes.Trim(data, `"`) {
		switch char {
		case 'c':
			as.Container = true

		case 'd':
			as.Dead = true

		case 'e':
			as.Edible = true

		case 'f':
			as.Fluid = true

		case 'g':
			as.Groupable = true

		case 'l':
			as.WieldedLeft = true

		case 'L':
			as.WieldedRight = true

		case 'm':
			as.Monster = true

		case 'r':
			as.Riftable = true

		case 't':
			as.Takeable = true

		case 'w':
			as.Worn = true

		case 'W':
			as.Wearable = true

		case 'x':
			as.Dangerous = true

		default:
			return fmt.Errorf("unknown attribute '%s'", string(char))
		}
	}

	return nil
}

// CharItemsList is a server-sent GMCP message listing items at the specified
// location.
type CharItemsList struct {
	Location string     `json:"location"`
	Items    []CharItem `json:"items"`
}

// Hydrate populates the message with data.
func (msg CharItemsList) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharItemsAdd is a server-sent GMCP message informing the client about an
// item being added to the specified location.
type CharItemsAdd struct {
	Location string   `json:"location"`
	Item     CharItem `json:"item"`
}

// Hydrate populates the message with data.
func (msg CharItemsAdd) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharItemsRemove is a server-sent GMCP message informing the client about an
// item being removed from the specified location.
type CharItemsRemove struct {
	Location string   `json:"location"`
	Item     CharItem `json:"item"`
}

// Hydrate populates the message with data.
func (msg CharItemsRemove) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharItemsUpdate is a server-sent GMCP message informing the client about an
// item being removed from the specified location.
type CharItemsUpdate struct {
	Location string   `json:"location"`
	Item     CharItem `json:"item"`
}

// Hydrate populates the message with data.
func (msg CharItemsUpdate) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharLogin is a client-sent GMCP message to log a character in.
type CharLogin struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

// String is the message's string representation.
func (msg CharLogin) String() string {
	data, _ := json.Marshal(msg)
	return fmt.Sprintf("Char.Login %s", data)
}

// CharSkillsGet is a client-sent GMCP message to request skill information.
type CharSkillsGet struct {
	Group string `json:"group,omitempty"`
	Name  string `json:"name,omitempty"`
}

// String is the message's string representation.
func (msg CharSkillsGet) String() string {
	if msg.Group == "" {
		msg.Name = ""
	}
	data, _ := json.Marshal(msg)
	return fmt.Sprintf("Char.Skills.Get %s", data)
}

// CharName is a server-sent GMCP message containing basic information about
// the player's character. Only sent on login.
type CharName struct {
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
}

// Hydrate populates the message with data.
func (msg CharName) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

type charSkillsGroup struct {
	Name     string
	Level    string
	Progress *int
}

// CharSkillsGroups is a server-sent GMCP message listing groups of skills
// available to the character.
type CharSkillsGroups []charSkillsGroup

// Hydrate populates the message with data.
func (msg CharSkillsGroups) Hydrate(data []byte) (ServerMessage, error) {
	var children []struct {
		Name string `json:"name"`
		Rank string `json:"rank"`
	}

	err := json.Unmarshal(data, &children)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		// Documentation says that `rank` is formatted as "Adept (1%)"
		// but at least Achaea instead seems to show only "Adept", so
		// we allow for both variants by just passing `rank` on as is.
		level, rank := splitRank(child.Rank)
		msg = append(msg, charSkillsGroup{
			Name:     child.Name,
			Level:    level,
			Progress: rank,
		})
	}

	return msg, nil
}

// CharSkillsList is a server-sent GMCP message listing skills within a group
// available to the character.
type CharSkillsList struct {
	Group        string   `json:"group"`
	List         []string `json:"list"`
	Descriptions []string `json:"descs"`
}

// Hydrate populates the message with data.
func (msg CharSkillsList) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharSkillsInfo is a server-sent GMCP message detailing information about a
// single skill.
type CharSkillsInfo struct {
	Group       string `json:"group"`
	Skill       string `json:"skill"`
	Information string `json:"info"`
}

// Hydrate populates the message with data.
func (msg CharSkillsInfo) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharStatusVars is a server-sent GMCP message listing character variables.
type CharStatusVars map[string]string

// Hydrate populates the message with data.
func (msg CharStatusVars) Hydrate(data []byte) (ServerMessage, error) {
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// CharStatus is a server-sent GMCP message containing character values. The
// initial message sent contains all values but subsequent messages only carry
// changes, with omitted properties assumed unchanged.
type CharStatus struct {
	Name     *string  `json:"name"`
	Fullname *string  `json:"fullname"`
	Age      *int     `json:"age,string"`
	Race     *string  `json:"race"`
	Level    *float64 `json:"-"`
	XP       *int     `json:"-"`
	Gender   *int     `json:"-"` // ISO/IEC 5218
}

// Hydrate populates the message with data.
func (msg CharStatus) Hydrate(data []byte) (ServerMessage, error) {
	type CharStatusAlias CharStatus
	var child struct {
		CharStatusAlias
		CLevel  *string `json:"level"`
		CGender *string `json:"gender"`
	}

	err := json.Unmarshal(data, &child)
	if err != nil {
		return nil, err
	}

	msg = (CharStatus)(child.CharStatusAlias)

	if child.CLevel != nil {
		lvl, rank := splitRank(*child.CLevel)
		level, _ := strconv.ParseFloat(lvl, 64)
		if rank == nil {
			return nil, fmt.Errorf(
				"failed parsing level '%s'", *child.CLevel,
			)
		}

		msg.Level = gox.NewFloat64(level)
		msg.XP = rank
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
func (msg CharVitals) Hydrate(data []byte) (ServerMessage, error) {
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

// CharVitalsStats is structured data extending CharVitals. Comparing different
// games, there is no overlap, and so we leave it here for reference only.
type CharVitalsStats struct {
}

// UnmarshalJSON hydrates CharVitalsStats from a list of unstructured strings.
func (stats *CharVitalsStats) UnmarshalJSON(data []byte) error {
	return nil
}
