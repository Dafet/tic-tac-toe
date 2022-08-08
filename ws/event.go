package ws

type testEvent struct {
	data string
}

type playerMatchedEvent struct {
	m PlayerMatch
}

type gameFinishedEvent struct {
	gameID     gameID
	isDraw     bool
	winnerID   string
	defeatedID string
}

type waitingTurnEvent struct {
	gameID           gameID
	turnMadeByPlayer string
}

type invalidCellIndexEvent struct {
	cellIndex int
	connID    string
	desc      string
}

type clientDisconnectEvent struct {
	connID string
}
