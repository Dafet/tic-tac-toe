package matchmaking

import (
	"math/rand"
	"tic-tac-toe/game/player"

	// "tic-tac-toe/player"
	"time"
)

type randomMatcher struct {
	playRdy map[string]player.Player
}

func (m randomMatcher) matchAsync() (Match, error) {
	if len(m.playRdy) < playerCount {
		// return Match{}, ErrNoEnoughPlayers
		panic("finish me")
	}

	return m.pickRandom()
}

func (m randomMatcher) pickRandom() (Match, error) {
	var keys = make([]string, 0, len(m.playRdy))
	for k := range m.playRdy {
		keys = append(keys, k)
	}

	p1i := rand.Intn(len(keys))
	p1 := keys[p1i]

	keys = popElemGen(keys, p1i)

	p2i := rand.Intn(len(keys))
	p2 := keys[p2i]

	return Match{
		Player1ID: p1,
		Player2ID: p2,
	}, nil
}

// if player count is "len"%2=0 -> pick random and return theirs id?
//highlight generic errs such as not enough players

func init() {
	rand.NewSource(time.Now().Unix())
}

func popElemGen[a any](slice []a, s int) []a {
	return append(slice[:s], slice[s+1:]...)
}

func popElem(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
