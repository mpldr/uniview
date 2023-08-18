package mansion

import (
	"context"
	"sync"
	"sync/atomic"

	"git.sr.ht/~poldi1405/glog"
)

type Mansion struct {
	ctx          context.Context
	cancel       context.CancelFunc
	rooms        map[string]*room
	roomsMtx     sync.RWMutex
	shuttingDown atomic.Bool
}

func New() *Mansion {
	ctx, cancel := context.WithCancel(context.Background())
	return &Mansion{
		ctx:    ctx,
		cancel: cancel,
		rooms:  make(map[string]*room),
	}
}

func (m *Mansion) GetRoom(name string) *room {
	glog.Tracef("requested room: %s", name)
	m.roomsMtx.RLock()
	if r, exists := m.rooms[name]; exists {
		m.roomsMtx.RUnlock()
		return r
	}

	glog.Tracef("creating new room %q", name)
	m.roomsMtx.RUnlock()
	m.roomsMtx.Lock()
	r := newRoom(m.ctx)
	m.rooms[name] = r
	m.roomsMtx.Unlock()

	return r
}

func (m *Mansion) Close() {
	glog.Debug("closing mansion and evicting tenants")
	m.cancel()
}
