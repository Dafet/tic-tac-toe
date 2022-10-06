package ws

import (
	"sync"
	"tic-tac-toe/game"

	"github.com/lithammer/shortuuid"
)

// review names
// split interface into 2 - gameManager + gameRepo.
type gameManager interface {
	startGame(p1, p2 string) gameID
	processTurn(id gameID, cellIndex int, playerID string) error
	getGameField(id gameID) (game.Grid, bool)
	getAnotherPlayerID(id gameID, playerID string) (string, bool)
	getGameByID(id gameID) (internalGame, bool)
	getGameByPlayer(playerID string) (internalGame, bool)
	finishGame(id gameID) error
}

type gameID string

func (g gameID) str() string {
	return string(g)
}

func newInMemGameManager(eventCh chan interface{}) *gameManagerInMem {
	return &gameManagerInMem{
		gameList:   make(map[gameID]internalGame),
		playerList: make(map[string]gameID),
		eventCh:    eventCh,
		gmu:        sync.Mutex{},
		pmu:        sync.Mutex{},
	}
}

type internalGame struct {
	p1, p2 string
	id     gameID
	game   *game.Game
}

func (g internalGame) getAnotherPlayerID(player string) string {
	var another string
	switch player {
	case g.p1:
		another = g.p2
	case g.p2:
		another = g.p1
	}
	return another
}

// mark finished games?
type gameManagerInMem struct {
	gameList   map[gameID]internalGame
	playerList map[string]gameID
	eventCh    chan interface{}
	gmu, pmu   sync.Mutex
}

func (m *gameManagerInMem) startGame(p1, p2 string) gameID {
	g := game.Start(p1, p2)

	gameID := gameID(shortuuid.New())

	m.gmu.Lock()
	defer m.gmu.Unlock()

	m.gameList[gameID] = internalGame{
		id:   gameID,
		p1:   p1,
		p2:   p2,
		game: g,
	}

	m.pmu.Lock()
	defer m.pmu.Unlock()

	m.playerList[p1] = gameID
	m.playerList[p2] = gameID

	return gameID
}

func (m *gameManagerInMem) processTurn(id gameID, cellIndex int, playerID string) error {
	g, found := m.gameList[id]
	if !found {
		return errGameNotFound
	}

	result := g.game.MakeTurn(cellIndex, playerID)
	if result.Err != nil {
		return result.Err // make more verbose error
	}

	// review
	if result.IsFinal {
		e := m.makeGameFinishedEvent(result, id, g)
		m.eventCh <- e

		return nil
	}

	m.eventCh <- waitingTurnEvent{
		gameID:           id,
		turnMadeByPlayer: playerID,
	}

	return nil
}

func (m *gameManagerInMem) getGameField(id gameID) (game.Grid, bool) {
	m.gmu.Lock()
	defer m.gmu.Unlock()

	g, ok := m.gameList[id]
	if !ok {
		return game.Grid{}, false
	}

	return g.game.GetGridCopy(), true
}

func (m *gameManagerInMem) getAnotherPlayerID(id gameID, playerID string) (string, bool) {
	m.gmu.Lock()
	defer m.gmu.Unlock()

	g, ok := m.gameList[id]
	if !ok {
		return "", false
	}

	var anotherPlayer string

	switch playerID {
	case g.p1:
		anotherPlayer = g.p2
	case g.p2:
		anotherPlayer = g.p1
	}

	return anotherPlayer, true
}

func (m *gameManagerInMem) getGameByID(id gameID) (internalGame, bool) {
	m.gmu.Lock()
	defer m.gmu.Unlock()

	g, ok := m.gameList[id]
	return g, ok
}

func (m *gameManagerInMem) getGameByPlayer(playerID string) (internalGame, bool) {
	m.pmu.Lock()
	defer m.pmu.Unlock()

	return m.getGameByID(m.playerList[playerID])
}

func (m *gameManagerInMem) finishGame(id gameID) error {
	var p1, p2 string

	g, ok := m.getGameByID(id)
	if ok {
		g.p1 = p1
		g.p2 = p2
	}

	m.gmu.Lock()
	delete(m.gameList, id)
	m.gmu.Unlock()

	m.pmu.Lock()
	delete(m.playerList, p1)
	delete(m.playerList, p2)
	m.pmu.Unlock()

	return nil
}

// review
func (m *gameManagerInMem) makeGameFinishedEvent(r game.TurnResult, id gameID, g internalGame) gameFinishedEvent {
	if r.State == game.Draw {
		return gameFinishedEvent{gameID: id, isDraw: true}
	}

	return gameFinishedEvent{
		gameID:     id,
		isDraw:     false,
		winnerID:   r.WinnerName,
		defeatedID: g.getAnotherPlayerID(r.WinnerName),
	}
}
