package client

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
)

var shutdownFuncs []func(context.Context) error

func serverShutdown(ctx context.Context, cancel context.CancelFunc, signals <-chan os.Signal, closeOnDone chan<- struct{}) {
	defer close(closeOnDone)

	select {
	case sig := <-signals:
		cancel()
		slog.Info("received signal. shutting down…", "signal", sig)
	case <-ctx.Done():
	}

	waitCtx, waitCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer waitCancel()

	var wg sync.WaitGroup
	for i, f := range shutdownFuncs {
		wg.Add(1)
		go func(ctx context.Context, id int, f func(context.Context) error) {
			defer wg.Done()

			slog.Debug("shutting down", "id", id)
			awaitStop := make(chan struct{})
			go func() {
				defer close(awaitStop)

				err := f(ctx)
				if err != nil {
					slog.Warn("failed to shutdown gracefully", "error", err)
				}
			}()
			select {
			case <-awaitStop:
			case <-ctx.Done():
				return
			}
			slog.Debug("completed shutdown", "id", id)
		}(waitCtx, i, f)
	}

	wg.Wait()
	slog.Debug("unlocking quit")
}
