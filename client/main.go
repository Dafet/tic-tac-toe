package main

import (
	"fmt"
	"os"
	"os/signal"
	log "tic-tac-toe/logger"
	"tic-tac-toe/ws"
)

const (
	defaultAddr = "localhost:8080"
	playerName  = "Sample"
)

var logger = log.NewDefaultZerolog()

func main() {
	conn, err := ws.Connect(defaultAddr)
	if err != nil {
		logger.Fatal().Err(err).Msg("error establishing connection")
	}
	defer conn.Close()

	// todo:
	// - make proper graceful shutdown?

	// finish := make(chan struct{})

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		<-interrupt

		logger.Info().Msg("cought interrupt signal, finishing")

		if err = conn.Close(); err != nil {
			logger.Fatal().Err(err).Msg("error closing connection")
		}

		os.Exit(1)
	}()

	var playerName = mustGetPlayerName()

	var msgChan = conn.ListenForServer()

	go sendSetUserDataMsg(conn, playerName)

	// go func() {
	// 	for {
	// 		logger.Info().Msg("sending test msg")
	// 		conn.SendMsg(ws.NewTestMsg("Josh"))
	// 		time.Sleep(time.Second * 1)
	// 	}
	// }()

	// rewrite - there is a bottleneck for msg receiving.
	for {
		m, moreMsg := <-msgChan

		if !moreMsg {
			logger.Info().Msg("no more messages from connection, finish listening")
			break
		}

		logger.Info().Msgf("received from server: %v", m)

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

func mustGetPlayerName() string {
	var name string

	fmt.Println("enter player's name: ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		logger.Fatal().Err(err).Msg("error reading player's name")
	}

	return name
}

func sendSetUserDataMsg(conn *ws.Connection, name string) {
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
