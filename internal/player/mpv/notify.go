package mpv

import "time"

func (m *MPV) NotifyPause() <-chan bool {
	return m.notifyPause
}

func (m *MPV) NotifySeek() <-chan time.Duration {
	return m.notifySeek
}

func (m *MPV) Quit() <-chan struct{} {
	return m.quitchan
}
