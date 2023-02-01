# Mancala

A game of mancala with a Golang server and React client. Also used PostgreSQL, WebSocket, RestAPI.

Watch a youtube video about the Mancala game [here](https://www.youtube.com/watch?v=-A-djjimCcM)

# Design and Features

1. Backend is written in Golang, using PostgreSQL as database, and WebSocket and REST API for communication.
2. To play a game, frontend creates a new lobby via HTTP request and then connects to a WebSocket session and plays the
   game.
3. On each turn and event, client and server share websocket messages.
4. On each player turn Mancala state is saved in database.

# How to run

* docker-compose up
    * This will build images for server and client and run them.
    * Server will be available on port 80
    * Client will be available on port 3000

* On browser open http://localhost:3000, click on "create game" and ten join it. You'll see simplified game board.

# Progress

## What's done

1. REST API with lobbies
2. WebSocket connection to handle game events.(unfinished)
3. Game logic and a few tests
4. Database connection and saving game state
5. Dockerfile for server and client
6. Basic frontend(unfinished)

## What to improve

1. Add more tests for game logic. I tested a few cases, but there are many more.
2. Add correct handling of WebSocket messages on frontend
2. Finish frontend(now it basically doesn't work, but you can see it sends correct "player turn" messages to server)
3. WebSocket handler is tightly coupled with implementation. Ideally it's better to split it.

