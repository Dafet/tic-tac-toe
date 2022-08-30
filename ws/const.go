package ws

import "errors"

// clear here

const (
	GameStartKind       = "game-start"
	MakeTurnKind        = "make-turn"
	WaitingTurnKind     = "waiting-player-turn"
	ErrCellOccupiedKind = "cell-occupied"
	GameFinishedKind    = "game-finished"
)

var (
	ErrPlayerNotFound       = errors.New(`player is not found`)
	ErrPlayerAlreadyInQueue = errors.New(`player is already in a game queue`)
	ErrPlayerNotInQueue     = errors.New(`player is not in a game queue`)
)

// TODO: make configurable needed values
const (
	upgradeConnHandlerPath = "/conn"
	port                   = ":8080"
	playerNameParam        = "player_name"

	// server kinds.
	connKind          = "connection"
	playerRdyKind     = "play-ready"
	disconnectKind    = "disconnect"
	testKind          = "test"
	setPlayerDataKind = "set-player-data"
)

var (
	errUnsupportedMsgKind = errors.New(`unsupported msg kind`)
	errEmptyMsgKind       = errors.New(`empty msg kind`)
	errGameNotFound       = errors.New(`game is not found`)
)
