// Page elements
let playerNameInput;
let joinGameButton;
let gameHeaderDiv;
let gameHeader;

// Data
let socket;
let gameID;
let playerName;
let opponentName;
let playerSymbol;
let gameState;
let yourTurn;
let cells;

window.addEventListener("DOMContentLoaded", (event) => {
    joinGameButton = document.getElementById("joinGameButton");
    joinGameButton.addEventListener("click", joinGame);
})

function joinGame() {
    playerNameInput = document.getElementById("nameInput");
    playerName = playerNameInput.value.trim();
    if (playerName === "") {
        alert("Please enter a username");
        return;
    }

    socket = new WebSocket("ws://localhost:3000/play");

    socket.onopen = function(event) {
        socket.send(JSON.stringify({ msgType: "connecting", playerName: playerName }));
        disallowJoinGame();
    }

    socket.onmessage = function(event) {
        data = JSON.parse(event.data);
        console.log(data);
        switch (data.msgType) {
            case "joinGame":
                gameID = data.gameID;
                opponentName = data.opponentName;

                // Clear the page
                document.body.innerHTML = ""

                gamePage = document.createElement("div");
                gamePage.id = "gamePage";
                document.body.appendChild(gamePage);

                // Set a title "Player vs Opponent"
                gameHeaderDiv = document.createElement("div");
                gameHeaderDiv.id = "gameHeaderDiv";
                gamePage.appendChild(gameHeaderDiv);
                gameHeader = document.createElement("h1");
                gameHeader.id = "gameHeader";
                gameHeader.textContent = playerName + " vs " + opponentName;
                gameHeaderDiv.appendChild(gameHeader);
            
                // Create grid
                const gameArea = document.createElement("div");
                gameArea.id = "gameArea";
                gamePage.appendChild(gameArea);
                const gridContainer = document.createElement('div');
                gridContainer.id = "gridContainer";
                gameArea.appendChild(gridContainer);
                const grid = document.createElement('div');
                grid.id = "grid";
                gridContainer.appendChild(grid);
                cells = [];
                for (let i = 0; i < 9; i++) {
                    const cell = document.createElement('div');
                    cells.push(cell);
                    cell.id = `cell${i+1}`;
                    cell.classList.add("cell");
                    grid.appendChild(cell);
                }
                break;
            case "gameStart":
            case "gameUpdate":
                gameState = data.gameState;
                yourTurn = data.yourTurn === "true" ? true : false;
                for (let i = 0; i < gameState.length; i++) {
                    updateCell(i, gameState[i], yourTurn);
                }
                if (data.msgType === "gameStart") {
                    playerSymbol = yourTurn ? "x" : "o";
                }
                console.log(playerSymbol);
                break;
        }
    }

    socket.onclose = function(event) {

    }

    socket.onerror = function(event) {

    }
}

function disallowJoinGame() {
    playerNameInput.disabled = true;
    joinGameButton.removeEventListener("click", joinGame);
    joinGameButton.style.backgroundColor = "rgb(30, 30, 43)";
    joinGameButton.style.color = "gray";
    joinGameButton.style.cursor = "default";
}

function makeMove(event) {
    let cell = event.target;
    for (let i = 0; i < gameState.length; i++) {
        cells[i].removeEventListener("click", makeMove);
        cells[i].style.cursor = "";
    }
    socket.send(JSON.stringify({ msgType: "gameUpdate", "cell": cell.id.slice(4) }))
}

function updateCell(i, cellState, yourTurn) {
    switch (cellState) {
        case "x":
            cells[i].textContent = "x";
            break;
        case "o":
            cells[i].textContent = "o";
            break;
        case "-":
            if (yourTurn) {
                cells[i].addEventListener("click", makeMove);
                cells[i].style.cursor = "pointer";
            }
            break;
    }
}