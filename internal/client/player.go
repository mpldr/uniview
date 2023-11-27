package client

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"errors"
	"fmt"

	"git.sr.ht/~mpldr/uniview/internal/config"
	"git.sr.ht/~mpldr/uniview/internal/player"
	"git.sr.ht/~mpldr/uniview/internal/player/mpv"
)

func getPlayer() (player.Interface, error) {
	userPreference := config.Client.General.Player
	if userPreference != "" {
		if factory, ok := playerBuilder[userPreference]; ok {
			return factory()
		}

		return nil, fmt.Errorf("unknown player %q", userPreference)
	}
	for _, player := range playerPriority {
		p, err := playerBuilder[player]()
		if err != nil {
			continue
		}
		return p, nil
	}
	return nil, errors.New("no player found")
}

var playerPriority = []string{
	"mpv",
}

var playerBuilder = map[string]func() (player.Interface, error){
	"mpv": mpv.New,
}
