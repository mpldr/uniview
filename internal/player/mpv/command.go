package mpv

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

func (p *MPV) send(cmd command) {
	if p.dead.Load() {
		return
	}
	cmdJSON, _ := json.Marshal(cmd)
	_, err := fmt.Fprintf(p.conn, "%s\n", cmdJSON)
	slog.Debug("sent command", "command", cmdJSON, "error", err)
}

func (p *MPV) getResponse(id int) response {
	for {
		p.responsesMtx.RLock()
		if res, exists := p.responses[id]; exists {
			p.responsesMtx.RUnlock()
			p.responsesMtx.Lock()
			delete(p.responses, id)
			p.responsesMtx.Unlock()

			return res
		}
		p.responsesMtx.RUnlock()
	}
}

type command struct {
	Command   []any `json:"command"`
	RequestID int   `json:"request_id,omitempty"`
}
