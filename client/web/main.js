const ws = new WsConn();
const grid = new Grid();

const startBtn = document.getElementById("start-btn");
const againBtn = document.getElementById("again-btn");

const playerName = "player";
let playerMark = "";
let gameID = "";

document.getElementById("start-btn").focus();

document.body.onload = function () {
    listenStartGameClick();
}

ws.setOnMsgCallback(function (e) {
    const data = JSON.parse(e.data);
    const msgData = data.Data;

    console.log("[debug] got server msg:", data);

    switch (data.Type) {
        // process game finished due to another player disconnect!
        case serverType.GAME_START:
            processGameStartType(msgData);
            break;
        case serverType.WAITING_TURN:
            processWaitingTurnKing(msgData);
            break;
        case serverType.ERR_CELL_OCCUPIED:
            console.log("[error] got server ERR_CELL_OCCUPIED error");
            break;
        case serverType.GAME_FINISHED:
            processGameFinishedType(msgData);
            break;
    }
});

function processGameStartType(data) {
    console.log("[info] starting game with id %s...",
        data.game_id)

    gameID = data.game_id;
    playerMark = data.mark;

    console.log("[debug] data.first_turn: ", data.first_turn);

    if (data.first_turn) {
        console.log("[debug] player is going first with: %s", playerMark);

        setYourTurnPlayerText();
        listenCellClicks();

        return
    }

    console.log("[debug] player is going second with: %s", playerMark);

    setWaitingAnotherPlayerTurnText()
}

function processWaitingTurnKing(data) {
    setYourTurnPlayerText();
    listenCellClicks();

    grid.elems = data.game_grid;
    grid.draw();
}

function processGameFinishedType(data) {
    if (data.opponent_disconnect) {
        console.log("[debug] opponent has disconnected");
    }

    grid.elems = data.game_grid;
    grid.draw();

    console.log("[debug] player won: ", data.player_won);

    if (data.opponent_disconnect) { setDisconnectFinishText(); } else if (data.is_draw) { setDrawFinishText(); }
    else if (data.player_won) { setWinFinishText(); }
    else { setLoseFinishText(); }

    runFinishPopup();
}