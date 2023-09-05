package client

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"os"
	"sync"
	"time"

	"git.sr.ht/~poldi1405/glog"
)

var shutdownFuncs []func(context.Context) error

func serverShutdown(ctx context.Context, cancel context.CancelFunc, signals <-chan os.Signal, closeOnDone chan<- struct{}) {
	defer close(closeOnDone)

	select {
	case sig := <-signals:
		cancel()
		glog.Infof("received %s, shutting down", sig)
	case <-ctx.Done():
	}

	waitCtx, waitCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer waitCancel()

	var wg sync.WaitGroup
	for i, f := range shutdownFuncs {
		wg.Add(1)
		go func(ctx context.Context, id int, f func(context.Context) error) {
			defer wg.Done()

			glog.Debugf("shutting down %d", id)
			awaitStop := make(chan struct{})
			go func() {
				defer close(awaitStop)

				err := f(ctx)
				if err != nil {
					glog.Warnf("failed to cleanly shutdown: %v", err)
				}
			}()
			select {
			case <-awaitStop:
			case <-ctx.Done():
				return
			}
			glog.Debugf("completed shutdown of %d", id)
		}(waitCtx, i, f)
	}

	wg.Wait()
	glog.Debug("unlocking quit")
}
