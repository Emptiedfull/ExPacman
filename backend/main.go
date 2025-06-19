package main

import (
	"fmt"
	"time"
)

var Lobbies = make(map[string]*Lobby)

func main() {

	// Game := InitializeGameState()
	// Game.MoveState[0] = Right
	// Game.MoveState[1] = Down

	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	Game.cancel()
	// }()
	// Game.startGame()
	createLobby()
	time.Sleep(1 * time.Second)
	createLobby()
	fmt.Println(Lobbies)
	fmt.Println("Total lobbies:", len(Lobbies))
	setUpServer()

	// board := ParseBoardString(boardString)
	// fmt.Println(board.visualize())

}
