package mansion

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"git.sr.ht/~mpldr/uniview/protocol"
)

type room struct {
	ctx           context.Context
	clientFeed    []*client
	clientFeedMtx sync.Mutex

	playbackStart  time.Time
	playbackPos    time.Duration
	playbackPaused bool

	password string
}

func newRoom(ctx context.Context) *room {
	r := &room{
		ctx:            ctx,
		playbackPos:    -1,
		playbackPaused: true,
	}
	return r
}

func (r *room) Client(feed protocol.UniView_RoomServer, id uint64) {
	r.clientFeedMtx.Lock()
	defer r.clientFeedMtx.Unlock()

	r.clientFeed = append(r.clientFeed, newClient(r.ctx, feed, id))
}

func (r *room) Broadcast(ev *protocol.RoomEvent, id uint64) {
	r.clientFeedMtx.Lock()
	defer r.clientFeedMtx.Unlock()

	slog.Debug("broadcasting message to clients", "count", len(r.clientFeed))

	max := len(r.clientFeed)
	for i := 0; i < max; i++ {
		for r.clientFeed[i].Dead() && ev.Type != protocol.EventType_EVENT_SERVER_CLOSE {
			slog.Debug("removing client", "client_id", r.clientFeed[i].id)
			r.clientFeed[i] = r.clientFeed[len(r.clientFeed)-1]
			r.clientFeed[len(r.clientFeed)-1] = nil
			r.clientFeed = r.clientFeed[:len(r.clientFeed)-1]
			max--
			if i == len(r.clientFeed) {
				return
			}
		}
		if r.clientFeed[i].id == id {
			slog.Debug("skipping client sender", "client_id", id)
			continue
		}
		slog.Debug("sending event to client", "event_type", ev.Type, "client_id", r.clientFeed[i].id)
		err := r.clientFeed[i].Send(ev)
		if err != nil {
			slog.Warn("failed to send to client", "client_id", r.clientFeed[i].id, "error", err)
		}
	}
}

func (r *room) SetPosition(pos time.Duration) {
	r.playbackStart = time.Now().Add(-1 * pos)
	r.playbackPos = pos
}

func (r *room) GetPosition() time.Duration {
	if r.playbackPaused {
		return r.playbackPos
	}
	return time.Since(r.playbackStart)
}

func (r *room) SetPause(pause bool) {
	r.playbackPaused = pause
}

func (r *room) GetPause() bool {
	return r.playbackPaused
}

func (r *room) Disconnect(id uint64) {
	r.clientFeedMtx.Lock()
	defer r.clientFeedMtx.Unlock()

	for k, v := range r.clientFeed {
		if v.id == id {
			r.clientFeed[k] = r.clientFeed[len(r.clientFeed)-1]
			r.clientFeed[len(r.clientFeed)-1] = nil
			r.clientFeed = r.clientFeed[:len(r.clientFeed)-1]
			return
		}
	}
}
