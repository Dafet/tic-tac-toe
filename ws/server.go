package ws

import (
	"fmt"
	log "tic-tac-toe/logger"

	"github.com/gorilla/websocket"
)

// this should be a facade for mm, websocket
// TODO:
//  - clear old games (due to disconnects etc)
//	- second player always goes first?
//	- send msg to clients if their game has broken (p2 disconnect etc)?
//  - add sending client msgs logs
//  - add retry msg to clients
//  - привести к единому виду терминалогию connID/playerID
//  - выделить единый ошибочный тип для Msg + cmd
//  - unqueue player on disconnect
//  - recover from panic

var (
	logger     = log.NewDefaultZerolog()
	msgfactory msgFactory
	cmdfactory cmdFactory
	// mmengine   MatchmakingEngine // initialize with empty mm engine?
)

func New() *Server {
	eventCh := make(chan interface{})

	s := &Server{
		eventCh: eventCh,
		handler: newWsHandler(eventCh),
	}

	msgfactory = newMsgFactory(s)
	cmdfactory = newServerCmdFactory(s)

	return s
}

type Server struct {
	handler  *wsHandler
	gm       gameManager
	mmengine MatchmakingEngine
	eventCh  chan interface{} // change type to concrete/abstract one
}

func (s *Server) SetMatchmakingEngine(mm MatchmakingEngine) {
	s.mmengine = mm
}

func (s *Server) Start() {
	go s.startEventLoop()

	s.initGameManager()
	s.initMatchmaking()

	s.handler.startHandle()
}

func (s *Server) initMatchmaking() {
	if s.mmengine == nil {
		logger.Fatal().Msg("matchmaking engine is not initialized")
	}

	matchRdy := make(chan PlayerMatch)

	// review
	go func() {
		for {
			m, ok := <-matchRdy
			if !ok {
				logger.Error().Msg("matchRdy chan is closed")
				return
			}

			s.eventCh <- playerMatchedEvent{m: m}
		}
	}()

	s.mmengine.Init(matchRdy)
}

func (s *Server) initGameManager() {
	s.gm = newGameManagerInMem(s.eventCh)
}

// process errors!
func (s *Server) sendGameStartMsgs(m PlayerMatch, id gameID) error {
	var (
		m1  = newGameStartMsg(true, id.str())
		m2  = newGameStartMsg(false, id.str())
		err error
	)

	if err = s.handler.sendMsg(m.Player1ID, m1); err != nil {
		err = fmt.Errorf(`player1 error: %w`, err)
		return err
	}

	if err = s.handler.sendMsg(m.Player2ID, m2); err != nil {
		err = fmt.Errorf(`player2 error: %w`, err)
		return err
	}

	return nil
}

// test func
func (s *Server) BroadcastMsg(msg *Msg) {
	for _, v := range s.handler.conns {
		b, err := serializeMsg(msg)
		if err != nil {
			logger.Fatal().Err(err).Msg("error serializing broadcast")
		}

		v.wsconn.c.WriteMessage(websocket.TextMessage, b)
	}
}

// test func
func (s *Server) LogConns() {
	for _, v := range s.handler.conns {
		logger.Debug().Msg("currently connected: " + v.name)
	}
}

func (s *Server) startEventLoop() {
	for {
		// start async func here for every event push?
		i := <-s.eventCh
		switch e := i.(type) {
		case testEvent:
			logger.Debug().Msg("got test event: " + e.data)
		case playerMatchedEvent: // pack into separate func?
			logger.Info().
				Str("player_1", e.m.Player1ID).
				Str("player_2", e.m.Player2ID).
				Msg("match is found, starting game")

			id := s.gm.startGame(e.m.Player1ID, e.m.Player2ID)

			if err := s.sendGameStartMsgs(e.m, id); err != nil {
				// try something here?
				logger.Error().Err(err).Msg("error sending start game msg")
			}
		case invalidCellIndexEvent:
			s.processInvalidCellIndexEvent(e)
		case waitingTurnEvent:
			s.processWaitingForTurnEvent(e)
		case gameFinishedEvent:
			s.processGameFinishedEvent(e)
		case clientDisconnectEvent:
			s.processClientDisconnectEvent(e)
		default:
			logger.Error().Msgf("unknown event type: %+v", e)
		}
	}
}

