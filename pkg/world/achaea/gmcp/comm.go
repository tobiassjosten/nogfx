package gmcp

var (
	_ ClientServerMessage = &CommChannelPlayers{}
)

// CommChannelPlayers is both a client-sent and server-sent GMCP message, to
// either request data or lists players and which channels (if any) they share
// with the player's character.
type CommChannelPlayers struct{}

// Hydrate populates the message with data.
func (msg CommChannelPlayers) Hydrate(_ []byte) (ServerMessage, error) {
	// @todo Implement this.
	return msg, nil
}

// String is the message's string representation.
func (msg CommChannelPlayers) String() string {
	return "Comm.Channel.Players"
}
