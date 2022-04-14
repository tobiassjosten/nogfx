package gmcp

var (
	_ Message = &IRERiftRequest{}
)

// IRERiftRequest is a client-sent GMCP message to request a list of items in
// the player's inventory.
type IRERiftRequest struct{}

// Hydrate populates the message with data.
func (msg IRERiftRequest) Hydrate(_ []byte) (Message, error) {
	// This is client-side only, so we'll never have to hydrate it.
	return nil, nil
}

// String is the message's string representation.
func (msg IRERiftRequest) String() string {
	return "IRE.Rift.Request"
}
