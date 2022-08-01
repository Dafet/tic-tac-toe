package ws

import (
	"errors"
	"fmt"
	"tic-tac-toe/matchmaking"
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
		return s.makeUserDataCmd(connID, msg)
	case playerRdyKind:
		logger.Debug().Msgf("player id: %s is rdy", connID)
	}

	return nil, fmt.Errorf(`unknown msg kind: %s`, msg.Kind)
}

func (s *serverCmdFactory) makeUserDataCmd(connID string, msg *Msg) (setUserDataCmd, error) {
	p, found := s.s.handler.getPlayerDataPtr(connID)
	if !found {
		return setUserDataCmd{}, ErrPlayerNotFound
	}

	data, ok := msg.Data.(SetUserDataMsg)
	if !ok {
		return setUserDataCmd{}, errors.New(`msg is not SetUserDataMsg`)
	}

	return setUserDataCmd{
		newName:        data.NewName,
		existingPlayer: p,
	}, nil
}

func (s *serverCmdFactory) makePlayerRdyCmd(connID string, msg Msg) (setUserDataCmd, error) {
	panic("implement me")
}

type setUserDataCmd struct {
	newName        string
	existingPlayer *playerData
}

func (c setUserDataCmd) apply() error {
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
	mm     *matchmaking.Engine
	connID string
	player *playerData
}

func (c playerRdyCmd) apply() error {
	return c.mm.QueuePlayer(c.connID, c.player.compilePlayer())
}
