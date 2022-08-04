package ws

import "errors"

// clear here

var (
	ErrPlayerNotFound = errors.New(`player is not found`)
)

const (
	upgradeConnHandlerPath = "/conn" // /conn?
	port                   = ":8080" // make configurable
	playerNameParam        = "player_name"
)

var (
	ErrUnsupportedMsgKind = errors.New(`unsupported msg kind`)
	ErrEmptyMsgKind       = errors.New(`empty msg kind`)
	ErrGameNotFound       = errors.New(`game is not found`)
)

const (
	// server kind?
	connKind       = "connection"
	playerRdyKind  = "play-ready"
	disconnectKind = "disconnect"
	testKind       = "test"

	setUserDataKind = "set-user-data"

	// client side?
	GameStartKind = "game-start"
	MakeTurnKind  = "make-turn"
	// player-turn-msg etc.
)
