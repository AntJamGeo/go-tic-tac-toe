# Online Multiplayer Tic-Tac-Toe
An online multiplayer tic-tac-toe game server.

## Build and Run

### Docker
If you have docker, you can run `docker compose build` to build the image. You can then run `docker compose up` to start the server.

### Make + Docker
If you make and docker, you can simply run `make build` to build the image. You can then run `make up` to start the server.

## Accessing the game
The game can be accessed by going to localhost:3000 in your browser. You should see the following screen:

![tic-tac-toe-home-page](https://github.com/AntJamGeo/go-tic-tac-toe/blob/main/images/tic-tac-toe-home-page.png)

When two players click to join a game, they are taken to the game page:

![tic-tac-toe-game-page](https://github.com/AntJamGeo/go-tic-tac-toe/blob/main/images/tic-tac-toe-game-page.png)

When the game concludes, the winning sequence will be highlighted:

![tic-tac-toe-game-page-end](https://github.com/AntJamGeo/go-tic-tac-toe/blob/main/images/tic-tac-toe-game-page-end.png)

A player can then choose to join another game or go back to the home page.

## Cleaning up
If you want to stop and remove your containers and other associated objects created from `up`, you can run the following:

### Docker
Stop and remove your container using `docker compose down`. Remove your image using `docker image rm go-tic-tac-toe`.

### Make + Docker
Stop and remove your container using `make down`, or remove both your container and image using `make clean`.

## Development
Some `.air.toml` files are included for hot reloading to make development easy with air. 

### Docker
You can build the development image by linking to the development compose file; run `docker compose -f docker-compose-dev.yml build`. This will build an image called `go-tic-tac-toe-dev`. The server can then be started in development mode by using `docker compose -f docker-compose-dev.yml up`. Once air has started, you will need to watch for changes by pressing "w". For clean up, you can run `docker compose -f docker-compose-dev.yml down`, and the image can be removed by running `docker image rm go-tic-tac-toe-dev`.

### Make + Docker
This is similar to the instructions for using docker, but with less verbose commands. These are:
* build: `make build-dev`
* compose up: `make up-dev`
* compose down: `make down-dev`
* image rm: `make clean-dev`

## To-Do List
A list of potential things for me to do in the future:

- [x] Handle draws
- [x] Validate that a move is allowed on server side
- [x] Dockerise
- [ ] Shift to a microservice architecture
- [ ] Add ability to cancel matchmaking after clicking to join
- [ ] Add different themes
- [ ] Create other game modes, including game modes for more than two players
- [ ] Add a chat for each game
- [ ] User account creation and authentication
- [ ] User game history and statistics
- [ ] Spectate other games
- [ ] Time-limited games
