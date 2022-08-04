package ws

import (
	"fmt"
	log "tic-tac-toe/logger"

	"github.com/gorilla/websocket"
)

// this should be a facade for mm, websocket
// TODO: clear old games (due to disconnects etc)

var (
	logger     = log.NewDefaultZerolog()
	msgfactory msgFactory
	cmdfactory cmdFactory
	// mmengine   MatchmakingEngine // initialize with empty mm engine?
)

func New() *Server {
	s := &Server{
		handler: newWsHandler(),
	}

	msgfactory = newMsgFactory(s)
	cmdfactory = newServerCmdFactory(s)

	return s
}

type Server struct {
	handler  *wsHandler
	gm       gameManager
	mmengine MatchmakingEngine
	eventCh  chan interface{}
}

func (s *Server) SetMatchmakingEngine(mm MatchmakingEngine) {
	s.mmengine = mm
}

func (s *Server) Start() {
	s.eventCh = make(chan interface{})
	go s.startEventLoop()

	s.initGM()
	s.initMM()

	s.handler.startHandle()
}

func (s *Server) initMM() {
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

func (s *Server) initGM() {
	s.gm = newGameManagerImpl()
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
		case gameFinishedEvent:
		default:
			logger.Error().Msgf("unknown event type: %+v", e)
		}
	}
}

// func (s *Server) listenForPlayerMatch() {
// 	for {
// 		m, ok := <-s.matchRdy
// 		if !ok {
// 			// review
// 			logger.Error().Msg("matchRdy chan is closed")
// 			return
// 		}

// 		logger.Debug().Msgf("match is rdy: %+v", m)

// 		id := s.gm.startGame(m.Player1ID, m.Player2ID)

// 		if err := s.sendGameStartMsgs(m, id); err != nil {
// 			// try something here?
// 			logger.Error().Err(err).Msg("error sending start game msg")
// 		}
// 	}
// }
