package game_test

import (
	"testing"
	"tic-tac-toe/game"
)

func TestStartGame(t *testing.T) {
	var table = map[string]struct {
		p1, p2 string
	}{
		"test": {
			p1: "123",
			p2: "123",
		},
	}

	for name, v := range table {
		t.Run(name, func(t *testing.T) {
			_ = game.Start(v.p1, v.p2)
			t.FailNow()
		})
	}

}
