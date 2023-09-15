package server

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"io"
	"log/slog"

	"git.sr.ht/~mpldr/uniview/internal/conman"
	"git.sr.ht/~mpldr/uniview/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (s *Server) Room(feed protocol.UniView_RoomServer) error {
	log := conman.GetLogger(feed.Context())
	log.Debug("new connection initialized. waiting for join event…")
	ev, err := feed.Recv()
	if err != nil {
		log.Warn("failed to receive initial message", "error", err)
		return status.Errorf(codes.Internal, "failed to receive join event: %v", err)
	}
	if ev.Type != protocol.EventType_EVENT_JOIN {
		log.Warn("received unexpected join event", "type", ev.Type)
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
	log = slog.With("grpc_client_id", id)
	log.Debug("client connected")
	room.Client(feed, id)
	defer room.Disconnect(id)
	if room.GetPosition() < 0 {
		log.Debug("unitialized room")
		if joinEv.Timestamp.AsDuration() < 0 {
			log.Debug("no timestamp provided")
			joinEv.Timestamp = durationpb.New(0)
		}
		log.Debug("initialized room", "timestamp", joinEv.Timestamp.AsDuration())
		room.SetPosition(joinEv.Timestamp.AsDuration())
		room.SetPause(false)
	} else {
		ev := &protocol.RoomEvent{
			Type: protocol.EventType_EVENT_PAUSE,
			Event: &protocol.RoomEvent_PauseEvent{
				PauseEvent: &protocol.PlayPause{
					Pause:     room.GetPause(),
					Timestamp: durationpb.New(room.GetPosition()),
				},
			},
		}
		log.Debug("sending initial state to client", "state", ev)
		feed.Send(ev)
	}

	for {
		ev, err = feed.Recv()
		switch {
		case err == io.EOF:
			log.Debug("closed connection")
			return nil
		case err != nil:
			select {
			case <-feed.Context().Done():
				return nil
			default:
				log.Error("failed to read value", "error", err)
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
			log.Debug("client disconnected")
			return nil
		}
		log.Debug("received event", "type", ev.Type)
		room.Broadcast(ev, id)
	}
}
