function processStartFocus(e) {
    if (event.keyCode == 13) {
        handleStartBtn();
    }
}

function listenStartGameClick() {
    startBtn.addEventListener("click", handleStartBtn);
    againBtn.addEventListener("click", handleStartBtn);
};

function handleStartBtn() {
    console.log("[info] preparing game start");

    hideFinishPopup(); // if present

    grid.reset();
    grid.draw();

    setWaitingAnotherPlayerText();

    ws.sendMsg(new Msg().makeSetPlayerName(playerName));
    ws.sendMsg(new Msg().makePlayerRdyMsg(playerName));

    stopStartClickListening();
}

function stopStartClickListening() {
    startBtn.removeEventListener("click", handleStartBtn);
    againBtn.removeEventListener("click", handleStartBtn);
}

function listenCellClicks() {
    const cellBtn = document.querySelectorAll(".cell")
    cellBtn.forEach(function (btn) {
        btn.addEventListener("click", handleCellClick);
    });
}

function handleCellClick(event) {
    const cellStr = event.target.id.replace("cell-", "");
    const cell = parseInt(cellStr, 10)
    console.log("clicked on grid cell ID: ", cell);

    if (!grid.fillCell(cell, playerMark)) {
        console.log("[debug] clicked on occupied cell");
        animateOccupiedCellClick(event.target.id);

        return
    };
    grid.draw();

    ws.sendMsg(new Msg().makeMakeTurnMsg(cell, gameID));

    stopCellClicksListening();
    setWaitingAnotherPlayerTurnText();
}

function animateOccupiedCellClick(cellID) {
    const cell = document.getElementById(cellID);
    if (cell == undefined) {
        return
    }

    cell.animate([
        { background: "#d85959" },
        { background: "white" },
    ], { duration: 1000, iterations: 1 })
}

function setWaitingAnotherPlayerText() {
    const text = document.getElementById("status-text");
    text.textContent = "looking for opponent...";
    text.style.visibility = "unset";
}

function setYourTurnPlayerText() {
    const text = document.getElementById("status-text");
    text.textContent = "now is your turn, mark is: " + playerMark;
}

function setWaitingAnotherPlayerTurnText() {
    const text = document.getElementById("status-text");
    text.textContent = "waiting for opponent's turn...";
}

function stopCellClicksListening() {
    const cellBtn = document.querySelectorAll(".cell")
    cellBtn.forEach(function (btn) {
        btn.removeEventListener("click", handleCellClick);
    });
}

function setGameFinishedText() {
    const text = document.getElementById("status-text");
    text.textContent = "game has finishied";
}

function setWinFinishText() {
    const text = document.getElementsByClassName("finish-text")[0];
    text.innerHTML = 'game finished: you <span style="color: #4bc447;">WIN</span>'
}

function setLoseFinishText() {
    const text = document.getElementsByClassName("finish-text")[0];
    text.innerHTML = 'game finished: you <span style="color: #b36a3d;">LOSE</span>'
}

function setDrawFinishText() {
    const text = document.getElementsByClassName("finish-text")[0];
    text.innerHTML = 'game finished: <span style="color: #3db38e;">DRAW</span>'
}

function setDisconnectFinishText() {
    const text = document.getElementsByClassName("finish-text")[0];
    text.innerHTML = 'game finished: <span style="color: #2c697d;">opponent disconected</span>'
}

function runFinishPopup() {
    const pop = document.getElementById("finish-popup");
    pop.style.display = "block";

    const elem = document.getElementsByClassName("close")[0];
    elem.onclick = function () {
        pop.style.display = "none";
    }

    window.onclick = function (event) {
        if (event.target == pop) {
            pop.style.display = "none";
        }
    }

    document.getElementById("again-btn").focus();

    listenStartGameClick();
    setGameFinishedText();
}

function hideFinishPopup() {
    const pop = document.getElementById("finish-popup");
    pop.style.display = "none";
}

