package achaea

import (
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	agmcp "github.com/tobiassjosten/nogfx/pkg/gmcp/achaea"
)

// Character is the currently logged in character.
type Character struct {
	Name  string
	Title string

	Class string

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

	Bleed int
	Rage  int

	Karma    int
	Kai      int
	Stance   string
	Ferocity int
	Spec     string
}

// FromCharName updates the character from a Char.Name GMCP message.
func (c *Character) FromCharName(msg *gmcp.CharName) {
	c.Name = msg.Name
	c.Title = msg.Fullname
}

// FromCharStatus updates the character from a Char.Status GMCP message.
func (c *Character) FromCharStatus(msg *agmcp.CharStatus) {
	if msg.Name != nil {
		c.Name = *msg.Name
	}
	if msg.Fullname != nil {
		c.Title = *msg.Fullname
	}
	if msg.Class != nil {
		c.Class = *msg.Class
	}
	if msg.Level != nil {
		c.Level = int(*msg.Level)
	}
}

// FromCharVitals updates the character from a Char.Vitals GMCP message.
func (c *Character) FromCharVitals(msg *agmcp.CharVitals) {
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

	c.Bleed = msg.Stats.Bleed
	c.Rage = msg.Stats.Rage

	if msg.Stats.Ferocity != nil {
		c.Ferocity = *msg.Stats.Ferocity
	}
	if msg.Stats.Kai != nil {
		c.Kai = *msg.Stats.Kai
	}
	if msg.Stats.Karma != nil {
		c.Karma = *msg.Stats.Karma
	}
	if msg.Stats.Spec != nil {
		c.Spec = *msg.Stats.Spec
	}
	if msg.Stats.Stance != nil {
		c.Stance = *msg.Stats.Stance
	}
}
