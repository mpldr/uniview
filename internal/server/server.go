package server

import (
	"errors"
	"fmt"
	"io"

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
	room.Client(feed)

	for {
		ev, err = feed.Recv()
		switch {
		case err == io.EOF:
			glog.Debugf("closed connection")
			return io.EOF
		case err != nil:
			glog.Errorf("feed: failed to read value: %v", err)
			return fmt.Errorf("error while receiving: %w", err)
		}
		room.Broadcast(ev)
	}
}
