package game

import (
	"errors"
	"tic-tac-toe/game/mark"
	"tic-tac-toe/game/player"
)

// - add more verbose errors?

var (
	ErrCellOccupied    = errors.New(`cell is already occupied`)
	ErrWrongPlayerTurn = errors.New(`invalid player turn: waiting for another player to make a turn`)
	ErrUnknownPlayer   = errors.New(`unknown player`)
	ErrInvalidIndex    = errors.New(`invalid cell index: must be in range 0 - 8`)
)

type Field [9]mark.Mark

func Start(p1Name, p2Name string) *Game {
	players := make(map[string]*player.Player)

	players[p1Name] = &player.Player{Name: p1Name, Mark: Player1Mark, FirstTurn: true}
	players[p2Name] = &player.Player{Name: p2Name, Mark: Player2Mark, FirstTurn: false}

	return &Game{
		players: players,
		grid:    newGrid(),
		tm:      newTurnManager(),
	}
}

type Game struct {
	players map[string]*player.Player
	grid    *grid
	tm      *turnManager
}

type State int

type TurnResult struct {
	IsFinal    bool
	State      State
	WinnerName string
	Err        error
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

	if g.isOver(playerName) {
		return g.makeFinalResult()
	}

	return TurnResult{IsFinal: false, Err: nil}
}

func (g *Game) GetField() Field {
	f := Field{}
	for i, d := range g.grid.data {
		f[i] = d
	}
	return f
}

// make private?
func (f *Field) InitNone() {
	for i := 0; i < 9; i++ {
		f[i] = mark.None
	}
}

func (f *Field) hasNone() bool {
	for _, m := range f {
		if m == mark.None {
			return true
		}
	}
	return false
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

func (g *Game) getPlayerMark(name string) mark.Mark {
	p, found := g.players[name]
	if !found {
		return mark.None
	}

	return p.Mark
}

// return ptr here?
// func (g *Game) getPlayer(name string) player.Player {
// 	p, ok := g.players[name]
// 	if !ok {
// 		return player.Player{}
// 	}

// 	return *p
// }

func (g *Game) isOver(playerName string) bool {
	if !g.isWinPossible() {
		return false
	}

	var (
		in = g.grid.data
		m  = g.getPlayerMark(playerName)
	)

	// make abstraction out of this (strategy pattern or something)
	for _, winIndexes := range winCombos {
		var temp string

		for _, winIndex := range winIndexes {
			temp += string(in[winIndex])
		}

		if len(temp) != 3 {
			continue
		}

		if m == mark.X && temp == winX {
			return true
		}

		if m == mark.O && temp == winO {
			return true
		}
	}

	if !g.grid.hasNone() {
		return true
	}

	return false
}

func (g *Game) isWinPossible() bool {
	return g.tm.currentTurn >= 5
}

// review
func (g *Game) makeFinalResult() TurnResult {
	winner, mark, found := g.getWinnerName()
	if !found {
		return TurnResult{IsFinal: true, State: Draw, WinnerName: "", Err: nil}
	}

	return TurnResult{IsFinal: true, State: getGameState(mark), WinnerName: winner, Err: nil}
}

// review
func (g *Game) getWinnerName() (string, mark.Mark, bool) {
	var in = g.grid.data

	// make abstraction out of this (strategy pattern or something)
	for _, winIndexes := range winCombos {
		var temp string

		for _, winIndex := range winIndexes {
			temp += string(in[winIndex])
		}

		if len(temp) != 3 {
			continue
		}

		switch temp {
		case winX:
			return g.getPlayerNameByMark(mark.X), mark.X, true
		case winO:
			return g.getPlayerNameByMark(mark.O), mark.O, true
		}
	}

	return "", mark.None, false
}

func (g *Game) getPlayerNameByMark(m mark.Mark) string {
	for _, pl := range g.players {
		if m == pl.Mark {
			return pl.Name
		}
	}

	return ""
}

func newGrid() *grid {
	fld := Field{}
	fld.InitNone()
	return &grid{data: fld}
}

type grid struct {
	data Field
}

func (g *grid) placeMark(i int, m mark.Mark) error {
	g.data[i] = m
	return nil
}

func (g *grid) isEligiblePlacement(i int, m mark.Mark) error {
	if i > len(Field{})-1 {
		return ErrInvalidIndex
	}

	if g.data[i] != mark.None {
		return ErrCellOccupied
	}
	return nil
}

func (g *grid) hasNone() bool {
	return g.data.hasNone()
}

func newTurnManager() *turnManager {
	return &turnManager{
		currentTurn: 0,
		nextMark:    Player1Mark,
	}
}

type turnManager struct {
	nextMark    mark.Mark
	currentTurn int
}

func (t *turnManager) makeTurn(m mark.Mark) error {
	t.currentTurn++
	t.switchNextMark()

	return nil
}

func (t *turnManager) validateTurn(m mark.Mark) error {
	if m != t.nextMark {
		return ErrWrongPlayerTurn
	}

	return nil
}

func (t *turnManager) switchNextMark() {
	switch t.nextMark {
	case Player1Mark:
		t.nextMark = Player2Mark
	case Player2Mark:
		t.nextMark = Player1Mark
	}
}

func getGameState(m mark.Mark) State {
	switch m {
	case mark.X:
		return WinX
	case mark.O:
		return WinO
	}

	return 0
}
