package agmcp

import (
	"encoding/json"
	"fmt"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

var (
	_ gmcp.ClientMessage = &CoreSupportsAdd{}
	_ gmcp.ClientMessage = &CoreSupportsRemove{}
	_ gmcp.ClientMessage = &CoreSupportsSet{}
)

// CoreSupports is a list of potentially supported modules.
type CoreSupports struct {
	gmcp.CoreSupports
	IRERift   *int
	IRETarget *int
}

// Strings transforms CoreSupports to a list of strings.
func (msg CoreSupports) Strings() []string {
	list := msg.CoreSupports.Strings()
	if msg.IRERift != nil {
		list = append(list, fmt.Sprintf("IRE.Rift %d", *msg.IRERift))
	}
	if msg.IRETarget != nil {
		list = append(list, fmt.Sprintf("IRE.Target %d", *msg.IRETarget))
	}

	return list
}

// String is the message's string representation.
func (msg CoreSupports) String() string {
	data, _ := json.Marshal(msg.Strings())
	return string(data)
}

// CoreSupportsSet is a client-sent GMCP message containing supported modules.
type CoreSupportsSet struct {
	CoreSupports
}

// String is the message's string representation.
func (msg CoreSupportsSet) String() string {
	return fmt.Sprintf("Core.Supports.Set %s", msg.CoreSupports)
}

// CoreSupportsAdd is a client-sent GMCP message adding supported modules.
type CoreSupportsAdd struct {
	CoreSupports
}

// String is the message's string representation.
func (msg CoreSupportsAdd) String() string {
	return fmt.Sprintf("Core.Supports.Add %s", msg.CoreSupports)
}

// CoreSupportsRemove is a client-sent GMCP message removing supported modules.
type CoreSupportsRemove struct {
	CoreSupports
}

// String is the message's string representation.
func (msg CoreSupportsRemove) String() string {
	return fmt.Sprintf("Core.Supports.Remove %s", msg.CoreSupports)
}
