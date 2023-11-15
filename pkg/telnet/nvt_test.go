package telnet_test

import (
	"net"

	"github.com/tobiassjosten/nogfx/pkg/telnet"
)

// Verify interface fulfilments.
var _ net.Conn = &telnet.NVT{}
