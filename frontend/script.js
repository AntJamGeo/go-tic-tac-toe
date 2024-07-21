// Home page elements
let homePage;
let ticTacToeHeader;
let joinGameUI;
let nameInputLabel;
let nameInput;
let joinGameButton;

// Game page elements
let gamePage;
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

    homePage = createDiv("home-page", document.body, ["page"]);

    ticTacToeHeader = createHeader("tic-tac-toe-header", homePage, 1, "tic-tac-toe");

    joinGameUI = createDiv("join-game-UI", homePage);

    nameInputLabel = createDiv("name-input-label", joinGameUI);
    nameInputLabel.innerText = "enter a username";

    nameInput = createInput("name-input", joinGameUI, "text", "1");
    nameInput.addEventListener("keydown", function(event) {
        if (event.key === "Enter") {
            joinGameNewUserName();
        }
    })

    joinGameButton = createButton("join-game-button", joinGameUI, "2", "join game", joinGameNewUserName);
}

function loadGamePage() {
    // Clear the page
    document.body.innerHTML = ""

    gamePage = createDiv("game-page", document.body, ["page"]);

    // Set a title "Player vs Opponent"
    gameHeaderDiv = createDiv("game-header-div", gamePage);
    gameHeader = createHeader("game-header", gameHeaderDiv, 1, `${playerName} vs ${opponentName}`);

    // Create grid
    gameArea = createDiv("game-area", gamePage);
    gridContainer = createDiv("grid-container", gameArea);
    grid = createDiv("grid", gridContainer, ["unselectable"]);
    cells = [];
    for (let i = 0; i < 9; i++) {
        const cell = createDiv(`cell${i}`, grid, ["cell"])
        cells.push(cell);
    }

    // Add area for notifications and "play again" buttons for when the game finishes
    notificationArea = createDiv("notification-area", gamePage);
    notificationTextArea = createDiv("notification-text-area", notificationArea);
    notification = createP("notification", notificationTextArea, "loading game...");
}

// If pressing the "back to home" button after game, must disconnect from the game first
function backToHome() {
    socket.send(JSON.stringify({ reqType: "disconnect" }));
    loadHomePage();
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
    socket.send(JSON.stringify({ reqType: "disconnect" }));

    joinNewGameButton.removeEventListener("click", joinGameSameUserName);
    joinNewGameButton.classList.add("button-clicked");
    joinNewGameButton.setAttribute("tabindex", "");

    joinGame();
}

// Main logic with game loop
function joinGame() {
    socket = new WebSocket("ws://localhost:3000/play");

    socket.onopen = function(event) {
        socket.send(JSON.stringify({ reqType: "game-Connect", playerName: playerName }));
    }

    socket.onmessage = function(event) {
        data = JSON.parse(event.data);
        console.log(data);
        switch (data.rspType) {
            case "game-Start":
            case "game-Update":
                gameState = data.gameState;
                yourTurn = data.yourTurn === "true" ? true : false;

                if (data.rspType === "game-Start") {
                    gameID = data.gameID;
                    opponentName = data.opponentName;
                    playerSymbol = yourTurn ? "x" : "o";
                    loadGamePage();
                }

                notification.innerText = yourTurn ? "your turn" : "waiting for opponent..."
                for (let i = 0; i < gameState.length; i++) {
                    updateCell(i, gameState[i], yourTurn);
                }
                break;
            case "game-Won":
                gameState = data.gameState;
                winner = data.winner === "true" ? true : false;
                winningCells = data.cells;

                notification.innerText = winner ? "you win!" : "better luck next time..."

                for (let i = 0; i < gameState.length; i++) {
                    updateCell(i, gameState[i], false);
                }
                updateWinningCells(winningCells);

                notificationButtonArea = createDiv("notification-button-area", notificationArea);
                joinNewGameButton = createButton("join-new-game-button", notificationButtonArea, "1", "join new game", joinGameSameUserName);
                backToHomeButton = createButton("back-to-home-button", notificationButtonArea, "2", "back to home", backToHome);
                break;
        }
    }

    socket.onclose = function(event) {
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
    socket.send(JSON.stringify({ reqType: "game-Move", "cell": cell.id.slice(4) }))
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
function createElement(tagName, idName, parentElement, classes = []) {
    newElement = document.createElement(tagName);
    newElement.id = idName;
    newElement.classList.add(...classes);
    parentElement.appendChild(newElement);
    return newElement;
}

function createDiv(idName, parentElement, classes = []) {
    return createElement("div", idName, parentElement, classes);
}

function createButton(idName, parentElement, tabindex, innerText, listener, classes = []) {
    button = createElement("div", idName, parentElement, ["button", "unselectable"]);
    button.classList.add(...classes);
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

function createHeader(idName, parentElement, level, innerText, classes = []) {
    if (level > 6 || level < 1) {
        throw Error("header level should be between 1 and 6");
    }
    header = createElement(`h${level}`, idName, parentElement, classes);
    header.innerText = innerText;
    return header
}

function createInput(idName, parentElement, type, tabindex = "", classes = []) {
    input = createElement("input", idName, parentElement, classes);
    input.type = type;
    input.setAttribute("tabindex", tabindex);
    return input;
}

function createP(idName, parentElement, innerText, classes = []) {
    p = createElement("p", idName, parentElement, classes);
    p.innerText = innerText;
    return p
}