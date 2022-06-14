package achaea

import (
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"
	igmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
)

// Target represents who or what is being targeted for skills and attacks.
type Target struct {
	*pkg.Target
	client   pkg.Client
	isPlayer bool
}

// NewTarget creates a new target object with the given client.
func NewTarget(client pkg.Client) *Target {
	target := &Target{client: client}
	target.Target = pkg.NewTarget(target.Set)
	return target
}

// PkgTarget converts our game-specific Target to the general pkg struct.
func (tgt *Target) PkgTarget() *pkg.Target {
	return tgt.Target
}

// Set executes the actual target change.
func (tgt *Target) Set(name string, _ *pkg.Target) {
	if tgt.isPlayer {
		return
	}

	if name == "" {
		tgt.client.Send([]byte("settarget none"))
		return
	}

	tgt.client.Send([]byte("settarget " + name))
}

// FromRoomInfo handles targeting when moving between rooms (areas, in effect).
func (tgt *Target) FromRoomInfo(msg *gmcp.RoomInfo) {
	room := navigation.RoomFromGMCP(msg)
	npcs := tgt.npcs(room)
	tgt.Target.SetCandidates(npcs)
}

// FromCharItemsList builds the list of NPCs in the room and retargets.
func (tgt *Target) FromCharItemsList(msg *gmcp.CharItemsList) {
	if msg.Location != "room" {
		return
	}

	present := []string{}
	for _, item := range msg.Items {
		present = append(present, item.Name)
	}
	tgt.Target.SetPresent(present)
}

// FromCharItemsAdd adds an NPC to the room list and retargets.
func (tgt *Target) FromCharItemsAdd(msg *gmcp.CharItemsAdd) {
	if msg.Location != "room" {
		return
	}

	tgt.Target.AddPresent(msg.Item.Name)
}

// FromCharItemsRemove removes an NPC to the room list and retargets.
func (tgt *Target) FromCharItemsRemove(msg *gmcp.CharItemsRemove) {
	if msg.Location != "room" {
		return
	}

	tgt.Target.RemovePresent(msg.Item.Name)
}

// FromCharStatus updates the current target.
func (tgt *Target) FromCharStatus(msg *agmcp.CharStatus) {
	if msg.Target != nil {
		tgt.Name = strings.ToLower(*msg.Target)
	}
}

// FromIRETargetSet updates the player status of the current target.
func (tgt *Target) FromIRETargetSet(msg *igmcp.IRETargetSet) {
	// This message works so inconsistenyly that we can only rely
	// on it for knowing that non-numbers equals a player.
	if msg.Target != "" {
		_, err := strconv.Atoi(msg.Target)
		tgt.isPlayer = err != nil
	}

	if msg.Target == "" || tgt.isPlayer {
		tgt.Health = -1
	}
}

// FromIRETargetInfo updates the current NPC-target's health.
func (tgt *Target) FromIRETargetInfo(msg *igmcp.IRETargetInfo) {
	tgt.Health = msg.Health
}

func (tgt *Target) npcs(room *navigation.Room) []string {
	if room == nil || room.Area == nil {
		return []string{}
	}

	// An important property of these lists is their order of importance,
	// where the most dangerous NPC is first and the rest in falling order.
	switch room.Area.ID {
	case 137:
		return []string{"shaman", "warrior", "manticore", "villager"}
	}

	return []string{}
}
