package ws

type testEvent struct {
	data string
}

type playerMatchedEvent struct {
	m PlayerMatch
}

type gameFinishedEvent struct {
	gameID string
}
