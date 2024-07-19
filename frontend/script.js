let playerNameInput;
let joinGameButton;

window.addEventListener("DOMContentLoaded", (event) => {
    joinGameButton = document.getElementById("joinGameButton");
    joinGameButton.addEventListener("click", joinGame);
})

function joinGame() {
    playerNameInput = document.getElementById("nameInput");
    const playerName = playerNameInput.value.trim();
    if (playerName === "") {
        alert("Please enter a username");
        return;
    }

    const socket = new WebSocket("ws://localhost:3000/play");

    socket.onopen = function(event) {
        socket.send(JSON.stringify({ msgType: "connecting", playerName: playerName }));
        playerNameInput.disabled = true;
        joinGameButton.removeEventListener("click", joinGame);
        joinGameButton.classList.add("joinGameButtonClicked");
        joinGameButton.style.backgroundColor = "rgb(30, 30, 43)";
        joinGameButton.style.color = "gray";
        joinGameButton.style.cursor = "default";
    }
}