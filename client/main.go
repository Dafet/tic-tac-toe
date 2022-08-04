package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"tic-tac-toe/game"
	"tic-tac-toe/game/mark"
	log "tic-tac-toe/logger"
	"tic-tac-toe/ws"
	"time"
)

const (
	defaultAddr = "localhost:8080"
	playerName  = "Sample"
)

var (
	logger = log.NewDefaultZerolog()

	conn    *ws.Connection
	gameFld game.Field
	gameID  string
)

func main() {
	var err error
	if conn, err = ws.Connect(defaultAddr); err != nil {
		logger.Fatal().Err(err).Msg("error establishing connection")
	}
	defer conn.Close()

	gameFld.InitNone()

	go initInterruptSignal()

	// todo:
	// - make proper graceful shutdown?

	// finish := make(chan struct{})

	var (
		playerName = mustGetPlayerName()
		msgChan    = conn.ListenForServer()
	)

	// implement msg send queue?

	go sendSetPlayerDataMsg(conn, playerName)
	time.Sleep(time.Millisecond * 30)
	go sendRdyMsg(conn, playerName)

	// go func() {
	// 	for {
	// 		logger.Info().Msg("sending test msg")
	// 		conn.SendMsg(ws.NewTestMsg("Josh"))
	// 		time.Sleep(time.Second * 1)
	// 	}
	// }()

	// rewrite - there is a bottleneck for msg receiving.
	// review - huge tabs
	for {
		m, moreMsg := <-msgChan

		if !moreMsg {
			logger.Info().Msg("no more messages from connection, finish listening")
			break
		}

		switch m.Kind {
		case ws.GameStartKind:
			msg, err := ws.MakeGameStartMsg(m)
			if err != nil {
				logger.Fatal().Err(err).Msg("error asserting game msg")
			}

			gameID = msg.GameID

			logger.Info().Msgf("starting game with id: %s, first turn: %v, mark: %v",
				gameID,
				msg.FirstTurn,
				msg.Mark.Str())

			// hardcoded init - get fld state from msg?

			drawGrid(gameFld)

			if msg.FirstTurn {
				// make a turn as well as send msg
				processFirstTurn(msg.Mark)
			}

		case "next turn":

		default:
			logger.Info().Msgf("received from server: %+v", m)
		}

		// time.Sleep(time.Millisecond * 500)

		// - check msg type
		// 	- if (start game) : draw board

		// if msg type is start game -> initiate game start?
	}

	fmt.Println("thanks for playing!")

	// connect to server
	// send rdy signal
	// establish msg receive loop
}

func initInterruptSignal() {
	var err error

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt

	logger.Info().Msg("cought interrupt signal, finishing")

	if err = conn.Close(); err != nil {
		logger.Fatal().Err(err).Msg("error closing connection")
	}

	os.Exit(1)
}

func mustGetPlayerName() string {
	var name string

	fmt.Println("enter player's name: ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		logger.Fatal().Err(err).Msg("error reading player's name")
	}

	return name
}

func sendSetPlayerDataMsg(conn *ws.Connection, name string) {
	err := conn.SendMsg(ws.NewSetUserDataMsg(name))
	if err != nil {
		logger.Fatal().Err(err).Msg("error starting game")
	}
}

func sendRdyMsg(conn *ws.Connection, name string) {
	err := conn.SendMsg(ws.NewPlayerRdyMsg(name))
	if err != nil {
		logger.Fatal().Err(err).Msg("error starting game")
	}
}

func drawGrid(fld game.Field) {
	b := strings.Builder{}
	defer b.Reset()

	for i, mark := range fld {
		if (i+3)%3 == 0 {
			b.WriteString("\n")
		}

		b.WriteString(string(mark))
	}

	fmt.Println(b.String())
}

func processFirstTurn(m mark.Mark) error {
	var cell int

	logger.Info().Msg("you are making first turn, write cell index (1-9)")

	_, err := fmt.Scanln(&cell)
	if err != nil {
		logger.Fatal().Err(err).Msg("error getting player's input")
	}

	i := cell - 1

	if i < 0 {
		logger.Fatal().Msg("invalid cell index: must not be less than 1") // "less than 1 is confusing - rewrite"
	}

	gameFld[i] = m

	err = conn.SendMsg(ws.NewMakeTurnMsg(i, gameID))
	if err != nil {
		logger.Fatal().Err(err).Msg("error sending makeTurn msg")
	}

	return nil
}
