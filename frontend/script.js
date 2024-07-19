// Page elements
let playerNameInput;
let joinGameButton;
let gameHeaderDiv;
let gameHeader;

// Data
let gameID;
let playerName;
let opponentName;

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

    const socket = new WebSocket("ws://localhost:3000/play");

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
                const cells = [];
                for (let i = 0; i < 9; i++) {
                    const cell = document.createElement('div');
                    cells.push(cell);
                    cell.id = `cell${i+1}`;
                    cell.classList.add("cell");
                    grid.appendChild(cell);
                }
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