package ws

import (
	log "tic-tac-toe/logger"
	"tic-tac-toe/matchmaking"

	"github.com/gorilla/websocket"
)

// this should be a facade for mm, websocket

var (
	logger     = log.NewDefaultZerolog()
	msgfactory msgFactory
	cmdfactory cmdFactory
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
	handler *wsHandler
	// gm       gameManager
	mmengine matchmaking.Engine
}

func (s *Server) Start() {
	s.handler.startHandle()
}

// test func
func (s *Server) BroadcastMsg(msg string) {
	for _, v := range s.handler.conns {
		v.wsconn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

// test func
func (s *Server) LogConns() {
	for _, v := range s.handler.conns {
		logger.Debug().Msg("currently connected: " + v.name)
	}
}

func (s *Server) startGame(m matchmaking.Match) {
	go func() {
		// game := game.Start(m.Player1ID, m.Player2ID)

		// send msg to client Player1ID (game start, grid, turn)
		// send msg to client Player2ID (game start, grid, turn)

		// start loop for receiving data from conn
	}()
}

// func newGameManager() *gameManager {
// 	fld := game.Field{}
// 	fld.InitNone()
// 	return &gameManager{
// 		field: fld,
// 		game:  &game.Game{},
// 	}
// }

// type gameManager struct {
// 	field game.Field
// 	game  *game.Game
// }

// func (g *gameManager) startGame(m matchmaking.Match) {
// 	g.game = game.Start(m.Player1ID, m.Player2ID)
// }
