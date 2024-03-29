package ironrealms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// IRETargetSet is a GMCP message to either set or verify the setting of the
// in-game target variable.
type IRETargetSet struct {
	Target string
}

// ID is the prefix before the message's data.
func (*IRETargetSet) ID() string {
	return "IRE.Target.Set"
}

// Marshal converts the message to a string.
func (msg *IRETargetSet) Marshal() string {
	return fmt.Sprintf("IRE.Target.Set %q", msg.Target)
}

// Unmarshal populates the message with data.
func (msg *IRETargetSet) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	err := json.Unmarshal(data, &msg.Target)
	if err != nil {
		return err
	}

	return nil
}

// IRETargetInfo is a GMCP message with additional information about the
// current active server side target.
type IRETargetInfo struct {
	Identity    string `json:"id"`
	Health      int    `json:"-"`
	Description string `json:"short_desc"`
}

// ID is the prefix before the message's data.
func (*IRETargetInfo) ID() string {
	return "IRE.Target.Info"
}

// Marshal converts the message to a string.
func (msg *IRETargetInfo) Marshal() string {
	proxy := struct {
		*IRETargetInfo
		PHealth string `json:"hpperc"`
	}{
		IRETargetInfo: msg,
	}

	if msg.Health != 0 {
		proxy.PHealth = fmt.Sprintf("%d%%", msg.Health)
	}

	data, _ := json.Marshal(proxy)

	return fmt.Sprintf("IRE.Target.Info %s", string(data))
}

// Unmarshal populates the message with data.
func (msg *IRETargetInfo) Unmarshal(data []byte) error {
	data = bytes.TrimPrefix(data, []byte(msg.ID()+" "))

	proxy := struct {
		*IRETargetInfo
		CHealth string `json:"hpperc"`
	}{
		IRETargetInfo: msg,
	}

	err := json.Unmarshal(data, &proxy)
	if err != nil {
		return err
	}

	*msg = *proxy.IRETargetInfo

	if proxy.CHealth != "" {
		if proxy.CHealth[len(proxy.CHealth)-1] == '%' {
			proxy.CHealth = proxy.CHealth[:len(proxy.CHealth)-1]
		}

		health, err := strconv.Atoi(proxy.CHealth)
		if err != nil {
			return err
		}

		msg.Health = health
	}

	return nil
}
