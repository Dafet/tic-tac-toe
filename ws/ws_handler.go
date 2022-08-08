// websocket server for tic-tac-toe game
package ws

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/lithammer/shortuuid"
	"github.com/rs/zerolog"
)

// TODO: implement ping, pong
// add locks for reading (conns)?

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // todo: perform some actual checks
	},
}

func newWsHandler(eventCh chan interface{}) *wsHandler {
	return &wsHandler{
		eventCh: eventCh,
		conns:   make(map[string]*playerData),
	}
}

type wsHandler struct {
	// conns map[string]*websocket.Conn
	conns      map[string]*playerData
	eventCh    chan interface{}
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
		logger := logger.With().Str("conn_id", connID).Logger()

		logger.Info().Msg("listening for conn")
		h.readClientMsgs(connID, logger)
		logger.Info().Msg("stop listening for conn")

		h.flushConn(connID)
	}()
}

func (h *wsHandler) storeConn(c *websocket.Conn, playerName string) (key string) {
	connKey := shortuuid.New()
	h.conns[connKey] = &playerData{
		wsconn: Connection{c: c},
		name:   playerName,
	}
	return connKey
}

func (h *wsHandler) readClientMsgs(connID string, logger zerolog.Logger) error {
	conn, ok := h.conns[connID]
	if !ok {
		return errors.New(`connection is not established`)
	}

	var readErr error

outer:
	for {
		mt, msgRaw, err := conn.wsconn.ReadMessage()

		logger.Info().
			Str("msg_raw", string(msgRaw)).
			Int("websocket_msg_type", mt).
			Msg("[incoming] client's msg")

		switch {
		case IsCloseError(err):
			h.eventCh <- clientDisconnectEvent{connID: connID}
			return nil
		case err != nil:
			logger.Error().Err(err).Msg("error reading client msg, finish listening")
			readErr = err
			break outer
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

		logger.Info().Str("cmd_name", reflect.TypeOf(cmd).Name()).Msg("applying command")

		if err = cmd.apply(); err != nil {
			logger.Error().Err(err).Msg("error applying cmd")
			continue
		}
	}

	if readErr != nil {
		return readErr
	}

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

func (h *wsHandler) flushConn(key string) error {
	conn, ok := h.conns[key]
	if !ok {
		return nil
	}

	if err := conn.wsconn.Close(); err != nil {
		return fmt.Errorf(`error closing conn: %w`, err)
	}

	delete(h.conns, key)

	return nil
}

func (h *wsHandler) sendMsg(connID string, msg *Msg) error {
	c, ok := h.conns[connID]
	if !ok {
		return errors.New(`connID is not found`)
	}

	msgRaw, err := jsoniter.Marshal(msg)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling msg for logging")
	}

	logger.Info().
		Str("conn_id", connID).
		Str("msg_raw", string(msgRaw)).
		Int("websocket_msg_type", websocket.TextMessage). // hardcoded
		Msg("[outgoing] sending msg to client")

	return c.wsconn.SendMsg(msg)
}

// move up?
type playerData struct {
	// wsconn *websocket.Conn
	wsconn Connection
	name   string
}
