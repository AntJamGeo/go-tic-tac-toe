// Page elements
let playerNameInput;
let joinGameButton;
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

                // Set a title "Player vs Opponent"
                gameHeader = document.createElement("h1");
                gameHeader.textContent = playerName + " vs " + opponentName;
                document.body.appendChild(gameHeader);
            
                // Create grid
                const gridContainer = document.createElement('div');
                gridContainer.id = "grid";
                document.body.appendChild(gridContainer);
                const cells = [];
                for (let i = 0; i < 9; i++) {
                    const cell = document.createElement('div');
                    cells.push(cell);
                    cell.id = `cell${i+1}`;
                    cell.classList.add('grid-item');
                    gridContainer.appendChild(cell);
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