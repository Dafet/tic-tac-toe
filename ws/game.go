package ws

import (
	"tic-tac-toe/game"

	"github.com/lithammer/shortuuid"
)

// review names

type gameManager interface {
	startGame(p1, p2 string) gameID
	processTurn(id gameID, cellIndex int, playerID string) error
}

type gameID string

func (g gameID) str() string {
	return string(g)
}

func newGameManagerImpl() *gameManagerImpl {
	return &gameManagerImpl{
		list: make(map[gameID]internalGame),
	}
}

type gameManagerImpl struct {
	list map[gameID]internalGame
}

type internalGame struct {
	p1   string
	p2   string
	game *game.Game
}

func (m *gameManagerImpl) startGame(p1, p2 string) gameID {
	g := game.Start(p1, p2)

	id := gameID(shortuuid.New())

	// lock here?
	m.list[id] = internalGame{
		p1:   p1,
		p2:   p2,
		game: g,
	}

	return id
}

// incapsulate errors here?
func (m *gameManagerImpl) processTurn(id gameID, cellIndex int, playerID string) error {
	g, found := m.list[id]
	if !found {
		return ErrGameNotFound
	}

	// validate turn itself

	// make turn

	// check turn errors

	// process result.IsFinal bool

	result := g.game.MakeTurn(cellIndex, playerID)
	if result.Err != nil {
		return result.Err // make verbose error
	}

	return nil
}
