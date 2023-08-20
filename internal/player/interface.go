package player

import "time"

type Interface interface {
	Pause(bool) error
	LoadFile(string) error
	Seek(time.Duration) error
	GetPlaybackPos() (time.Duration, error)
	NotifySeek() <-chan time.Duration
	NotifyPause() <-chan bool
	Quit() <-chan struct{}
	Close()
}
