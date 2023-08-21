package mpv

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"git.sr.ht/~mpldr/uniview/internal/player"
	"git.sr.ht/~poldi1405/glog"
)

func (p *MPV) Pause(state bool) error {
	if p.dead.Load() {
		return player.ErrPlayerDead
	}
	req := rand.Int()
	glog.Tracef("sending command to set pause to %t", state)
	p.send(command{
		Command:   []any{"set_property", "pause", state},
		RequestID: req,
	})
	res := p.getResponse(req)
	glog.Tracef("received response %#v", res)
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
	glog.Tracef("sending command to load %q", path)
	p.send(command{
		Command:   []any{"loadfile", path},
		RequestID: req,
	})
	res := p.getResponse(req)
	glog.Tracef("received response %#v", res)
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
	glog.Tracef("sending command to seek to %s", ts)
	p.send(command{
		Command:   []any{"set_property", "time-pos", float64(ts/time.Millisecond) / 1000},
		RequestID: req,
	})
	res := p.getResponse(req)
	glog.Tracef("received response %#v", res)
	if res.Error != "" && res.Error != "success" {
		return errors.New(res.Error)
	}
	return nil
}

func (p *MPV) GetPlaybackPos() (time.Duration, error) {
	if p.dead.Load() {
		return 0, player.ErrPlayerDead
	}
	req := rand.Int()
	glog.Trace("sending command to query playback position")
	p.send(command{
		Command:   []any{"get_property", "time-pos"},
		RequestID: req,
	})
	res := p.getResponse(req)
	glog.Tracef("received response %#v", res)
	if pos, ok := res.Data.(float64); ok {
		ts := time.Duration(pos*1000) * time.Millisecond
		return ts, nil
	}
	return 0, fmt.Errorf("get-pos: got '%s' of type %T instead of float64", res.Data, res.Data)
}
