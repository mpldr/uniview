package mpv

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"git.sr.ht/~mpldr/uniview/internal/player"
	"github.com/adrg/xdg"
	"github.com/fsnotify/fsnotify"
)

type MPV struct {
	conn net.Conn
	cmd  *exec.Cmd

	responses         map[int]response
	responsesMtx      sync.RWMutex
	responseBroadcast *sync.Cond

	notifyPause chan bool
	notifySeek  chan time.Duration
	quitchan    chan struct{}
	playerReady chan struct{}

	notifySeekInternal chan struct{}
	notifyIdle         chan struct{}
	commands           chan command
	pauseState         bool

	dead      *atomic.Bool
	dropSeek  *atomic.Bool
	dropPause *atomic.Bool
}

func New() (*MPV, error) {
	runtimeDir := xdg.RuntimeDir
	socketpath, err := os.MkdirTemp(runtimeDir, "uniview_")
	if err != nil {
		return nil, fmt.Errorf("failed to create socket dir: %w", err)
	}
	socketPath := filepath.Join(socketpath, "uniview.sock")

	// find mpv binary
	mpvPath, err := exec.LookPath("mpv")
	if err != nil {
		return nil, fmt.Errorf("could not find executable 'mpv': %w", err)
	}

	// start filesystem watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for socket: %w", err)
	}

	err = watcher.Add(filepath.Dir(socketPath))
	if err != nil {
		return nil, fmt.Errorf("failed to watch: %w", err)
	}

	// start mpv
	p := &MPV{
		responses:          make(map[int]response),
		responseBroadcast:  sync.NewCond(&sync.Mutex{}),
		notifyPause:        make(chan bool, 1),
		notifySeek:         make(chan time.Duration, 1),
		notifySeekInternal: make(chan struct{}, 1),
		notifyIdle:         make(chan struct{}, 1),
		quitchan:           make(chan struct{}, 1),
		playerReady:        make(chan struct{}),
		commands:           make(chan command, 16),
		dead:               &atomic.Bool{},
		dropSeek:           &atomic.Bool{},
		dropPause:          &atomic.Bool{},
	}
	p.cmd = exec.Command(mpvPath, "--input-ipc-server="+socketPath, "--player-operation-mode=pseudo-gui", "--idle")
	if err := p.cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting mpv: %w", err)
	}
	go func(p *MPV) {
		p.cmd.Wait()
		p.dead.Store(true)
	}(p)

	// wait for socket
outer:
	for {
		select {
		case ev := <-watcher.Events:
			slog.Debug("filewatcher received event", "event", ev)
			if ev.Has(fsnotify.Create) {
				watcher.Close()
				break outer
			}
		case <-time.After(1 * time.Second):
			p.cmd.Process.Kill()
			return nil, player.ErrPlayerNotStarted
		}
	}

	// connect to socket
	p.conn, err = net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("error connecting to socket: %w", err)
	}

	go p.monitor()

	<-p.notifyIdle
	slog.Debug("confirmed player idle state")

	go p.pollPause()
	go p.handleSeekEvents()

	return p, nil
}

func (p *MPV) Name() string {
	return "mpv"
}

func (p *MPV) Close() {
	if p.dead.Load() {
		return
	}
	p.cmd.Process.Kill()
}
