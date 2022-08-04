package ws

// todo:
//  - add some player repo interface?

type PlayerMatch struct {
	Player1ID string
	Player2ID string
}

// rename + enrich with player params?
type MatchmakingEngine interface {
	QueuePlayer(id string) error
	UnqueuePlayer(id string) error
	Init(result chan PlayerMatch) // review
}
