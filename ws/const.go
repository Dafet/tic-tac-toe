package ws

import "errors"

// clear here

const (
	GameStartType       = "game-start"
	MakeTurnType        = "make-turn"
	WaitingTurnType     = "waiting-player-turn"
	ErrCellOccupiedType = "cell-occupied" // redundant?
	GameFinishedType    = "game-finished"
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

	// server types.
	connType          = "connection"
	playerRdyType     = "play-ready"
	disconnectType    = "disconnect"
	testType          = "test"
	setPlayerDataType = "set-player-data"
)

var (
	errUnsupportedMsgType = errors.New(`unsupported msg type`)
	errEmptyMsgType       = errors.New(`empty msg type`)
	errGameNotFound       = errors.New(`game is not found`)
)
