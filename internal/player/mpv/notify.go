package mpv

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

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

func (m *MPV) PlayerReady() <-chan struct{} {
	return m.playerReady
}
