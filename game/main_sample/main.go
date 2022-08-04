package main

import (
	"fmt"
	"os"
	"strings"
	"tic-tac-toe/game"
	"tic-tac-toe/game/player"
	l "tic-tac-toe/logger"
)

var logger = l.NewDefaultZerolog()

const (
	p1Name = "Steeve"
	p2Name = "Jake"
)

func main() {
	consoleSample()
	// debugVersion()
}

func consoleSample() {
	g := game.Start(p1Name, p2Name)

	logger.Info().Msg("starting game")

	logger.Info().Msg("player1 is " + p1Name)
	logger.Info().Msg("player2 is " + p2Name)

	finish := make(chan player.Player)

	go func() {
		p := <-finish
		logger.Info().Msgf("game is finished, winner: %+v", p.Name)
		os.Exit(0)
	}()

	for {
		drawGrid(g.GetField())

		logger.Info().Msg("waiting for player1 cell index: ")

		var p1Turn int
		_, err := fmt.Scanln(&p1Turn)
		if err != nil {
			logger.Fatal().Err(err).Msg("error reading p1 turn")
		}

		var result game.TurnResult
		result = g.MakeTurn(p1Turn, p1Name)

		if result.Err != nil {
			logger.Fatal().Err(result.Err).Msg("error making turn as p1")
		}

		if result.IsFinal {
			finish <- result.Winner
		}

		drawGrid(g.GetField())

		logger.Info().Msg("waiting for player2 cell index: ")

		var p2Turn int
		_, err = fmt.Scanln(&p2Turn)
		if err != nil {
			logger.Fatal().Err(result.Err).Msg("error reading p1 turn")
		}

		result = g.MakeTurn(p2Turn, p2Name)

		if result.Err != nil {
			logger.Fatal().Err(result.Err).Msg("error making turn as p2")
		}

		drawGrid(g.GetField())

		if result.IsFinal {
			finish <- result.Winner
		}
	}
}

func debugVersion() {
	g := game.Start(p1Name, p2Name)

	logger.Info().Msg("starting game")

	logger.Info().Msg("player1 is " + p1Name)
	logger.Info().Msg("player2 is " + p2Name)

	for {
		drawGrid(g.GetField())

		logger.Info().Msg("waiting for player1 cell index: ")

		var p1Turn = 3

		var result game.TurnResult
		result = g.MakeTurn(p1Turn, p1Name)

		if result.Err != nil {
			logger.Fatal().Err(result.Err).Msg("error making turn as p1")
		}

		if result.IsFinal {
			logger.Info().Msgf("game is finished, winner: %+v", result.Winner)
		}

		drawGrid(g.GetField())

		logger.Info().Msg("waiting for player2 cell index: ")

		var p2Turn = 4

		result = g.MakeTurn(p2Turn, p2Name)

		if result.Err != nil {
			logger.Fatal().Err(result.Err).Msg("error making turn as p2")
		}

		drawGrid(g.GetField())

		if result.IsFinal {
			logger.Info().Msgf("game is finished, winner: %+v", result.Winner)
		}
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

	logger.Info().Msg(b.String())
}

func someFunc() {

}
