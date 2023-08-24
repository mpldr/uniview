package server

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"git.sr.ht/~mpldr/uniview/internal/mansion"
	"git.sr.ht/~mpldr/uniview/protocol"
)

type Server struct {
	protocol.UnimplementedUniViewServer
	Rooms *mansion.Mansion
}
