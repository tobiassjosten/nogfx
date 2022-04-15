package gmcp

var (
	_ ClientMessage = &IRERiftRequest{}
)

// IRERiftRequest is a client-sent GMCP message to request a list of items in
// the player's inventory.
type IRERiftRequest struct{}

// String is the message's string representation.
func (msg IRERiftRequest) String() string {
	return "IRE.Rift.Request"
}
