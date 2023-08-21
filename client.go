package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"git.sr.ht/~mpldr/uniview/internal/player"
	"git.sr.ht/~mpldr/uniview/internal/player/mpv"
	"git.sr.ht/~mpldr/uniview/protocol"
	"git.sr.ht/~poldi1405/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

func startClient() error {
	if len(os.Args) < 4 {
		fmt.Println("Usage: uniview [server] [room] [file]")
		os.Exit(2)
	}

	var p player.Interface
	var err error

	glog.Debug("starting player…")
	p, err = mpv.New()
	if err != nil {
		return fmt.Errorf("failed to start mpv: %w", err)
	}
	defer p.Close()

	glog.Debugf("loading file %q…", os.Args[3])
	err = p.LoadFile(os.Args[3])
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	glog.Debugf("connecting to remote…")
	gconn, err := grpc.Dial(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer gconn.Close()

	glog.Debug("requesting handle…")
	cl := protocol.NewUniViewClient(gconn)
	stream, err := cl.Room(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get room feed: %w", err)
	}

	pos, err := p.GetPlaybackPos()
	if err != nil {
		pos = -1
	}

	glog.Debugf("Joining room %q", os.Args[2])
	err = stream.Send(&protocol.RoomEvent{
		Type: protocol.EventType_EVENT_JOIN,
		Event: &protocol.RoomEvent_Join{
			Join: &protocol.RoomJoin{
				Name:      os.Args[2],
				Timestamp: durationpb.New(pos),
				Url:       os.Args[3],
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	ignoreSeek := &atomic.Bool{}

	glog.Debug("waiting for remote events…")
	go recvChanges(stream, p, ignoreSeek)

	glog.Debug("waiting for player events…")
	for {
		select {
		case timestamp := <-p.NotifySeek():
			if ignoreSeek.Load() {
				glog.Debug("seek ignored")
				continue
			}
			glog.Debugf("local: seek to %s detected", timestamp)
			stream.Send(&protocol.RoomEvent{
				Type: protocol.EventType_EVENT_JUMP,
				Event: &protocol.RoomEvent_JumpEvent{
					JumpEvent: &protocol.PlaybackJump{
						Timestamp: durationpb.New(timestamp),
					},
				},
			})
		case pause := <-p.NotifyPause():
			if ignoreSeek.Load() {
				glog.Debug("pause ignored")
				continue
			}
			glog.Debugf("local: pause state change to %t detected", pause)
			pos, err := p.GetPlaybackPos()
			if err != nil {
				glog.Warnf("failed to get playback position: %v", err)
				break
			}
			glog.Debugf("local: pause state triggered at %s", pos)
			stream.Send(&protocol.RoomEvent{
				Type: protocol.EventType_EVENT_PAUSE,
				Event: &protocol.RoomEvent_PauseEvent{
					PauseEvent: &protocol.PlayPause{
						Pause:     pause,
						Timestamp: durationpb.New(pos),
					},
				},
			})
		case <-p.Quit():
			stream.SendMsg(&protocol.RoomEvent{
				Type: protocol.EventType_EVENT_CLIENT_DISCONNECT,
			})
			stream.CloseSend()
			return nil
		}
	}
}

func recvChanges(cl protocol.UniView_RoomClient, p player.Interface, ignoreSeek *atomic.Bool) {
	for {
		ev, err := cl.Recv()
		if err != nil {
			glog.Errorf("receive failed: %v", err)
			os.Exit(1)
		}

		glog.Debugf("received %s", ev.Type)

		ignoreSeek.Store(true)
		switch ev.Type {
		case protocol.EventType_EVENT_PAUSE:
			pauseEv := ev.GetPauseEvent()
			if pauseEv == nil {
				glog.Warn("received empty pause event")
				continue
			}

			glog.Debugf("remote: pause %t", pauseEv.Pause)
			p.Pause(pauseEv.Pause)
			glog.Debugf("remote: pause jump to %s", pauseEv.Timestamp.AsDuration())
			p.Seek(pauseEv.Timestamp.AsDuration())
		case protocol.EventType_EVENT_JUMP:
			jumpEv := ev.GetJumpEvent()
			if jumpEv == nil {
				glog.Warn("received empty jump event")
				continue
			}

			glog.Debugf("remote: jump to %s", jumpEv.Timestamp.AsDuration())
			p.Seek(jumpEv.Timestamp.AsDuration())
		}
		time.Sleep(50 * time.Millisecond)
		ignoreSeek.Store(false)
	}
}
