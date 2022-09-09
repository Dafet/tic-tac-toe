package ws

import (
	"encoding/json"
	"fmt"
	"tic-tac-toe/game"
	"tic-tac-toe/game/mark"

	jsoniter "github.com/json-iterator/go"
)

// msgs can be server or client only - separate through pkgs?

type Msg struct {
	Type string
	Data interface{}
}

func NewSetUserDataMsg(name string) *Msg {
	return &Msg{
		Type: setPlayerDataType,
		Data: SetUserDataMsg{NewName: name},
	}
}

func NewPlayerRdyMsg(name string) *Msg {
	return &Msg{
		Type: playerRdyType,
	}
}

type SetUserDataMsg struct {
	NewName string `json:"new_name"`
}

// private?
func newGameStartMsg(firstTurn bool, id string) *Msg {
	var mark mark.Mark
	if firstTurn {
		mark = game.Player1Mark
	} else {
		mark = game.Player2Mark
	}

	return &Msg{
		Type: GameStartType,
		Data: GameStartMsg{
			FirstTurn: firstTurn,
			GameID:    id,
			Mark:      mark,
		},
	}
}

type GameStartMsg struct {
	FirstTurn bool      `json:"first_turn"`
	GameID    string    `json:"game_id"`
	Mark      mark.Mark `json:"mark"`
}

func NewMakeTurnMsg(cellIndex int, gameID string) *Msg {
	return &Msg{
		Type: MakeTurnType,
		Data: &MakeTurnMsg{
			CellIndex: cellIndex,
			GameID:    gameID,
		},
	}
}

type MakeTurnMsg struct {
	CellIndex int    `json:"cell_index"`
	GameID    string `json:"game_id"`
}

func newWaitingTurnMsg(id gameID, gr game.Grid) *Msg {
	return &Msg{
		Type: WaitingTurnType,
		Data: WaitingTurnMsg{
			GameID: id.str(),
			Grid:   gr,
		},
	}
}

type WaitingTurnMsg struct {
	GameID string    `json:"game_id"`
	Grid   game.Grid `json:"game_grid"`
}

func newGameFinishedWinMsg(id gameID, fld game.Grid) *Msg {
	return &Msg{
		Type: GameFinishedType,
		Data: GameFinishedMsg{GameID: id.str(), PlayerWon: true, Field: fld},
	}
}

func newGameFinishedDefeatMsg(id gameID, fld game.Grid) *Msg {
	return &Msg{
		Type: GameFinishedType,
		Data: GameFinishedMsg{GameID: id.str(), PlayerWon: false, Field: fld},
	}
}

func newGameFinishedDrawMsg(id gameID, fld game.Grid) *Msg {
	return &Msg{
		Type: GameFinishedType,
		Data: GameFinishedMsg{GameID: id.str(), PlayerWon: false, IsDraw: true, Field: fld},
	}
}

func newGameFinishedDisconnectMsg(id gameID, fld game.Grid) *Msg {
	return &Msg{
		Type: GameFinishedType,
		Data: GameFinishedMsg{GameID: id.str(), OpponentDisconnect: true, Field: fld},
	}
}

type GameFinishedMsg struct {
	GameID             string    `json:"game_id"`
	PlayerWon          bool      `json:"player_won"`
	IsDraw             bool      `json:"is_draw"`
	OpponentDisconnect bool      `json:"opponent_disconnect"`
	Field              game.Grid `json:"game_grid"`
}

func newErrorMsg(kind, desc string) *Msg {
	return &Msg{
		Type: kind,
		Data: ErrorMsg{
			Desc: desc,
		},
	}
}

// add fields?
type ErrorMsg struct {
	Desc string `json:"desc"`
}

func serializeMsg(msg *Msg) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if data, err = json.Marshal(msg); err != nil {
		return []byte{}, err
	}

	return data, nil
}

func DeserializeMsg(raw []byte) (*Msg, error) {
	var (
		m   *Msg
		err error
	)

	if err = json.Unmarshal(raw, &m); err != nil {
		return &Msg{}, err
	}

	return m, nil
}

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

	switch msg.Type {
	case setPlayerDataType:
		return m.makeSetUserDataMsg(msg.Data)
	case playerRdyType:
		return m.makePlayerRdyMsg()
	case MakeTurnType:
		return m.makeTurnMsg(msg.Data)
	case "":
		return nil, errEmptyMsgType
	}

	return nil, errUnsupportedMsgType
}

func (m *serverMsgFactory) makePlayerRdyMsg() (*Msg, error) {
	return &Msg{Type: playerRdyType}, nil
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

	return &Msg{Type: setPlayerDataType, Data: dest}, nil
}

func (m *serverMsgFactory) makeTurnMsg(data interface{}) (*Msg, error) {
	temp, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf(`error marshaling data: %w`, err)
	}

	var dest MakeTurnMsg
	err = jsoniter.Unmarshal(temp, &dest)
	if err != nil {
		return nil, fmt.Errorf(`error unmarshaling data: %w`, err)
	}

	return &Msg{Type: MakeTurnType, Data: dest}, nil
}
