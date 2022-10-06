// const wsURL = "ws://185.195.27.142:8080/conn"
const wsURL = "ws://localhost:8080/conn"

const type = {
    SET_PLAYER_DATA: "set-player-data",
    PLAY_RDY: "play-ready",
    MAKE_TURN: "make-turn",
}

const serverType = {
    GAME_START: "game-start",
    WAITING_TURN: "waiting-player-turn",
    ERR_CELL_OCCUPIED: "cell-occupied",
    GAME_FINISHED: "game-finished",
}

class WsConn {
    ws;

    constructor(onOpenf) {
        this.ws = new WebSocket(wsURL);
        this.ws.onopen = onOpenf;
        this.ws.onerror = function (e) {
            console.log("[error] connecting to ws server: ", e.type);
            alert("error connecting to ws server")
        }
    }

    info() {
        console.log("ws url: ", this.ws.url);
        console.log("ws rdy state: ", this.ws.readyState);
    }

    setOnMsgCallback(f) {
        this.ws.onmessage = f;
    }

    sendMsg(msg) {
        this.ws.send(msg);
    }
}


class Msg {
    makeSetPlayerName(name) {
        const data = {
            type: type.SET_PLAYER_DATA,
            data: {
                new_name: name
            }
        }
        return JSON.stringify(data);
    }

    makePlayerRdyMsg() {
        const data = {
            type: type.PLAY_RDY
        }
        return JSON.stringify(data)
    }

    makeMakeTurnMsg(cellID, gameID) {
        const data = {
            type: type.MAKE_TURN,
            data: {
                cell_index: cellID, game_id: gameID
            }
        }
        return JSON.stringify(data)
    }
}
