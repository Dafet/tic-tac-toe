package game

import (
	"errors"
)

// - X always goes first
// - add more verbose errors?

var (
	ErrCellOccupied    = errors.New(`cell is already occupied`)
	ErrWrongPlayerTurn = errors.New(`invalid player turn: waiting for another player to make a turn`)
	ErrUnknownPlayer   = errors.New(`unknown player`)
	ErrInvalidIndex    = errors.New(`invalid cell index: must be in range 0 - 8`)
)

type Player struct {
	Name      string
	Mark      Mark
	FirstTurn bool
}

type Mark string

const (
	X Mark = "x"
	O Mark = "o"

	None Mark = "-"
)

type Field [9]Mark

func Start(p1Name, p2Name string) *Game {
	players := make(map[string]*Player)

	players[p1Name] = &Player{Name: p1Name, Mark: X, FirstTurn: true}
	players[p2Name] = &Player{Name: p2Name, Mark: O, FirstTurn: false}

	return &Game{
		players: players,
		grid:    newGrid(),
		tm:      newTurnManager(),
	}
}

type Game struct {
	players map[string]*Player
	grid    *grid
	tm      *turnManager
}

type TurnResult struct {
	IsFinal bool
	Err     error
	Winner  Player
}

func (g *Game) MakeTurn(cellIndex int, playerName string) TurnResult {
	var err error

	if err = g.validateTurn(cellIndex, playerName); err != nil {
		return TurnResult{Err: err}
	}

	var mark = g.getPlayerMark(playerName)
	if err = g.tm.makeTurn(mark); err != nil {
		return TurnResult{Err: err}
	}

	if err = g.grid.placeMark(cellIndex, mark); err != nil {
		return TurnResult{Err: err}
	}

	if g.isWin(playerName) {
		return g.makeWinResult(playerName)
	}

	return TurnResult{IsFinal: false, Err: nil}
}

func (g *Game) GetGrid() Field {
	f := Field{}
	for i, d := range g.grid.data {
		f[i] = d
	}
	return f
}

func (g *Game) validateTurn(cellIndex int, playerName string) error {
	var err error

	if !g.playerExists(playerName) {
		return ErrUnknownPlayer
	}

	var mark = g.getPlayerMark(playerName)
	if err = g.tm.validateTurn(mark); err != nil {
		return err
	}

	err = g.grid.isEligiblePlacement(cellIndex, mark)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) playerExists(name string) bool {
	_, exists := g.players[name]
	return exists
}

func (g *Game) getPlayerMark(name string) Mark {
	p, found := g.players[name]
	if !found {
		return None
	}

	return p.Mark
}

// return ptr here?
func (g *Game) getPlayer(name string) Player {
	p, ok := g.players[name]
	if !ok {
		return Player{}
	}

	return *p
}

func (g *Game) isWin(playerName string) bool {
	if !g.isWinPossible() {
		return false
	}

	var (
		in   = g.grid.data
		mark = g.getPlayerMark(playerName)
	)

	for _, winIndexes := range winCombos {
		var temp string

		for _, winIndex := range winIndexes {
			temp += string(in[winIndex])
		}

		if len(temp) != 3 {
			continue
		}

		if mark == X && temp == winX {
			return true
		}

		if mark == O && temp == winO {
			return true
		}
	}

	return false
}

func (g *Game) isWinPossible() bool {
	return g.tm.currentTurn >= 5
}

func (g *Game) makeWinResult(playerName string) TurnResult {
	return TurnResult{
		IsFinal: true,
		Err:     nil,
		Winner:  g.getPlayer(playerName),
	}
}

func newGrid() *grid {
	fld := Field{}
	fld.initNone()
	return &grid{data: fld}
}

type grid struct {
	data Field
}

func (g *grid) placeMark(i int, m Mark) error {
	g.data[i] = m
	return nil
}

func (g *grid) isEligiblePlacement(i int, m Mark) error {
	if i > len(Field{})-1 {
		return ErrInvalidIndex
	}

	if g.data[i] != None {
		return ErrCellOccupied
	}
	return nil
}

const (
	player1Mark = X
	player2Mark = O
)

func newTurnManager() *turnManager {
	return &turnManager{
		currentTurn: 0,
		nextMark:    player1Mark,
	}
}

type turnManager struct {
	nextMark    Mark
	currentTurn int
}

func (t *turnManager) makeTurn(m Mark) error {
	t.currentTurn++
	t.switchNextMark()

	return nil
}

func (t *turnManager) validateTurn(m Mark) error {
	if m != t.nextMark {
		return ErrWrongPlayerTurn
	}

	return nil
}

func (t *turnManager) switchNextMark() {
	switch t.nextMark {
	case player1Mark:
		t.nextMark = player2Mark
	case player2Mark:
		t.nextMark = player1Mark
	}
}

func (f *Field) initNone() {
	for i := 0; i < 9; i++ {
		f[i] = None
	}
}
