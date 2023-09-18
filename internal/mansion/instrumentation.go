package mansion

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

func (m *Mansion) GetRoomCount() int {
	m.roomsMtx.RLock()
	defer m.roomsMtx.RUnlock()
	return len(m.rooms)
}

func (m *Mansion) SessionCounter() uint64 {
	return m.clientID.Load()
}

func (m *Mansion) ActiveSessions() int {
	var count int

	m.roomsMtx.RLock()
	defer m.roomsMtx.RUnlock()
	for _, r := range m.rooms {
		r.clientFeedMtx.Lock()
		for _, c := range r.clientFeed {
			if !c.Dead() {
				count++
			}
		}
		r.clientFeedMtx.Unlock()
	}

	return count
}
