// This is a sample client.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"tic-tac-toe/game"
	"tic-tac-toe/game/mark"
	log "tic-tac-toe/logger"
	"tic-tac-toe/ws"
	"time"
)

const (
	defaultAddr = "localhost:8080"
)

var (
	logger = log.NewDefaultZerolog()

	playerName = ""

	conn       *ws.Connection
	gameFld    game.Field
	gameID     string
	playerMark mark.Mark
)

func init() {
	linuxClear()
}

func main() {
	var err error
	if conn, err = ws.Connect(defaultAddr); err != nil {
		logger.Fatal().Err(err).Msg("error establishing connection")
	}
	defer conn.Close()

	gameFld.InitNone()

	go initInterruptSignal()

	playerName = mustGetPlayerName()

	go sendSetPlayerDataMsg(conn, playerName)
	time.Sleep(time.Millisecond * 30)
	go sendRdyMsg(conn, playerName)

	fmt.Println("looking for a game...")

	var msgChan = listenForServer()

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
			playerMark = msg.Mark

			fmt.Println("starting game...")
			fmt.Println("game id: ", gameID)
			fmt.Println("mark   : ", playerMark)

			drawGrid(gameFld)

			if msg.FirstTurn {
				fmt.Println("you are making turn, write cell index (1-9)")
				processTurn(playerMark)
				redrawGameData()
			}

			fmt.Println("waiting for another player to make a turn...")
		case ws.WaitingTurnKind:
			msg, err := ws.MakeWaitingTurnMsg(m)
			if err != nil {
				logger.Fatal().Err(err).Msg("error asserting game msg")
			}

			gameFld = msg.Field
			redrawGameData()

			fmt.Println("you are making turn, write cell index (1-9)")

			processTurn(playerMark)

			redrawGameData()
			fmt.Println("waiting for another player to make a turn...")
		case ws.ErrCellOccupiedKind:
			fmt.Println("cell is already occupied, choose another")

			processTurn(playerMark)

			redrawGameData()
		case ws.GameFinishedKind:
			msg, err := ws.MakeGameFinishedMsg(m)
			if err != nil {
				logger.Fatal().Err(err).Msg("error asserting game msg")
			}

			gameFld = msg.Field
			redrawGameData()

			switch {
			case msg.PlayerWon:
				fmt.Println("game is over, you won!")
			case msg.IsDraw:
				fmt.Println("game is over, draw")
			case msg.OpponentDisconect:
				fmt.Println("game is over, your opponent has disconnected")
			default:
				fmt.Println("game is over, you have been defeated")
			}

			processRetryGame()
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

func listenForServer() <-chan *ws.Msg {
	// todo:
	// - make proper graceful shutdown?
	// - add logs?

	// process interrupt?
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	msgChan := make(chan *ws.Msg)

	go func() {
		defer close(msgChan)
		for {
			_, raw, err := conn.ReadMessage()

			if err != nil {
				// todo: what to do in case of errors
				logger.Warn().Err(err).Msg("error reading msg from server")
				break
			}

			msg, err := ws.DeserializeMsg(raw)
			if err != nil {
				logger.Error().Err(err).Msg("error reading server msg")
				continue
			}

			msgChan <- msg
		}
	}()

	return msgChan
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

func processTurn(m mark.Mark) error {
	var cellIndex int

	for {
		var inputRaw string
		_, err := fmt.Scan(&inputRaw)
		if err != nil {
			fmt.Println("error reading player's input: ", err)
			continue
		}

		input, err := strconv.Atoi(inputRaw)
		if err != nil {
			fmt.Println("invalid input, must be number from 1 to 9, write again ", err)
			continue
		}

		cellIndex = input - 1

		if cellIndex < 0 || cellIndex > 8 {
			fmt.Println("invalid number, must be from 1 to 9")
			continue
		}

		break
	}

	if placedMark := gameFld[cellIndex]; placedMark != mark.None {
		fmt.Printf("cell is already occupied with [%s], place another mark \n", m.Str())
		return processTurn(m)
	}

	gameFld[cellIndex] = m

	msg := ws.NewMakeTurnMsg(cellIndex, gameID)
	err := conn.SendMsg(msg)
	if err != nil {
		logger.Fatal().Err(err).Msg("error sending makeTurn msg")
	}

	return nil
}

func processRetryGame() {
	fmt.Println("start a new game? (y/n)")

	var input string
	for {
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Println("error reading player's input: ", err)
			continue
		}

		if strings.ToLower(input) != "y" && strings.ToLower(input) != "n" {
			continue
		}

		break
	}

	if input == "y" {
		gameFld.InitNone()
		sendRdyMsg(conn, playerName)
		linuxClear()

		fmt.Println("looking for a game...")
		return
	}

	if err := conn.Close(); err != nil {
		logger.Fatal().Err(err).Msg("error closing connection")
	}

	os.Exit(0)
}

func redrawGameData() {
	linuxClear()

	fmt.Println("game id: ", gameID)
	fmt.Println("mark: ", playerMark)
	drawGrid(gameFld)
}

func linuxClear() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
