package achaea

import (
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/world/achaea/agmcp"
)

// Character is the currently logged in character.
type Character struct {
	Name  string
	Title string

	Level int
	XP    int

	Balance     bool
	Equilibrium bool

	Health       int
	MaxHealth    int
	Mana         int
	MaxMana      int
	Endurance    int
	MaxEndurance int
	Willpower    int
	MaxWillpower int
}

// FromCharName updates the character from a Char.Name GMCP message.
func (c *Character) FromCharName(msg gmcp.CharName) {
	c.Name = msg.Name
	c.Title = msg.Fullname
}

// FromCharStatus updates the character from a Char.Status GMCP message.
func (c *Character) FromCharStatus(msg agmcp.CharStatus) {
	if msg.Level != nil {
		c.Level = *msg.Level
	}
}

// FromCharVitals updates the character from a Char.Vitals GMCP message.
func (c *Character) FromCharVitals(msg agmcp.CharVitals) {
	c.XP = msg.NL

	c.Balance = msg.Bal
	c.Equilibrium = msg.Eq

	c.Health = msg.HP
	c.MaxHealth = msg.MaxHP
	c.Mana = msg.MP
	c.MaxMana = msg.MaxMP
	c.Endurance = msg.EP
	c.MaxEndurance = msg.MaxEP
	c.Willpower = msg.WP
	c.MaxWillpower = msg.MaxWP
}
