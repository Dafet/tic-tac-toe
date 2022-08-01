package player

import "tic-tac-toe/game/mark"

type Player struct {
	Name      string
	Mark      mark.Mark
	FirstTurn bool
}
