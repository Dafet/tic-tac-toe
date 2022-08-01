// websocket server for tic-tac-toe game
package ws

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"tic-tac-toe/game/mark"
	"tic-tac-toe/game/player"

	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
)

// todo: implement ping, pong, close msgs

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // todo: perform some actual checks
	},
}

func newWsHandler() *wsHandler {
	return &wsHandler{
		conns: make(map[string]*playerData),
	}
}

// don't reference whole server in msg structs?
type wsHandler struct {
	// conns map[string]*websocket.Conn
	conns      map[string]*playerData
	playerLock sync.Mutex
}

// todo: add separate async func?
func (h *wsHandler) startHandle() {
	http.HandleFunc(upgradeConnHandlerPath, h.upgradeConn)
	logger.Info().Msgf("websoket server is listening on %s 🔥", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal().Err(err).Msg("error listening server")
	}
}

// todo: rename + rewrite (too many responsibilities and operations)
func (h *wsHandler) upgradeConn(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errMsg := "error upgrading connection: " + err.Error()
		fmt.Println(errMsg)
		fmt.Fprintln(w, errMsg)
	}

	// defer connection.Close() // move

	go func() {
		connID := h.storeConn(c, "") // review playerName param.
		logger.Debug().Msg("stored connection under: " + connID)

		h.readClientMsg(connID)

		logger.Info().Msg("stop listening for: " + connID)

		h.flushConn(connID) // delete from rdyConns map as well?
	}()

	// w.WriteHeader(200) // ok? no - this causing panic
}

func (h *wsHandler) storeConn(c *websocket.Conn, playerName string) (key string) {
	connKey := shortuuid.New()
	h.conns[connKey] = &playerData{
		wsconn: c,
		name:   playerName,
	}
	return connKey
}

func (h *wsHandler) flushConn(key string) error {
	conn, ok := h.conns[key]
	if !ok {
		return nil
	}

	if err := conn.wsconn.Close(); err != nil {
		return err
	}

	delete(h.conns, key)

	return nil
}

// there are 1:1 ratio conn to func call (issue?)
func (h *wsHandler) readClientMsg(connID string) error {
	conn, ok := h.conns[connID]
	if !ok {
		return errors.New(`connection is not established`)
	}

	var readErr error

	for {
		mt, msgRaw, err := conn.wsconn.ReadMessage()

		if mt == websocket.CloseMessage {
			logger.Debug().Msg("client sent close message, flushing connection")

			if err = h.flushConn(connID); err != nil {
				logger.Error().Err(err).Msg("error flushing connection")
			}
		}

		if err != nil {
			logger.Error().Err(err).Msg("error reading msg from connection, stopping to listen")
			readErr = err
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана

			// fix: "websocket: close 1006 (abnormal closure): unexpected EOF"
		}

		msg, err := msgfactory.make(msgRaw)
		if err != nil {
			logger.Error().Err(err).Msg("error compiling msg")
			continue
		}

		cmd, err := cmdfactory.make(connID, msg)
		if err != nil {
			logger.Error().Err(err).Msg("error creating cmd")
			continue
		}

		if err = cmd.apply(); err != nil {
			// add cmd name into log?
			logger.Error().Err(err).Str("cmd_name", reflect.TypeOf(cmd).Name()).Msg("error applying cmd")
		}
	}

	if readErr != nil {
		return readErr
	}

	// panic("why are we here?") // change panic/text
	return nil
}

func (h *wsHandler) getPlayerDataPtr(connID string) (*playerData, bool) {
	if connID == "" {
		return &playerData{}, false
	}

	h.playerLock.Lock()
	defer h.playerLock.Unlock()

	p, found := h.conns[connID]
	if !found {
		return &playerData{}, false
	}

	return p, true
}

func compilePlayerName(r *http.Request) string {
	name := r.URL.Query().Get(playerNameParam)
	if name != "" {
		return name
	}

	return shortuuid.New()
}

// move up?
type playerData struct {
	wsconn *websocket.Conn
	name   string
}

func (p *playerData) compilePlayer() player.Player {
	return player.Player{
		Name:      p.name,
		Mark:      mark.X, // currently hardcoded - fix - insert into playerData struct more data?
		FirstTurn: false,
	}
}
