package server

import (
	"errors"
	"fmt"
	"io"
	"math/rand"

	"git.sr.ht/~mpldr/uniview/protocol"
	"git.sr.ht/~poldi1405/glog"
)

func (s *Server) Room(feed protocol.UniView_RoomServer) error {
	glog.Trace("new connection initialized. waiting for join event…")
	ev, err := feed.Recv()
	if err != nil {
		glog.Warnf("failed to receive initial message: %v", err)
		return fmt.Errorf("failed to receive initial message: %w", err)
	}
	if ev.Type != protocol.EventType_EVENT_JOIN {
		glog.Warnf("received unexpected join event: %s", ev.Type)
		return fmt.Errorf("received unexpected join event: %s", ev.Type)
	}

	joinEv := ev.GetJoin()
	if joinEv == nil {
		return errors.New("missing join event")
	}

	room := s.Rooms.GetRoom(joinEv.Name)
	id := rand.Int63()
	glog.Debugf("client has been assigned id %d", id)
	room.Client(feed, id)

	for {
		ev, err = feed.Recv()
		switch {
		case err == io.EOF:
			glog.Debugf("closed connection")
			return io.EOF
		case err != nil:
			select {
			case <-feed.Context().Done():
				return nil
			default:
				glog.Errorf("feed: failed to read value: %v", err)
				return fmt.Errorf("error while receiving: %w", err)
			}
		}
		switch ev.Type {
		case protocol.EventType_EVENT_CLIENT_DISCONNECT:
			room.Disconnect(id)
			glog.Debugf("client %d disconnected", id)
			return nil
		}
		glog.Debugf("received %s from %d", ev.Type, id)
		room.Broadcast(ev, id)
	}
}