func (s *Server) processWaitingForTurnEvent(e waitingTurnEvent) {
	// send e.playerID waiting for turn msg

	fld, found := s.gm.getGameField(e.gameID)
	if !found {
		logger.Warn().Msgf("cannot find gameID: [%s] for waitingNextTurn msg", e.gameID.str())
		// how to process it?
		return
	}

	playerID, found := s.gm.getAnotherPlayerID(e.gameID, e.turnMadeByPlayer)
	if !found {
		logger.Warn().Msgf("cannot get another playerID for waitingNextTurn msg, gameID: [%s], previous turn was made by: [%s]", e.gameID.str(), e.turnMadeByPlayer)
		// how to process it?
		return
	}

	m := newWaitingTurnMsg(e.gameID, fld)
	err := s.handler.sendMsg(playerID, m)
	if err != nil {
		// how to process correctly - retry logic?
		logger.Error().Err(err).Msgf("error sending '%s' msg", WaitingTurnType)
	}
}

func (s *Server) processInvalidCellIndexEvent(e invalidCellIndexEvent) {
	m := newErrorMsg(ErrCellOccupiedType, e.desc)
	err := s.handler.sendMsg(e.connID, m)
	if err != nil {
		// how to process correctly - retry logic?
		logger.Error().Err(err).Msgf("error sending '%s' msg", WaitingTurnType)
	}
}

func (s *Server) processGameFinishedEvent(e gameFinishedEvent) {
	fld, found := s.gm.getGameField(e.gameID)
	if !found {
		logger.Warn().Msgf("cannot find gameID: [%s] for waitingNextTurn msg", e.gameID.str())
		// how to process it?
		return
	}

	// make func for this?
	if e.isDraw {
		logger.Info().Msgf("game [%s] has finished with [draw]", e.gameID)
		g, ok := s.gm.getGameByID(e.gameID)
		if !ok {
			logger.Error().Msgf("game [%s] is not found", e.gameID)
			return
		}

		var drawMsg = newGameFinishedDrawMsg(e.gameID, fld)

		if err := s.handler.sendMsg(g.p1, drawMsg); err != nil {
			logger.Error().Err(err).Msg("error sending draw msg to player1")
		}

		if err := s.handler.sendMsg(g.p2, drawMsg); err != nil {
			logger.Error().Err(err).Msg("error sending draw msg to player2")
		}

		return
	}

	// make func for this?
	var (
		winMsg    = newGameFinishedWinMsg(e.gameID, fld)
		defeatMsg = newGameFinishedDefeatMsg(e.gameID, fld)
		err       error
	)

	logger.Info().Msgf("game [%s] has finished, winner [%s]", e.gameID, e.winnerID)

	if err = s.handler.sendMsg(e.winnerID, winMsg); err != nil {
		logger.Error().Err(err).Msg("error sending win msg")
	}

	if err = s.handler.sendMsg(e.defeatedID, defeatMsg); err != nil {
		logger.Error().Err(err).Msg("error sending defeat msg")
	}

	if err = s.gm.finishGame(e.gameID); err != nil {
		logger.Error().Err(err).Msg("error finishing game")
	}
}

func (s *Server) processClientDisconnectEvent(e clientDisconnectEvent) {
	if e.connID == "" {
		return
	}

	logger := logger.With().Str("conn_id", e.connID).Logger()

	var err error

	logger.Info().Msg("unqueueing player")
	if err = s.mmengine.UnqueuePlayer(e.connID); err != nil {
		logger.Warn().Err(err).Str("conn_id", e.connID).Msg("error unqueueing player")
		// return
	}

	g, found := s.gm.getGameByPlayer(e.connID)
	if found {
		anotherPlayer := g.getAnotherPlayerID(e.connID)
		dcMsg := newGameFinishedDisconnectMsg(g.id, g.game.GetGridCopy())

		if err = s.handler.sendMsg(anotherPlayer, dcMsg); err != nil {
			logger.Error().Err(err).Msg("error sending finish msg")
		}
	} else {
		logger.Info().Msg("player currenly is not playing any game")
	}

	s.gm.finishGame(g.id)
}
