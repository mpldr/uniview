package mansion

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"git.sr.ht/~mpldr/uniview/protocol"
)

type Mansion struct {
	ctx      context.Context
	cancel   context.CancelFunc
	rooms    map[string]*room
	roomsMtx sync.RWMutex
	clientID atomic.Uint64
	shutDown chan struct{}
}

func New() *Mansion {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Mansion{
		ctx:      ctx,
		cancel:   cancel,
		rooms:    make(map[string]*room),
		shutDown: make(chan struct{}),
	}
	go m.housekeeping()
	return m
}

func (m *Mansion) housekeeping() {
	for {
		select {
		case <-m.ctx.Done():
			m.roomsMtx.Lock()
			slog.Debug("evacuating rooms", "count", len(m.rooms))
			var wg sync.WaitGroup
			for n, r := range m.rooms {
				wg.Add(1)
				go func(r *room, name string) {
					defer wg.Done()
					defer slog.Debug("cleared room", "id", name)

					r.Broadcast(&protocol.RoomEvent{
						Type: protocol.EventType_EVENT_SERVER_CLOSE,
					}, 0)
				}(r, n)
			}
			wg.Wait()
			close(m.shutDown)
			return
		case <-time.After(5 * time.Minute):
			slog.Debug("running housekeeper")
			m.roomsMtx.RLock()
			var wg sync.WaitGroup
			for n, r := range m.rooms {
				wg.Add(1)
				go func(r *room, name string) {
					defer wg.Done()
					defer slog.Debug("pinging room", "id", name)

					r.Broadcast(&protocol.RoomEvent{
						Type: protocol.EventType_EVENT_SERVER_PING,
					}, 0)
				}(r, n)
			}
			wg.Wait()
			slog.Debug("finished pinging")
			m.roomsMtx.RUnlock()

			m.roomsMtx.Lock()
			for n, r := range m.rooms {
				if len(r.clientFeed) == 0 {
					slog.Debug("deleted room", "id", n)
					delete(m.rooms, n)
				}
			}
			slog.Debug("housekeeper finished")
			m.roomsMtx.Unlock()
		}
	}
}

func (m *Mansion) GetRoom(name string) (*room, uint64) {
	slog.Debug("requested room", "id", name)
	m.roomsMtx.RLock()
	if r, exists := m.rooms[name]; exists {
		m.roomsMtx.RUnlock()
		id := m.clientID.Add(1)
		return r, id
	}

	slog.Debug("creating new room", "id", name)
	m.roomsMtx.RUnlock()
	m.roomsMtx.Lock()
	r := newRoom(m.ctx)
	m.rooms[name] = r
	m.roomsMtx.Unlock()

	id := m.clientID.Add(1)
	return r, id
}

func (m *Mansion) Close() {
	slog.Debug("closing mansion and evicting tenants")
	m.cancel()
	select {
	case <-m.shutDown:
	case <-time.After(5 * time.Second):
	}
}

func (m *Mansion) Closing() bool {
	select {
	case <-m.ctx.Done():
		return true
	default:
		return false
	}
}
