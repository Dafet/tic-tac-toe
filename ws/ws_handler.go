// websocket server for tic-tac-toe game
package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/lithammer/shortuuid"
	"github.com/rs/zerolog"
)

const (
	pingInterval = time.Second * 10
	pingDeadline = pingInterval
)

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

type playerData struct {
	wsconn Connection
	name   string
}

type wsHandler struct {
	conns      map[string]*playerData
	eventCh    chan interface{}
	playerLock sync.Mutex
}

func (h *wsHandler) startHandle() {
	http.HandleFunc(upgradeConnHandlerPath, h.upgradeConn)
	logger.Info().Msgf("websoket server is listening on %s 🔥", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal().Err(err).Msg("error listening server")
	}
}

func (h *wsHandler) upgradeConn(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errMsg := "error upgrading connection: " + err.Error()
		fmt.Println(errMsg)
		fmt.Fprintln(w, errMsg)
	}

	h.handleConnection(c)
}

func (h *wsHandler) handleConnection(c *websocket.Conn) {
	var (
		connID      = h.storeConn(c, "")
		logger      = logger.With().Str("conn_id", connID).Logger()
		ctx, cancel = context.WithCancel(context.Background())
		conn        = &Connection{c: c}
	)

	wsPinger{c: conn, logger: logger, cancel: cancel}.startPing()
	conn.c.SetCloseHandler(func(code int, text string) error {
		h.eventCh <- clientDisconnectEvent{connID: connID}
		h.flushConn(connID)
		return nil
	})

	logger.Info().Msg("listening for conn")
	h.readClientMsgs(connID, logger, ctx)
	logger.Info().Msg("stop listening for conn")

	h.flushConn(connID)
}

// review playerName param.
func (h *wsHandler) storeConn(c *websocket.Conn, playerName string) (key string) {
	connKey := shortuuid.New()
	h.conns[connKey] = &playerData{
		wsconn: Connection{c: c},
		name:   playerName,
	}
	return connKey
}

type wsConnMsgData struct {
	msgType int
	raw     []byte
	err     error
}

func listenWsConnMsgs(pd *playerData) chan *wsConnMsgData {
	msgCh := make(chan *wsConnMsgData)

	go func() {
		for {
			mt, raw, err := pd.wsconn.ReadRawMessage()

			logger.Info().
				Str("msg_raw", string(raw)).
				Int("websocket_msg_type", mt).
				Msg("[incoming] client's msg")

			msgCh <- &wsConnMsgData{msgType: mt, raw: raw, err: err}
		}
	}()

	return msgCh
}

func (h *wsHandler) readClientMsgs(connID string, logger zerolog.Logger, ctx context.Context) error {
	pd, ok := h.conns[connID]
	if !ok {
		return errors.New(`connection is not established`)
	}

	msgs := listenWsConnMsgs(pd)

	for {
		select {
		case md := <-msgs:
			if md.err != nil {
				return md.err
			}

			msg, err := msgfactory.make(md.raw)
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
		case <-ctx.Done():
			return nil
		}
	}
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
		Int("websocket_msg_type", websocket.TextMessage). // hardcoded text msg
		Msg("[outgoing] sending msg to client")

	return c.wsconn.SendMsg(msg)
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

type wsPinger struct {
	c      *Connection
	logger zerolog.Logger
	cancel context.CancelFunc
}

func (p wsPinger) startPing() {
	go func() {
		time.Sleep(time.Second * 5)

		pingOk := make(chan struct{})

		p.c.SetPongHandler(func(appData string) error {
			pingOk <- struct{}{}
			return nil
		})

		for {
			time.Sleep(pingInterval)

			p.ping()

			select {
			case <-pingOk:
				continue
			case <-time.After(pingDeadline):
				p.logger.Info().Msgf("ping exceeds deadline: %s, cancelling", pingDeadline.String())
				p.cancel()
				return
			}
		}
	}()
}

func (p wsPinger) ping() {
	go p.c.c.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5))
}
