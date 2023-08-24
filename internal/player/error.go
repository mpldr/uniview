package player

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import "errors"

var (
	ErrPlayerDead       = errors.New("player has quit")
	ErrPlayerNotStarted = errors.New("player did not start")
)
