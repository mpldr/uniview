package mpv

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"bufio"
	"encoding/json"
	"io"
	"log/slog"
	"math/rand"
	"time"
)

func (p *MPV) monitor() {
	defer p.conn.Close()
	defer close(p.notifySeek)
	defer close(p.notifyPause)

	rd := bufio.NewReader(p.conn)
	for {
		// message, err := bufio.NewReader(p.conn).ReadBytes('\n')
		message, err := rd.ReadBytes('\n')
		if err == io.EOF {
			p.dead.Store(true)
			p.quitchan <- struct{}{}
			return
		}
		var res response
		err = json.Unmarshal(message, &res)
		slog.Debug("received message", "message", message)
		if err != nil {
			slog.Warn("received non-understood message", "message", message)
			continue
		}

		switch {
		case res.Event == "seek":
			if p.dropSeek.Load() {
				slog.Debug("dropped seek")
				p.dropSeek.Store(false)
				continue
			}
			select {
			case p.notifySeekInternal <- struct{}{}:
			default:
			}
		case res.Event == "idle":
			select {
			case p.notifyIdle <- struct{}{}:
			default:
			}
		case res.Event == "end-file":
			if res.Reason == "quit" {
				p.dead.Store(true)
				p.quitchan <- struct{}{}
				return
			}
		case res.Event == "file-loaded":
			select {
			case <-p.playerReady:
			default:
				close(p.playerReady)
			}
		case res.RequestID != 0:
			p.responsesMtx.Lock()
			p.responses[res.RequestID] = res
			p.responsesMtx.Unlock()
		}
	}
}

func (p *MPV) pollPause() {
	for {
		<-time.After(50 * time.Millisecond)
		req := rand.Int()
		p.send(command{
			Command:   []any{"get_property", "pause"},
			RequestID: req,
		})
		res := p.getResponse(req)
		if pause, ok := res.Data.(bool); ok {
			if pause != p.pauseState {
				p.pauseState = pause
				if p.dropPause.Load() {
					p.dropPause.Store(false)
					continue
				}
				select {
				case p.notifyPause <- p.pauseState:
				default:
				}
			}
		}
	}
}

func (p *MPV) handleSeekEvents() {
	for range p.notifySeekInternal {
		pos, err := p.GetPlaybackPos()
		if err == nil {
			select {
			case p.notifySeek <- pos:
			default:
			}
		}
	}
}

type response struct {
	Error     string `json:"error"`
	Event     string `json:"event"`
	ID        int    `json:"id"`
	Data      any    `json:"data"`
	Name      string `json:"name"`
	RequestID int    `json:"request_id"`
	Reason    string `json:"reason"`
}
