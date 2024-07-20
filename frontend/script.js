// Home page elements
let homePage;
let ticTacToeHeader;
let joinGameUI;
let nameInputLabel;
let nameInput;
let joinGameButton;

// Game page elements
let gameHeaderDiv;
let gameHeader;
let gameArea;
let gridContainer;
let grid;
let cells;
let notificationArea;
let notificationTextArea;
let notification;
let notificationButtonArea;
let joinNewGameButton;
let backToHomeButton;

// Data
let socket;
let gameID;
let playerName;
let opponentName;
let playerSymbol;
let gameState;
let yourTurn;
let winningCells;

window.addEventListener("DOMContentLoaded", (event) => {
    loadHomePage();
})

// Page loaders
function loadHomePage() {
    document.body.innerHTML = "";

    homePage = createElement("div", "homePage", document.body);

    ticTacToeHeader = createElement("h1", "ticTacToeHeader", homePage);
    ticTacToeHeader.innerText = "tic-tac-toe";

    joinGameUI = createElement("div", "joinGameUI", homePage);

    nameInputLabel = createElement("div", "nameInputLabel", joinGameUI);
    nameInputLabel.innerText = "enter a username";

    nameInput = createElement("input", "nameInput", joinGameUI);
    nameInput.type = "text";
    nameInput.setAttribute("tabindex", "1");
    nameInput.addEventListener("keydown", function(event) {
        if (event.key === "Enter") {
            joinGameNewUserName();
        }
    })

    joinGameButton = createButton("joinGameButton", joinGameUI, "2", "join game", joinGameNewUserName);
}

function loadGamePage() {
    // Clear the page
    document.body.innerHTML = ""

    gamePage = createElement("div", "gamePage", document.body);

    // Set a title "Player vs Opponent"
    gameHeaderDiv = createElement("div", "gameHeaderDiv", gamePage);
    gameHeader = createElement("h1", "gameHeader", gameHeaderDiv);
    gameHeader.textContent = playerName + " vs " + opponentName;

    // Create grid
    gameArea = createElement("div", "gameArea", gamePage);
    gridContainer = createElement("div", "gridContainer", gameArea)
    grid = createElement("div", "grid", gridContainer);
    grid.classList.add("unselectable");
    cells = [];
    for (let i = 0; i < 9; i++) {
        const cell = createElement("div", `cell${i}`, grid);
        cell.classList.add("cell");
        cells.push(cell);
    }

    // Add area for notifications and "play again" buttons for when the game finishes
    notificationArea = createElement("div", "notificationArea", gamePage);
    notificationTextArea = createElement("div", "notificationTextArea", notificationArea);
    notification = createElement("p", "notification", notificationTextArea);
    notification.innerText = "loading game..."
}

// Functions for joining games
// New user name is set when joining game from home page
function joinGameNewUserName() {
    playerName = nameInput.value.trim();
    if (playerName === "") {
        alert("Please enter a username");
        return;
    }

    nameInput.disabled = true;
    joinGameButton.removeEventListener("click", joinGameNewUserName);
    joinGameButton.classList.add("button-clicked");
    joinGameButton.setAttribute("tabindex", "");

    joinGame();
}

// Same user name is used when joining a game straight after finishing another
function joinGameSameUserName() {
    joinNewGameButton.removeEventListener("click", joinGameSameUserName);
    joinNewGameButton.classList.add("button-clicked");
    joinNewGameButton.setAttribute("tabindex", "");

    joinGame();
}

// Main logic with game loop
function joinGame() {
    socket = new WebSocket("ws://localhost:3000/play");

    socket.onopen = function(event) {
        socket.send(JSON.stringify({ msgType: "connecting", playerName: playerName }));
    }

    socket.onmessage = function(event) {
        data = JSON.parse(event.data);
        console.log(data);
        switch (data.msgType) {
            case "joinGame":
                gameID = data.gameID;
                opponentName = data.opponentName;
                loadGamePage();
                break;
            case "gameStart":
            case "gameUpdate":
                gameState = data.gameState;
                yourTurn = data.yourTurn === "true" ? true : false;

                notification.innerText = yourTurn ? "your turn" : "waiting for opponent..."
                for (let i = 0; i < gameState.length; i++) {
                    updateCell(i, gameState[i], yourTurn);
                }
                if (data.msgType === "gameStart") {
                    playerSymbol = yourTurn ? "x" : "o";
                }
                break;
            case "gameWon":
                gameState = data.gameState;
                winner = data.winner === "true" ? true : false;
                winningCells = data.cells;

                notification.innerText = winner ? "you win!" : "better luck next time..."

                for (let i = 0; i < gameState.length; i++) {
                    updateCell(i, gameState[i], false);
                }
                updateWinningCells(winningCells);
                break;
        }
    }

    socket.onclose = function(event) {
        notificationButtonArea = createElement("div", "notificationButtonArea", notificationArea);
        joinNewGameButton = createButton("joinNewGameButton", notificationButtonArea, "1", "join new game", joinGameSameUserName);
        backToHomeButton = createButton("backToHomeButton", notificationButtonArea, "2", "back to home", loadHomePage);
    }

    socket.onerror = function(event) {

    }
}

// Game updates
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

function updateWinningCells(winningCells) {
    for (let i = 0; i < winningCells.length; i++) {
        cells[winningCells[i]].style.backgroundColor = "aqua";
        cells[winningCells[i]].style.color = "rgb(51, 51, 51)";
    }
}

// Element creation helper functions
function createElement(tagName, idName, parentElement) {
    newElement = document.createElement(tagName);
    newElement.id = idName;
    parentElement.appendChild(newElement);
    return newElement;
}

function createButton(idName, parentElement, tabindex, innerText, listener) {
    button = createElement("div", idName, parentElement);
    button.classList.add("button", "unselectable");
    button.setAttribute("tabindex", tabindex);
    button.innerText = innerText;
    button.addEventListener("click", listener);
    button.addEventListener("keydown", function(event) {
        if (event.key === "Enter") {
            listener();
        }
    })
    return button;
}