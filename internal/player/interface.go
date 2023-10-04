// Package player provides the universal interface for remote-controlling
// a player. It is assumed that the constructor takes care of launching the
// player properly, without requiring manual user intervention.
package player

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import "time"

// Interface is used to programmatically control the player
type Interface interface {
	// Pause is used to set the pause state of the player
	Pause(bool) error
	// LoadFile is used to load a file into the player
	LoadFile(string) error
	// Seek is used to seek a position in the player
	Seek(time.Duration) error
	// GetPlaybackPos is used to get the current playback position
	GetPlaybackPos() (time.Duration, error)
	// GetPauseState is used to get the current pause state
	GetPauseState() bool
	// NotifySeek sends the current playback position when a seek
	// event is reported by the player
	NotifySeek() <-chan time.Duration
	// NotifyPause sends the current pause state when a pause state
	// change is reported by the player
	NotifyPause() <-chan bool
	// Quit reports that the underlying player has terminated
	Quit() <-chan struct{}
	// Close is used to terminate the underlying player
	Close()
	// PlayerReady reports that the underlying player is ready to
	// receive commands
	PlayerReady() <-chan struct{}
	// Name reports the name of the underlying player
	Name() string
}
