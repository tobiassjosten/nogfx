package gmcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/icza/gox/gox"
)

// CharSkillsGet is a GMCP message to request skill information.
type CharSkillsGet struct {
	Group string `json:"group,omitempty"`
	Name  string `json:"name,omitempty"`
}

// ID is the prefix before the message's data.
func (msg *CharSkillsGet) ID() string {
	return "Char.Skills.Get"
}

// Marshal converts the message to a string.
func (msg *CharSkillsGet) Marshal() string {
	proxy := struct {
		*CharSkillsGet
		Name string `json:"name,omitempty"`
	}{msg, msg.Name}

	if msg.Group == "" {
		proxy.Name = ""
	}

	return Marshal(proxy)
}

// Unmarshal populates the message with data.
func (msg *CharSkillsGet) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

type charSkillsGroup struct {
	Name     string `json:"name"`
	Rank     string `json:"rank"`
	Progress *int   `json:"-"`
}

// CharSkillsGroups is a GMCP message listing groups of skills available to
// the character.
type CharSkillsGroups []charSkillsGroup

// ID is the prefix before the message's data.
func (msg *CharSkillsGroups) ID() string {
	return "Char.Skills.Groups"
}

// Marshal converts the message to a string.
func (msg *CharSkillsGroups) Marshal() string {
	type Proxy struct {
		charSkillsGroup
		Rank string `json:"rank"`
	}

	proxies := []Proxy{}

	for _, group := range *msg {
		proxy := Proxy{group, group.Rank}

		if group.Progress != nil {
			proxy.Rank = fmt.Sprintf("%s (%d%%)", group.Rank, *group.Progress)
		}

		proxies = append(proxies, proxy)
	}

	data, _ := json.Marshal(proxies)

	return fmt.Sprintf("%s %s", msg.ID(), string(data))
}

// Unmarshal populates the message with data.
func (msg *CharSkillsGroups) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(bytes.TrimPrefix(data, []byte(msg.ID())))

	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}

	for i, group := range *msg {
		parts := strings.SplitN(group.Rank, " ", 2)

		// Documentation says that `rank` is formatted as "Adept (1%)"
		// but some games (like Achaea) only show "Adept".
		if len(parts) == 1 {
			continue
		}

		progressStr := strings.Trim(parts[1], "(%)")
		progress, err := strconv.Atoi(progressStr)

		if err != nil {
			return fmt.Errorf("failed parsing rank progress: %w", err)
		}

		(*msg)[i].Rank = parts[0]
		(*msg)[i].Progress = gox.NewInt(progress)
	}

	return nil
}

// CharSkillsList is a GMCP message listing skills within a group available to
// the character.
type CharSkillsList struct {
	Group        string   `json:"group"`
	List         []string `json:"list"`
	Descriptions []string `json:"descs"`
}

// ID is the prefix before the message's data.
func (msg *CharSkillsList) ID() string {
	return "Char.Skills.List"
}

// Marshal converts the message to a string.
func (msg *CharSkillsList) Marshal() string {
	proxy := struct {
		*CharSkillsList
		List         []string `json:"list"`
		Descriptions []string `json:"descs"`
	}{msg, msg.List, msg.Descriptions}

	if msg.List == nil {
		proxy.List = []string{}
	}

	if msg.Descriptions == nil {
		proxy.Descriptions = []string{}
	}

	return Marshal(proxy)
}

// Unmarshal populates the message with data.
func (msg *CharSkillsList) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}

// CharSkillsInfo is a GMCP message detailing information about a single skill.
type CharSkillsInfo struct {
	Group string `json:"group"`
	Skill string `json:"skill"`
	Info  string `json:"info"`
}

// ID is the prefix before the message's data.
func (msg *CharSkillsInfo) ID() string {
	return "Char.Skills.Info"
}

// Marshal converts the message to a string.
func (msg *CharSkillsInfo) Marshal() string {
	return Marshal(msg)
}

// Unmarshal populates the message with data.
func (msg *CharSkillsInfo) Unmarshal(data []byte) error {
	return Unmarshal(data, msg)
}
