package client

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"git.sr.ht/~mpldr/uniview/internal/client/api"
	"git.sr.ht/~mpldr/uniview/internal/player"
	"git.sr.ht/~mpldr/uniview/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

func StartClient(u *url.URL) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	waitForExit := make(chan struct{})

	sigs := make(chan os.Signal, 8)
	go serverShutdown(ctx, cancel, sigs, waitForExit)
	signal.Notify(sigs, os.Interrupt)

	slog.Debug("starting player…")
	p, err := getPlayer()
	if err != nil {
		return fmt.Errorf("failed to start player: %w", err)
	}
	defer p.Close()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	if u.Query().Has("insecure") || insecureByDefault(u.Hostname()) {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if insecureByDefault(u.Hostname()) {
			q := u.Query()
			q.Add("insecure", "")
			u.RawQuery = q.Encode()
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	status := api.StatusConnectionConnecting
	go StartRestServer(context.Background(), p, &status, u)

	slog.Debug("loading file", "file", u.Query().Get("file"))
	err = p.LoadFile(u.Query().Get("file"))
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	if u.Port() == "" {
		u.Host += ":443"
	}

	slog.Debug("connecting to remote", "host", u.Host, "insecure_explicit", u.Query().Has("insecure"), "insecure_ip", insecureByDefault(u.Hostname()))
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	gconn, err := grpc.DialContext(dialCtx, u.Host, opts...)
	dialCancel()
	if err != nil {
		return fmt.Errorf("failed to connect to server %q: %w", u.Host, err)
	}
	defer gconn.Close()

	status = api.StatusConnectionOk

	slog.Debug("requesting handle")
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

	slog.Debug("joining room", "room", room)
	err = stream.Send(&protocol.RoomEvent{
		Type: protocol.EventType_EVENT_JOIN,
		Event: &protocol.RoomEvent_Join{
			Join: &protocol.RoomJoin{
				Name:      room,
				Timestamp: durationpb.New(pos),
				Url:       u.Query().Get("file"),
				Password:  u.Fragment,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	slog.Debug("waiting for remote events")
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

	slog.Debug("waiting for shutdown")
	<-waitForExit
	slog.Debug("completed shutdown")

	return nil
}

func sendPlayerEvents(ctx context.Context, p player.Interface, cl protocol.UniView_RoomClient) {
	slog.Debug("waiting for player events")

loop:
	for {
		select {
		case timestamp := <-p.NotifySeek():
			slog.Debug("seek detected", "origin", "local", "seek_to", timestamp)
			cl.Send(&protocol.RoomEvent{
				Type: protocol.EventType_EVENT_JUMP,
				Event: &protocol.RoomEvent_JumpEvent{
					JumpEvent: &protocol.PlaybackJump{
						Timestamp: durationpb.New(timestamp),
					},
				},
			})
		case pause := <-p.NotifyPause():
			slog.Debug("pause state change detected", "origin", "local", "state", pause)
			pos, err := p.GetPlaybackPos()
			if err != nil {
				slog.Warn("failed to get playback position", "error", err)
				break
			}
			slog.Debug("sending pause state update", "state", pos, "timestamp", pos)
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
				slog.Error("receive failed", "error", err)
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
			slog.Debug("received event", "type", ev.Type)

			<-p.PlayerReady()
			switch ev.Type {
			case protocol.EventType_EVENT_PAUSE:
				pauseEv := ev.GetPauseEvent()
				if pauseEv == nil {
					slog.Warn("received empty pause event")
					continue
				}

				slog.Debug("pause state change detected", "origin", "remote", "state", pauseEv.Pause)
				p.Pause(pauseEv.Pause)
				slog.Debug("seek detected", "origin", "remote", "seek_to", pauseEv.Timestamp.AsDuration())
				p.Seek(pauseEv.Timestamp.AsDuration())
			case protocol.EventType_EVENT_JUMP:
				jumpEv := ev.GetJumpEvent()
				if jumpEv == nil {
					slog.Warn("received empty jump event")
					continue
				}

				slog.Debug("seek detected", "origin", "remote", "seek_to", jumpEv.Timestamp.AsDuration())
				p.Seek(jumpEv.Timestamp.AsDuration())
			case protocol.EventType_EVENT_SERVER_CLOSE:
				slog.Info("received shutdown notification from the server. Disconnecting.")
			case protocol.EventType_EVENT_SERVER_PING: // ignore
			default:
				slog.Warn("received unknown event", "type", ev.Type)
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

	ipaddr, err := net.ResolveIPAddr("ip", host)
	slog.Debug("resolving host using system resolver", "ip", ipaddr.IP, "error", err)
	ip := ipaddr.IP
	if err != nil {
		ip = net.ParseIP(host)
		slog.Debug("parsing as literal IP", "host", host, "ip", ip)
	}
	return ip.IsLoopback() || ip.IsPrivate()
}
