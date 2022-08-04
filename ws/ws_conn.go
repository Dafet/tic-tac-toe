package ws

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

// todo: websocket.CloseMessage <- use me

type Connection struct {
	c *websocket.Conn
}

// todo: make buffer chan as return?
// Msg type as a return?
func (c *Connection) ListenForServer() <-chan *Msg {
	// todo:
	// - make proper graceful shutdown?
	// - add logs?

	// process interrupt?
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	msgChan := make(chan *Msg)

	go func() {
		defer close(msgChan)
		for {
			logger.Debug().Msg("reading from server conn")

			// websocket.TextMessage
			_, raw, err := c.c.ReadMessage()
			// logger.Debug().Msgf("type is %v", t)

			if err != nil {
				// todo: what to do in case of errors
				logger.Warn().Err(err).Msg("error reading msg from server")
				break
			}

			msg, err := deserializeMsg(raw)
			if err != nil {
				logger.Error().Err(err).Msg("error reading server msg")
				continue
			}

			msgChan <- msg
		}
	}()

	return msgChan
}

func (c *Connection) SendMsg(m *Msg) error {
	data, err := serializeMsg(m)
	if err != nil {
		return fmt.Errorf(`error serializing message: %w`, err)
	}

	return c.c.WriteMessage(websocket.TextMessage, data)
}

func (c *Connection) Close() error {
	if c.c == nil {
		return nil
	}

	// todo: review func

	// msg := []byte("interrupt signal")

	// err := c.c.WriteMessage(websocket.CloseNormalClosure, msg)

	// websocket.FormatCloseMessage()

	// if err != nil {
	// 	return fmt.Errorf(`error sending CloseNormalClosure msg: %w`, err)
	// }

	// cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "add your message here")
	if err := c.c.WriteMessage(websocket.CloseMessage, []byte("normal close")); err != nil {
		// handle error
		logger.Error().Err(err).Msg("error writing close msg before closing connection")
	}

	return c.c.Close()
}

// review
func Connect(addr string) (*Connection, error) {
	// flag.Parse()
	// log.SetFlags(0)

	// validate addr

	// playerQuery := fmt.Sprintf("%s=%s", playerNameParam, playerName)

	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   upgradeConnHandlerPath,
	}
	// log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &Connection{}, err
	}
	// defer c.Close()

	conn := Connection{
		c: c,
	}

	return &conn, nil
}

// review double marshal
func MakeGameStartMsg(msg *Msg) (*GameStartMsg, error) {
	raw, err := json.Marshal(msg.Data)
	if err != nil {
		return &GameStartMsg{}, fmt.Errorf(`error marshaling: %w`, err)
	}

	var dest *GameStartMsg
	if err = json.Unmarshal(raw, &dest); err != nil {
		return &GameStartMsg{}, fmt.Errorf(`error unmarshaling: %w`, err)
	}

	return dest, nil
}
