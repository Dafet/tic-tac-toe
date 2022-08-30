package game

import (
	"errors"
	"tic-tac-toe/game/mark"
)

var (
	ErrCellOccupied    = errors.New(`cell is already occupied`)
	ErrWrongPlayerTurn = errors.New(`invalid player turn: waiting for another player to make a turn`)
	ErrUnknownPlayer   = errors.New(`unknown player`)
	ErrInvalidIndex    = errors.New(`invalid cell index: must be in range 0 - 8`)
)

const (
	Player1Mark = mark.X
	Player2Mark = mark.O
)

var (
	// to win a game certain indexes should be matched with below ones
	winComboRow1      = []int{0, 1, 2}
	winComboRow2      = []int{3, 4, 5}
	winComboRow3      = []int{6, 7, 8}
	winComboCol1      = []int{0, 3, 6}
	winComboCol2      = []int{1, 4, 7}
	winComboCol3      = []int{2, 5, 8}
	winComboDiagonal1 = []int{0, 4, 8}
	winComboDiagonal2 = []int{2, 4, 6}

	winCombos = [][]int{winComboRow1, winComboRow2, winComboRow3, winComboCol1, winComboCol2, winComboCol3, winComboDiagonal1, winComboDiagonal2}

	// emptyTurnResult = TurnResult{}
)

const (
	winX = "xxx"
	winO = "ooo"
)
