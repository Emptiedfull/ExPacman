package main

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

	setUpServer()

	// board := ParseBoardString(boardString)
	// fmt.Println(board.visualize())

}
