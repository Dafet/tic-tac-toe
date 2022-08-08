package ws

import (
	"errors"
	"fmt"
	"tic-tac-toe/game"
)

type cmd interface {
	apply() error
}

type cmdFactory interface {
	make(connID string, msg *Msg) (cmd, error)
}

func newServerCmdFactory(s *Server) cmdFactory {
	return &serverCmdFactory{s: s}
}

type serverCmdFactory struct {
	s *Server
}

func (s *serverCmdFactory) make(connID string, msg *Msg) (cmd, error) {
	if msg == nil {
		return nil, errors.New(`msg is nil`)
	}

	switch msg.Kind {
	case setUserDataKind:
		return s.makeSetPlayerDataCmd(connID, msg)
	case playerRdyKind:
		return s.makePlayerRdyCmd(connID)
		// do nothing here
	// case GameStartKind:
	case MakeTurnKind:
		return s.makeMakeTurnCmd(connID, msg, s.s.eventCh)
	}

	return nil, fmt.Errorf(`unknown msg kind: %s`, msg.Kind)
}

func (s *serverCmdFactory) makeSetPlayerDataCmd(connID string, msg *Msg) (setPlayerDataCmd, error) {
	p, found := s.s.handler.getPlayerDataPtr(connID)
	if !found {
		return setPlayerDataCmd{}, ErrPlayerNotFound
	}

	data, ok := msg.Data.(SetUserDataMsg)
	if !ok {
		return setPlayerDataCmd{}, errors.New(`msg is not SetUserDataMsg`)
	}

	return setPlayerDataCmd{
		newName:        data.NewName,
		existingPlayer: p,
	}, nil
}

func (s *serverCmdFactory) makeMakeTurnCmd(connID string, msg *Msg, events chan interface{}) (makeTurnCmd, error) {
	data, ok := msg.Data.(MakeTurnMsg)
	if !ok {
		return makeTurnCmd{}, errors.New(`msg is not SetUserDataMsg`)
	}

	c := makeTurnCmd{
		cellIndex: data.CellIndex,
		connID:    connID,
		gameID:    gameID(data.GameID),
		gm:        s.s.gm,
		eventCh:   events,
	}

	return c, nil
}

func (s *serverCmdFactory) makePlayerRdyCmd(connID string) (playerRdyCmd, error) {
	return playerRdyCmd{
		connID: connID,
		engine: s.s.mmengine,
	}, nil
}

type setPlayerDataCmd struct {
	newName        string
	existingPlayer *playerData
}

func (c setPlayerDataCmd) apply() error {
	if c.newName == "" {
		return nil // error here?
	}

	if c.existingPlayer == nil {
		return errors.New(`player is nil`)
	}

	c.existingPlayer.name = c.newName

	return nil
}

type playerRdyCmd struct {
	connID string
	engine MatchmakingEngine
}

func (c playerRdyCmd) apply() error {
	logger.Info().Str("conn_id", c.connID).Msg("queueing player")
	return c.engine.QueuePlayer(c.connID)
}

type makeTurnCmd struct {
	cellIndex int
	connID    string
	gameID    gameID
	gm        gameManager
	eventCh   chan interface{}
}

func (c makeTurnCmd) apply() error {
	// review error handling
	err := c.gm.processTurn(c.gameID, c.cellIndex, c.connID)
	if err != nil {

		// review - make proper error processing - some util func?
		if err == game.ErrCellOccupied {
			c.eventCh <- invalidCellIndexEvent{
				cellIndex: c.cellIndex,
				connID:    c.connID,
				desc:      fmt.Sprintf("cell [%v] is already occupied", c.cellIndex),
			}
		}

		return fmt.Errorf(`error processing turn: %w`, err)
	}

	return nil
}
