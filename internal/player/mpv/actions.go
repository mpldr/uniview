package mpv

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"git.sr.ht/~mpldr/uniview/internal/player"
)

func (p *MPV) Pause(state bool) error {
	if p.dead.Load() {
		return player.ErrPlayerDead
	}
	req := rand.Int()
	slog.Debug("sending pause command", "desired_state", state)
	p.dropPause.Store(true)
	p.send(command{
		Command:   []any{"set_property", "pause", state},
		RequestID: req,
	})
	res := p.getResponse(req)
	slog.Debug("received response", "response", res)
	if res.Error != "" && res.Error != "success" {
		return errors.New(res.Error)
	}
	return nil
}

func (p *MPV) LoadFile(path string) error {
	if p.dead.Load() {
		return player.ErrPlayerDead
	}
	req := rand.Int()
	slog.Debug("sending command to load file", "file", path)
	p.send(command{
		Command:   []any{"loadfile", path},
		RequestID: req,
	})
	res := p.getResponse(req)
	slog.Debug("received response", "response", res)
	if res.Error != "" && res.Error != "success" {
		return errors.New(res.Error)
	}
	return nil
}

func (p *MPV) Seek(ts time.Duration) error {
	if p.dead.Load() {
		return player.ErrPlayerDead
	}
	req := rand.Int()
	slog.Debug("sending command to seek", "seek_to", ts)
	p.dropSeek.Store(true)
	p.send(command{
		Command:   []any{"set_property", "time-pos", float64(ts/time.Millisecond) / 1000},
		RequestID: req,
	})
	res := p.getResponse(req)
	slog.Debug("received response", "response", res)
	if res.Error != "" && res.Error != "success" {
		return errors.New(res.Error)
	}
	return nil
}

func (p *MPV) GetPauseState() bool {
	return p.pauseState
}

func (p *MPV) GetPlaybackPos() (time.Duration, error) {
	if p.dead.Load() {
		return 0, player.ErrPlayerDead
	}
	req := rand.Int()
	slog.Debug("sending command to query playback position")
	p.send(command{
		Command:   []any{"get_property", "time-pos"},
		RequestID: req,
	})
	res := p.getResponse(req)
	slog.Debug("received response", "response", res)
	if pos, ok := res.Data.(float64); ok {
		ts := time.Duration(pos*1000) * time.Millisecond
		return ts, nil
	}
	return 0, fmt.Errorf("get-pos: got '%s' of type %T instead of float64", res.Data, res.Data)
}
