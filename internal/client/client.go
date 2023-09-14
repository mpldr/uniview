package client

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"git.sr.ht/~mpldr/uniview/internal/client/api"
	"git.sr.ht/~mpldr/uniview/internal/player"
	"git.sr.ht/~mpldr/uniview/internal/player/mpv"
	"git.sr.ht/~mpldr/uniview/protocol"
	"git.sr.ht/~poldi1405/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

func StartClient(u *url.URL) error {
	var p player.Interface
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	waitForExit := make(chan struct{})

	sigs := make(chan os.Signal, 8)
	go serverShutdown(ctx, cancel, sigs, waitForExit)
	signal.Notify(sigs, os.Interrupt)

	glog.Debug("starting player…")
	p, err = mpv.New()
	if err != nil {
		return fmt.Errorf("failed to start mpv: %w", err)
	}
	defer p.Close()

	status := api.StatusConnectionConnecting
	go StartRestServer(context.Background(), p, &status, u)

	glog.Debugf("loading file %q…", u.Query().Get("file"))
	err = p.LoadFile(u.Query().Get("file"))
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	if u.Port() == "" {
		u.Host += ":443"
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	if u.Query().Has("insecure") || insecureByDefault(u.Hostname()) {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	glog.Debugf("connecting to remote…")
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	gconn, err := grpc.DialContext(dialCtx, u.Host, opts...)
	dialCancel()
	if err != nil {
		return fmt.Errorf("failed to connect to server %q: %w", u.Host, err)
	}
	defer gconn.Close()

	status = api.StatusConnectionOk

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

	room := strings.TrimPrefix(u.Path, "/")
	if len(room) == 0 {
		return errors.New("no room provided")
	}

	glog.Debugf("Joining room %q", room)
	err = stream.Send(&protocol.RoomEvent{
		Type: protocol.EventType_EVENT_JOIN,
		Event: &protocol.RoomEvent_Join{
			Join: &protocol.RoomJoin{
				Name:      room,
				Timestamp: durationpb.New(pos),
				Url:       u.Query().Get("file"),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	glog.Debug("waiting for remote events…")
	go func() {
		receiveEvents(ctx, p, stream)
		select {
		case <-ctx.Done():
		default:
			cancel()
		}
	}()
	go func() {
		sendPlayerEvents(ctx, p, stream)
		select {
		case <-ctx.Done():
		default:
			cancel()
		}
	}()

	glog.Debug("waiting for shutdown")
	<-waitForExit
	glog.Debug("completed shutdown")

	return nil
}

func sendPlayerEvents(ctx context.Context, p player.Interface, cl protocol.UniView_RoomClient) {
	glog.Debug("waiting for player events…")

loop:
	for {
		select {
		case timestamp := <-p.NotifySeek():
			glog.Debugf("local: seek to %s detected", timestamp)
			cl.Send(&protocol.RoomEvent{
				Type: protocol.EventType_EVENT_JUMP,
				Event: &protocol.RoomEvent_JumpEvent{
					JumpEvent: &protocol.PlaybackJump{
						Timestamp: durationpb.New(timestamp),
					},
				},
			})
		case pause := <-p.NotifyPause():
			glog.Debugf("local: pause state change to %t detected", pause)
			pos, err := p.GetPlaybackPos()
			if err != nil {
				glog.Warnf("failed to get playback position: %v", err)
				break
			}
			glog.Debugf("local: pause state triggered at %s", pos)
			cl.Send(&protocol.RoomEvent{
				Type: protocol.EventType_EVENT_PAUSE,
				Event: &protocol.RoomEvent_PauseEvent{
					PauseEvent: &protocol.PlayPause{
						Pause:     pause,
						Timestamp: durationpb.New(pos),
					},
				},
			})
		case <-ctx.Done():
			break loop
		case <-p.Quit():
			break loop
		}
	}
	cl.SendMsg(&protocol.RoomEvent{
		Type: protocol.EventType_EVENT_CLIENT_DISCONNECT,
	})
	cl.CloseSend()
}

func receiveEvents(ctx context.Context, p player.Interface, cl protocol.UniView_RoomClient) {
	events := make(chan *protocol.RoomEvent, 16)
	go func() {
		defer close(events)
		for {
			ev, err := cl.Recv()
			if err != nil {
				glog.Errorf("receive failed: %v", err)
				return
			}
			events <- ev
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case ev, open := <-events:
			if !open {
				return
			}
			glog.Debugf("received %s", ev.Type)

			<-p.PlayerReady()
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
			case protocol.EventType_EVENT_SERVER_CLOSE:
				glog.Info("received shutdown notification from the server. Disconnecting.")
			case protocol.EventType_EVENT_SERVER_PING: // ignore
			default:
				glog.Warnf("received unknown event: %s", ev.Type)
			}
		}
	}
}

func insecureByDefault(host string) bool {
	dnsServerList := []string{
		"9.9.9.9",        // Quad9
		"45.11.45.11",    // dns.sb
		"1.1.1.1",        // CloudFlare
		"8.8.8.8",        // Google
		"208.67.222.222", // OpenDNS
	}

	for _, dnsServer := range dnsServerList {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "udp", dnsServer)
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 128*time.Millisecond)
		ip, err := r.LookupIP(ctx, "ip", host)
		cancel()
		if err != nil {
			continue
		}
		if len(ip) == 0 {
			continue
		}
		return ip[0].IsLoopback() || ip[0].IsPrivate()
	}

	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return false
	}
	return ip.IP.IsLoopback()
}
