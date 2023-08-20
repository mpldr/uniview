package mansion

import (
	"context"
	"sync"
	"time"

	"git.sr.ht/~mpldr/uniview/protocol"
	"git.sr.ht/~poldi1405/glog"
)

type room struct {
	ctx           context.Context
	clientFeed    []*client
	clientFeedMtx sync.Mutex

	playbackStart time.Time
}

func newRoom(ctx context.Context) *room {
	r := &room{
		ctx: ctx,
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

	glog.Debugf("broadcasting message to %d clients", len(r.clientFeed))

	max := len(r.clientFeed)
	for i := 0; i < max; i++ {
		for r.clientFeed[i].Dead() {
			r.clientFeed[i] = r.clientFeed[len(r.clientFeed)-1]
			r.clientFeed[len(r.clientFeed)-1] = nil
			r.clientFeed = r.clientFeed[:len(r.clientFeed)-1]
			max--
			if i == len(r.clientFeed) {
				return
			}
		}
		if r.clientFeed[i].id == id {
			continue
		}
		err := r.clientFeed[i].Send(ev)
		if err != nil {
			glog.Warnf("failed to send to client %d: %v", i, err)
		}
	}
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
