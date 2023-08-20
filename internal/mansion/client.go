package mansion

import (
	"context"

	"git.sr.ht/~mpldr/uniview/protocol"
)

type client struct {
	ctx    context.Context
	cancel context.CancelFunc
	feed   protocol.UniView_RoomServer
	id     int64
}

func newClient(parent context.Context, feed protocol.UniView_RoomServer, id int64) *client {
	c := &client{
		feed: feed,
		id:   id,
	}
	c.ctx, c.cancel = context.WithCancel(parent)

	return c
}

func (c *client) Send(ev *protocol.RoomEvent) error {
	err := c.feed.Send(ev)
	if err != nil {
		c.cancel()
	}
	return err
}

func (c *client) Dead() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}