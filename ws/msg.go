package ws

import (
	"encoding/json"
	"tic-tac-toe/game"
)

// msgs can be server or client only
// separate through pkgs?

type Msg struct {
	Kind string
	Data interface{}
}

func NewSetUserDataMsg(name string) *Msg {
	return &Msg{
		Kind: setUserDataKind,
		Data: SetUserDataMsg{NewName: name},
	}
}

func NewPlayerRdyMsg(name string) *Msg {
	return &Msg{
		Kind: playerRdyKind,
	}
}

type SetUserDataMsg struct {
	NewName string `json:"new_name"`
}

func NewMakeTurnMsg(cellIndex int, playerID string) *MakeTurnMsg {
	return &MakeTurnMsg{
		g:         &game.Game{},
		cellIndex: cellIndex,
		playerID:  playerID,
	}
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

// ======================TESTING======================

type MakeTurnMsg struct {
	g         *game.Game `json:"-"`
	cellIndex int        `json:"cell_index"`
	playerID  string     `json:"player_id"`
}

func newMakeTurnMsg(g *game.Game, m *MakeTurnMsg) *MakeTurnMsg {
	return &MakeTurnMsg{
		g:         g,
		cellIndex: m.cellIndex,
		playerID:  m.playerID,
	}
}

func (m *MakeTurnMsg) process() error {
	// somehow compile player id based on conn id

	result := m.g.MakeTurn(m.cellIndex, m.playerID)

	logger.Debug().Interface("turn_result", result).Msg("a turn has been made")

	if result.Err != nil {
		return result.Err // text: error making turn?
	}

	// panic("not implemented") // TODO: Implement

	return nil
}

func (m *MakeTurnMsg) kind() string {
	return makeTurnKind
}
