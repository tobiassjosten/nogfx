package achaea

import (
	"strconv"
	"strings"

	"github.com/tobiassjosten/nogfx/pkg"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/agmcp"
)

// Target represents who or what is being targeted for skills and attacks. We
// typically display only the General target to the player but use the Specific
// target internally for executing skills against.
type Target struct {
	client pkg.Client

	Name     string
	Health   int
	IsPlayer bool

	room     *navigation.Room
	roomNPCs []string
}

// FromRoom handles auto targeting when moving between bashing zones.
func (target *Target) FromRoomInfo(msg gmcp.RoomInfo) {
	current := navigation.RoomFromGMCP(msg)
	if current == nil || current.Area == nil {
		return
	}

	previous := target.room
	target.room = current

	if target.IsPlayer {
		return
	}

	if previous != nil && previous.Area != nil && previous.Area.ID == current.Area.ID {
		return
	}

	cNPCs, pNPCs := areaNPCs(current), areaNPCs(previous)

	changeable := target.Name == ""
	for _, npc := range pNPCs {
		if target.Name == npc {
			changeable = true
			break
		}
	}

	if len(cNPCs) == 0 {
		if changeable && target.Name != "" {
			target.client.Send([]byte("settarget none"))
		}
		return
	}

	if changeable && target.Name != cNPCs[0] {
		target.client.Send([]byte("settarget " + cNPCs[0]))
		return
	}
}

func (target *Target) FromCharItemsList(msg gmcp.CharItemsList) {
	if msg.Location != "room" {
		return
	}

	target.roomNPCs = []string{}
	for _, item := range msg.Items {
		for _, anpc := range target.areaNPCs() {
			if strings.Index(item.Name, anpc) > 0 {
				target.roomNPCs = append(target.roomNPCs, anpc)
				break
			}
		}
	}

	target.retarget()
}

func (target *Target) FromCharItemsAdd(msg gmcp.CharItemsAdd) {
	if msg.Location != "room" {
		return
	}

	for _, anpc := range target.areaNPCs() {
		if strings.Index(msg.Item.Name, anpc) > 0 {
			target.roomNPCs = append(target.roomNPCs, anpc)
			break
		}
	}

	target.retarget()
}

func (target *Target) FromCharItemsRemove(msg gmcp.CharItemsRemove) {
	if msg.Location != "room" {
		return
	}

	for i, rnpc := range target.roomNPCs {
		if strings.Index(msg.Item.Name, rnpc) >= 0 {
			target.roomNPCs = append(
				target.roomNPCs[:i],
				target.roomNPCs[i+1:]...,
			)
			break
		}
	}

	target.retarget()
}

func (target *Target) FromCharStatus(msg agmcp.CharStatus) {
	if msg.Target != nil {
		target.Name = strings.ToLower(*msg.Target)
	}
}

func (target *Target) FromIRETargetSet(msg gmcp.IRETargetSet) {
	// This message works so inconsistenyly that we can only rely
	// on it for knowing that non-numbers equalling a player.
	if msg.Target != "" {
		_, err := strconv.Atoi(msg.Target)
		target.IsPlayer = err != nil
	}
}

func (target *Target) FromIRETargetInfo(msg gmcp.IRETargetInfo) {
	if msg.Health > 0 {
		// @todo hur vet vi när vi ska ta bort vital? timeout bara?
		target.Health = msg.Health
	}
}

func (target *Target) retarget() {
	if target.IsPlayer {
		return
	}

	newTarget := target.roomTarget()

	if newTarget != "" && newTarget != target.Name {
		target.client.Send([]byte("settarget " + newTarget))
	}
}

func (target *Target) roomTarget() string {
	for _, anpc := range target.areaNPCs() {
		for _, rnpc := range target.roomNPCs {
			if anpc == rnpc {
				return anpc
			}
		}
	}

	return ""
}

func (target *Target) areaNPCs() []string {
	return areaNPCs(target.room)
}

func areaNPCs(room *navigation.Room) []string {
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
