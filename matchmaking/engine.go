package matchmaking

import (
	"errors"
	"math/rand"
	"sync"
	log "tic-tac-toe/logger"
	"tic-tac-toe/ws"

	"time"
)

var (
	ErrAlreadyQueued = errors.New(`player is already queued`)
)

var logger = log.NewDefaultZerolog()

func NewRandom(tickDur time.Duration) *Random {
	return &Random{
		t:      time.NewTicker(tickDur),
		result: nil,
		queue:  []string{},
	}
}

// TODO: fix ineffective queue rewrites on each match?
type Random struct {
	t      *time.Ticker
	result chan ws.PlayerMatch
	queue  []string
	mu     sync.Mutex
}

func (r *Random) QueuePlayer(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if found := r.playerQueued(id); found {
		return ErrAlreadyQueued
	}

	r.queue = append(r.queue, id)

	return nil
}

func (r *Random) UnqueuePlayer(id string) error {
	panic("not implemented") // TODO: Implement
}

func (r *Random) Init(result chan ws.PlayerMatch) {
	r.result = result
	go r.startTick()
}

func (r *Random) startTick() {
	for {
		<-r.t.C // ticker
		// logger.Debug().Msg("ticking random mm engine")

		m, ok := r.match()
		if ok {
			r.result <- m
		}
	}
}

func (r *Random) match() (ws.PlayerMatch, bool) {
	if len(r.queue) < 2 {
		return ws.PlayerMatch{}, false
	}

	p1, p2 := r.pickIDs()

	match := ws.PlayerMatch{
		Player1ID: p1,
		Player2ID: p2,
	}

	return match, true
}

func (r *Random) pickIDs() (string, string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p1Index, p1 := pickRandom(r.queue)
	r.queue = popElem(p1Index, r.queue)

	p2Index, p2 := pickRandom(r.queue)
	r.queue = popElem(p2Index, r.queue)

	return p1, p2
}

func (e *Random) playerQueued(id string) bool {
	// e.mu.Lock()
	// defer e.mu.Unlock()

	// TODO: get rid of range - bit time complexity.
	for _, s := range e.queue {
		if s == id {
			return true
		}
	}

	return false
}

// if player count is "len"%2=0 -> pick random and return theirs id?
// highlight generic errs such as not enough players

func init() {
	rand.NewSource(time.Now().Unix()) // TODO: useless?
}

func popElem[a any](i int, slice []a) []a {
	return append(slice[:i], slice[i+1:]...)
}

// picks random element from "a" and returns element's index
func pickRandom[a any](slice []a) (int, a) {
	l := len(slice)
	n := rand.Intn(l)
	return n, slice[n]
}

// // func popElem(slice []string, s int) []string {
// // 	return append(slice[:s], slice[s+1:]...)
// }
