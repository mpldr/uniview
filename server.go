package main

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"git.sr.ht/~mpldr/uniview/internal/mansion"
	"git.sr.ht/~mpldr/uniview/internal/server"
	"git.sr.ht/~mpldr/uniview/protocol"
	"git.sr.ht/~poldi1405/glog"
	"google.golang.org/grpc"
)

func serverShutdown(signals <-chan os.Signal, srv *grpc.Server, rooms *mansion.Mansion) {
	sig := <-signals
	glog.Infof("received %s, shutting down", sig)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		awaitStop := make(chan struct{})
		go func() {
			srv.GracefulStop()
			close(awaitStop)
		}()
		select {
		case <-awaitStop:
		case <-time.After(5 * time.Second):
			srv.Stop()
		}
	}()

	go func() {
		defer wg.Done()
		rooms.Close()
	}()

	wg.Wait()
}

func startServer() error {
	sigs := make(chan os.Signal, 8)

	srv := grpc.NewServer()
	roomMan := mansion.New()
	protocol.RegisterUniViewServer(srv, &server.Server{
		Rooms: roomMan,
	})
	go serverShutdown(sigs, srv, roomMan)
	signal.Notify(sigs, os.Interrupt)

	glog.Debugf("starting listener on %s:%d", "0.0.0.0", 1558)
	lis, err := net.Listen("tcp", "0.0.0.0:1558")
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	return srv.Serve(lis)
}
