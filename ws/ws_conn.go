package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"

	"github.com/gorilla/websocket"
)

type Connection struct {
	c *websocket.Conn
}

func (c *Connection) ReadMessage() (mt int, data []byte, err error) {
	return c.c.ReadMessage()
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

	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "normal closure")
	err := c.c.WriteMessage(websocket.CloseMessage, msg)

	if err != nil {
		return err
	}

	return nil
}

func IsCloseError(err error) bool {
	codes := []int{websocket.CloseMessage, websocket.CloseNormalClosure, websocket.CloseGoingAway}
	close := websocket.IsCloseError(err, codes...)
	return close
}

func Connect(addr string) (*Connection, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   upgradeConnHandlerPath,
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &Connection{}, err
	}

	return &Connection{c: c}, nil
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

func MakeWaitingTurnMsg(msg *Msg) (*WaitingTurnMsg, error) {
	var dest WaitingTurnMsg
	if err := makeConcreteMsg(msg, &dest); err != nil {
		return &WaitingTurnMsg{}, err
	}

	return &dest, nil
}

func MakeGameFinishedMsg(msg *Msg) (*GameFinishedMsg, error) {
	var dest GameFinishedMsg
	if err := makeConcreteMsg(msg, &dest); err != nil {
		return &GameFinishedMsg{}, err
	}

	return &dest, nil
}

func makeConcreteMsg(msg *Msg, dest interface{}) error {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return errors.New(`dest must be a pointer`)
	}

	raw, err := json.Marshal(msg.Data)
	if err != nil {
		return fmt.Errorf(`error marshaling: %w`, err)
	}

	if err = json.Unmarshal(raw, &dest); err != nil {
		return fmt.Errorf(`error unmarshaling: %w`, err)
	}

	return nil
}
