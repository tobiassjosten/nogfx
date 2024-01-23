package gmcp

// CharLogin is a GMCP message to log a character in.
type CharLogin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// ID is the prefix before the message's data.
func (*CharLogin) ID() string {
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

// CharName is a GMCP message containing basic information about
// the player's character. Only sent on login.
type CharName struct {
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
}

// ID is the prefix before the message's data.
func (*CharName) ID() string {
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

// CharStatusVars is a GMCP message listing character variables.
type CharStatusVars map[string]string

// ID is the prefix before the message's data.
func (*CharStatusVars) ID() string {
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
