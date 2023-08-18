package mansion

import (
	"context"
	"sync"

	"git.sr.ht/~mpldr/uniview/protocol"
)

type room struct {
	ctx           context.Context
	clientFeed    []*client
	clientFeedMtx sync.Mutex
}

func newRoom(ctx context.Context) *room {
	r := &room{
		ctx: ctx,
	}
	return r
}

func (r *room) Client(feed protocol.UniView_RoomServer) {
	r.clientFeedMtx.Lock()
	defer r.clientFeedMtx.Unlock()

	r.clientFeed = append(r.clientFeed, newClient(r.ctx, feed))
}

func (r *room) Broadcast(ev *protocol.RoomEvent) {
	r.clientFeedMtx.Lock()
	defer r.clientFeedMtx.Unlock()

	max := len(r.clientFeed)
	for i := 0; i < max; i++ {
		for r.clientFeed[i].Dead() {
			r.clientFeed[i] = r.clientFeed[len(r.clientFeed)-1]
			r.clientFeed[len(r.clientFeed)-1] = nil
			r.clientFeed = r.clientFeed[:len(r.clientFeed)-1]
			max--
		}
		r.clientFeed[i].Send(ev)
	}
}
