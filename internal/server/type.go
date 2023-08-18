package server

import (
	"git.sr.ht/~mpldr/uniview/internal/mansion"
	"git.sr.ht/~mpldr/uniview/protocol"
)

type Server struct {
	protocol.UnimplementedUniViewServer
	Rooms *mansion.Mansion
}
