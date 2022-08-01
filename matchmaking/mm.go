package matchmaking

import (
	"errors"
	"sync"
	"tic-tac-toe/game/player"
)

const (
	playerCount = 2
)

var (
	ErrAlreadyQueued = errors.New(`player is already queued`)
	ErrNoMatcher     = errors.New(`matcher is not set`)
)

// func NewDebugEngine() Engine {

// }

// in-memory version - change to abstraction here?
type Engine struct {
	// rdy player list
	playRdy map[string]player.Player
	mu      sync.Mutex
	// matcher matcher
}

func (e *Engine) QueuePlayer(id string, p player.Player) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.playerQueued(id) {
		return ErrAlreadyQueued
	}

	e.playRdy[id] = p

	return nil
}

func (e *Engine) UnQueuePlayer(id string, p player.Player) {
	panic("implement me")
}

func (e *Engine) GetMatch() *Match {
	// we should block here
	panic("finish")
}

// move?
// transfer?
type Match struct {
	Player1ID string
	Player2ID string
}

func (e *Engine) TryMatchPlayers() error {
	// if e.matcher == nil {
	// 	return ErrNoMatcher
	// }

	// return e.matcher.match()
	panic("finish")
}

func (e *Engine) playerQueued(id string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, found := e.playRdy[id]
	return found
}
