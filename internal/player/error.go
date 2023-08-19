package player

import "errors"

var (
	ErrPlayerDead       = errors.New("player has quit")
	ErrPlayerNotStarted = errors.New("player did not start")
)
