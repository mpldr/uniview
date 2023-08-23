package server

import (
	"io"

	"git.sr.ht/~mpldr/uniview/protocol"
	"git.sr.ht/~poldi1405/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (s *Server) Room(feed protocol.UniView_RoomServer) error {
	glog.Trace("new connection initialized. waiting for join event…")
	ev, err := feed.Recv()
	if err != nil {
		glog.Warnf("failed to receive initial message: %v", err)
		return status.Errorf(codes.Internal, "failed to receive join event: %v", err)
	}
	if ev.Type != protocol.EventType_EVENT_JOIN {
		glog.Warnf("received unexpected join event: %s", ev.Type)
		return status.Errorf(codes.FailedPrecondition, "received unexpected join event: %s", ev.Type)
	}

	joinEv := ev.GetJoin()
	if joinEv == nil {
		return status.Error(codes.FailedPrecondition, "incomplete join event received")
	}

	if s.Rooms.Closing() {
		return status.Error(codes.Unavailable, "the server is shutting down")
	}

	room, id := s.Rooms.GetRoom(joinEv.Name)
	glog.Debugf("client has been assigned id %d", id)
	room.Client(feed, id)
	if room.GetPosition() < 0 {
		if joinEv.Timestamp.AsDuration() < 0 {
			joinEv.Timestamp = durationpb.New(0)
		}
		room.SetPosition(joinEv.Timestamp.AsDuration())
		room.SetPause(false)
	} else {
		feed.Send(&protocol.RoomEvent{
			Type: protocol.EventType_EVENT_PAUSE,
			Event: &protocol.RoomEvent_PauseEvent{
				PauseEvent: &protocol.PlayPause{
					Pause:     room.GetPause(),
					Timestamp: durationpb.New(room.GetPosition()),
				},
			},
		})
	}

	for {
		ev, err = feed.Recv()
		switch {
		case err == io.EOF:
			glog.Debugf("closed connection")
			return nil
		case err != nil:
			select {
			case <-feed.Context().Done():
				return nil
			default:
				glog.Errorf("feed: failed to read value: %v", err)
				return status.Errorf(codes.Internal, "failed to receive event: %v", err)
			}
		}
		switch ev.Type {
		case protocol.EventType_EVENT_PAUSE:
			if ev.GetPauseEvent().GetTimestamp().AsDuration() < 0 {
				break
			}
			room.SetPause(ev.GetPauseEvent().GetPause())
			room.SetPosition(ev.GetJoin().GetTimestamp().AsDuration())
		case protocol.EventType_EVENT_JUMP:
			room.SetPosition(ev.GetJumpEvent().GetTimestamp().AsDuration())
		case protocol.EventType_EVENT_CLIENT_DISCONNECT:
			room.Disconnect(id)
			glog.Debugf("client %d disconnected", id)
			return nil
		}
		glog.Debugf("received %s from %d", ev.Type, id)
		room.Broadcast(ev, id)
	}
}
