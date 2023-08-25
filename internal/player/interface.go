package player

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import "time"

type Interface interface {
	Pause(bool) error
	LoadFile(string) error
	Seek(time.Duration) error
	GetPlaybackPos() (time.Duration, error)
	GetPauseState() bool
	NotifySeek() <-chan time.Duration
	NotifyPause() <-chan bool
	Quit() <-chan struct{}
	Close()
	PlayerReady() <-chan struct{}
}
