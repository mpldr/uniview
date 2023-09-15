package main

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"git.sr.ht/~mpldr/uniview/internal/config"
	"git.sr.ht/~mpldr/uniview/internal/mansion"
	"git.sr.ht/~mpldr/uniview/internal/server"
	wraphttp "git.sr.ht/~mpldr/uniview/internal/server/http"
	"git.sr.ht/~mpldr/uniview/protocol"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

var shutdown []func()

func serverShutdown(signals <-chan os.Signal) {
	sig := <-signals
	slog.Info("signal received. shutting down", "signal", sig)
	var wg sync.WaitGroup

	for _, f := range shutdown {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(f)
	}

	wg.Wait()
}

func startServer() error {
	sigs := make(chan os.Signal, 8)

	roomMan := mansion.New()
	shutdown = append(shutdown, roomMan.Close)

	grpcsrv := grpc.NewServer()
	protocol.RegisterUniViewServer(grpcsrv, &server.Server{
		Rooms: roomMan,
	})
	shutdown = append(shutdown, grpcShutdown(grpcsrv))

	slog.Debug("starting listener", "bind_to", config.Server.General.Bind)
	lis, err := net.Listen("tcp", config.Server.General.Bind)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	var handler http.Handler
	handler, err = wraphttp.NewServer(grpcsrv)
	if err != nil {
		return fmt.Errorf("failed to wrap gRPC: %w", err)
	}

	go serverShutdown(sigs)
	signal.Notify(sigs, os.Interrupt)

	handler = h2c.NewHandler(handler, &http2.Server{})
	srv := &http.Server{
		Addr:    config.Server.General.Bind,
		Handler: handler,
	}

	shutdown = append(shutdown, httpShutdown(srv.Shutdown))

	err = srv.Serve(lis)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func httpShutdown(srv func(context.Context) error) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv(ctx)
	}
}

func grpcShutdown(srv *grpc.Server) func() {
	return func() {
		awaitStop := make(chan struct{})
		go func() {
			srv.GracefulStop()
			close(awaitStop)
		}()
		select {
		case <-awaitStop:
		case <-time.After(5 * time.Second):
			slog.Debug("force-stopping gRPC")
			srv.Stop()
		}
	}
}
