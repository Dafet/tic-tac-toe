package ws

import (
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
func (c *Connection) ListenForServer() <-chan interface{} {
	// todo:
	// - make proper graceful shutdown?
	// - add logs?

	// process interrupt?
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	msgChan := make(chan interface{})

	go func() {
		defer close(msgChan)
		for {
			logger.Debug().Msg("reading from server conn")

			// websocket.TextMessage
			_, message, err := c.c.ReadMessage()
			// logger.Debug().Msgf("type is %v", t)

			if err != nil {
				// todo: what to do in case of errors
				logger.Warn().Err(err).Msg("error reading msg from server")
				break
			}

			msgChan <- message

			// os.Exit(0)
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
