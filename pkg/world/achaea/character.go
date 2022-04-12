package achaea

import "github.com/tobiassjosten/nogfx/pkg/world/achaea/gmcp"

type Character struct {
	Name string

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

func (c *Character) fromCharName(msg gmcp.CharName) {
	c.Name = msg.Name
}

func (c *Character) fromCharVitals(msg gmcp.CharVitals) {
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
