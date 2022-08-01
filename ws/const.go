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
)

const (
	connKind        = "connection"
	playerRdyKind   = "play-ready"
	disconnectKind  = "disconnect"
	testKind        = "test"
	makeTurnKind    = "make-turn"
	setUserDataKind = "set-user-data"
	// player-turn-msg etc.
)
