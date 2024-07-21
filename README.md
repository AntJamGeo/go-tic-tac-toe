# Online Multiplayer Tic-Tac-Toe
An online multiplayer tic-tac-toe game with the backend written in Go.

## Build
To build, clone the repository, navigate to the `backend` directory, and run `go build -o ../server.exe ./cmd`. This will output `server.exe` in the repository's root directory.

## Run
To run, navigate to the root of the repository and run the command `./server.exe`. The game can be accessed by going to `localhost:3000` in your browser. You should see the following screen:

![tic-tac-toe-home-page](https://github.com/AntJamGeo/go-tic-tac-toe/blob/main/images/tic-tac-toe-home-page.png)

When two players click to join a game, they are taken to the game page:

![tic-tac-toe-game-page](https://github.com/AntJamGeo/go-tic-tac-toe/blob/main/images/tic-tac-toe-game-page.png)

When the game concludes, the winning sequence will be highlighted:

![tic-tac-toe-game-page-end](https://github.com/AntJamGeo/go-tic-tac-toe/blob/main/images/tic-tac-toe-game-page-end.png)

A player can then choose to join another game or go back to the home page.

## To-Do List
A list of potential things for me to do in the future:

- [ ] Handle draws
- [ ] Validate that a move is allowed on server side
- [ ] Add ability to cancel matchmaking after clicking to join
- [ ] Add different themes
- [ ] Create other game modes, including game modes for more than two players
- [ ] Add a chat for each game
- [ ] User account creation and authentication
- [ ] User game history and statistics
- [ ] Spectate other games
- [ ] Time-limited games
