package ws

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type msgFactory interface {
	make(raw []byte) (*Msg, error)
}

func newMsgFactory(s *Server) msgFactory {
	return &serverMsgFactory{s: s}
}

type serverMsgFactory struct {
	s *Server
}

func (m *serverMsgFactory) make(raw []byte) (*Msg, error) {
	var msg *Msg
	if err := json.Unmarshal(raw, &msg); err != nil {
		return &Msg{}, fmt.Errorf(`error unmarshaling json msg: %w`, err)
	}

	// somehow compile player id based on conn id?

	switch msg.Kind {
	case setUserDataKind:
		return m.makeSetUserDataMsg(msg.Data)
	case playerRdyKind:
		return m.makePlayerRdy()
	case "":
		return nil, ErrEmptyMsgKind
	}

	return nil, ErrUnsupportedMsgKind
}

func (m *serverMsgFactory) makePlayerRdy() (*Msg, error) {
	return &Msg{Kind: playerRdyKind}, nil
}

func (m *serverMsgFactory) makeSetUserDataMsg(data interface{}) (*Msg, error) {
	temp, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf(`error marshaling data: %w`, err)
	}

	var dest SetUserDataMsg
	err = jsoniter.Unmarshal(temp, &dest)
	if err != nil {
		return nil, fmt.Errorf(`error unmarshaling data: %w`, err)
	}

	return &Msg{Kind: setUserDataKind, Data: dest}, nil
}
