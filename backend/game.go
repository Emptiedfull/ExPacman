package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type UserMessage struct {
	Type        string      `json:"type"`
	UserID      string      `json:"id"`
	Content     string      `json:"content"`
	Direction   int         `json:"direction"`
	GameOptions GameOptions `json:"options"`
}

type GameOptions struct {
	Timeout   int `json:"duration"`
	GameSpeed int `json:"game_speed"`
}

type UserAlert struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Message string `json:"name"`
}
type PowerUpAlert struct {
	Type   string `json:"type`
	status bool   `json:"status"`
}

type StartAlert struct {
	Type   string        `json:"type"`
	ID     string        `json:"id"`
	Role   string        `json:"role"`
	Users  []UserInfo    `json:"users"`
	Pacman PacmanUpdate  `json:"pacman"`
	Enemy  []GhostUpdate `json:"enemy"`
}
type MoveState struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Direction move   `json:"direction"`
}

type UserInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Pacman bool   `json:"pacman"`
	Enemy  Enemy  `json:"enemy"`
	Host   bool   `json:"host"`
	You    bool   `json:"you"`
}

type BoardUpdate struct {
	Type   string        `json:"type"`
	ID     string        `json:"id"`
	Board  []string      `json:"board"`
	Scores []int         `json:"scores"`
	Pacman PacmanUpdate  `json:"pacman"`
	Ghost  []GhostUpdate `json:"enemy"`
}

type GhostUpdate struct {
	TargetX int `json:"target_x"`
	TargetY int `json:"target_y"`
	PosX    int `json:"pos_x"`
	PosY    int `json:"pos_y"`
	Dir     int `json:"dir"`
	Name    int `json:"name"`
}

type PacmanUpdate struct {
	TargetX int `json:"target_x"`
	TargetY int `json:"target_y"`
	PosX    int `json:"pos_x"`
	PosY    int `json:"pos_y"`
	Dir     int `json:"dir"`
}

type UserError struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Error string `json:"error"`
}

type UserInfoUpdate struct {
	Type  string     `json:"type"`
	ID    string     `json:"id"`
	Users []UserInfo `json:"users"`
}

type HostUpdate struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func listenForMessages(conn *websocket.Conn, user *User, lobby *Lobby) (err bool) {
	// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Welcome %s to lobby %s", user.Name, lobby.ID)))
	GenMessage := GenMessage{
		Type:    "Message",
		ID:      lobby.ID,
		Message: fmt.Sprintf("User %s has joined the lobby", user.Name),
		Sender:  "Server",
	}

	userInfoUpdate := UserInfoUpdate{
		Type:  "UserInfoUpdate",
		ID:    lobby.ID,
		Users: make([]UserInfo, 0, len(lobby.Users)),
	}
	for _, u := range lobby.Users {
		userInfoUpdate.Users = append(userInfoUpdate.Users, UserInfo{
			ID:     u.ID,
			Name:   u.Name,
			Pacman: u.pacman,
			Enemy:  u.Enemy,
			Host:   u.Host,
		})
	}

	sendUserInfoUpdate(userInfoUpdate, lobby)

	// jsonbyytes, error := json.Marshal(userInfoUpdate)
	// if error != nil {
	// 	fmt.Println("Marshal error:", err)

	// }
	// lobby.broadcast <- jsonbyytes

	jsonBytes, error := json.Marshal(GenMessage)
	if error != nil {
		fmt.Println("Marshal error:", err)
	}
	lobby.broadcast <- jsonBytes
	fmt.Printf("User %s connected to lobby %s\n", user.Name, lobby.ID)
	for {
		var msg UserMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read error:", err)
			return true
		}

		if msg.Type == "MoveState" {
			if lobby.GameState == nil {
				fmt.Println("Game has not started yet")
				user.socket.WriteMessage(websocket.TextMessage, []byte("Game has not started yet"))
				continue
			}

			if msg.Direction > 4 {
				fmt.Println("Invalid move direction")
				erorr := UserError{
					Type:  "Error",
					ID:    user.ID,
					Error: "Invalid Move Direction",
				}
				jsonBytes, err := json.Marshal(erorr)
				if err != nil {
					fmt.Println("Error marshaling")
					continue
				}
				lobby.broadcast <- jsonBytes
				continue
			}

			if user.pacman {
				lobby.GameState.MoveState[0] = move(msg.Direction)
			}

			x := int(user.Enemy)
			if x <= 4 || x != 0 {
				lobby.GameState.MoveState[x] = move(msg.Direction)
			}

		}

		if msg.Type == "StartGame" {
			if lobby.GameState != nil {
				continue
			}
			if lobby.host != user {
				fmt.Println("Only the host can start the game")
				user.socket.WriteMessage(websocket.TextMessage, []byte("Only the host can start the game"))
				continue
			}

			if len(lobby.Users) <= 1 {
				fmt.Println("Not enough Users to start game")
				userError := UserError{
					Type:  "Error",
					Error: "Not enough users",
					ID:    lobby.ID,
				}
				user.socket.WriteJSON(userError)

			}

			lobby.GameState = InitializeGameState(len(lobby.Users))
			lobby.AssignRoles()

			userInfoUpdate := UserInfoUpdate{
				Type:  "UserInfoUpdate",
				ID:    lobby.ID,
				Users: make([]UserInfo, 0, len(lobby.Users)),
			}

			for _, u := range lobby.Users {
				userInfoUpdate.Users = append(userInfoUpdate.Users, UserInfo{
					ID:     u.ID,
					Name:   u.Name,
					Pacman: u.pacman,
					Enemy:  u.Enemy,
					Host:   u.Host,
				})
			}

			sendUserInfoUpdate(userInfoUpdate, lobby)

			// jsonbytes, err := json.Marshal(userInfoUpdate)
			// if err != nil {
			// 	fmt.Println("Error marshaling user info update:", err)
			// // }

			// fmt.Println(msg)

			// lobby.broadcast <- jsonbytes

			if msg.GameOptions.GameSpeed == 0 {
				msg.GameOptions.GameSpeed = 250
			}
			if msg.GameOptions.Timeout == 0 {
				msg.GameOptions.Timeout = 10
			}

			go lobby.GameState.startGame(lobby, msg.GameOptions)
		}

	}

	fmt.Println("Message Disconnectd:", user.Name)
	return false
}

func (Lobby *Lobby) AssignRoles() {
	pacman := randHost(Lobby.Users)
	pacman.pacman = true
	pacman.Score = 0
	i := 1
	for _, v := range Lobby.Users {
		if v != pacman {
			v.pacman = false
			v.Enemy = Enemy(i)
			v.Score = 200
			i++
		}

	}

	for _, u := range Lobby.Users {
		fmt.Println("User:", u.Name, "Pacman:", u.pacman, "Enemy:", u.Enemy)
	}

	startAlert := StartAlert{

		Type:  "StartAlert",
		ID:    Lobby.ID,
		Users: make([]UserInfo, 0),
		Pacman: PacmanUpdate{
			TargetX: Lobby.GameState.PlayerPositions[0][0],
			TargetY: Lobby.GameState.PlayerPositions[0][1],
			PosX:    Lobby.GameState.PlayerPositions[0][0],
			PosY:    Lobby.GameState.PlayerPositions[0][1],
			Dir:     3,
		}, Enemy: make([]GhostUpdate, 0),
	}

	for _, User := range Lobby.Users {
		if User.Enemy != 0 {
			startAlert.Enemy = append(startAlert.Enemy, GhostUpdate{
				TargetX: Lobby.GameState.PlayerPositions[User.Enemy][0],
				TargetY: Lobby.GameState.PlayerPositions[User.Enemy][1],
				PosX:    Lobby.GameState.PlayerPositions[User.Enemy][0],
				PosY:    Lobby.GameState.PlayerPositions[User.Enemy][1],
				Dir:     3,
				Name:    int(User.Enemy),
			})
		}

		startAlert.Users = append(startAlert.Users, UserInfo{
			ID:     User.ID,
			Name:   User.Name,
			Pacman: User.pacman,
			Enemy:  User.Enemy,
			Host:   User.Host,
		})

	}

	fmt.Println("starting game with users:", startAlert.Users)

	for _, User := range Lobby.Users {
		if User.pacman {
			startAlert.Role = fmt.Sprintf("%d", 0)
		} else {
			startAlert.Role = fmt.Sprintf("%d", User.Enemy)
		}
		jsonBytes, err := json.Marshal(startAlert)
		if err != nil {
			fmt.Println("Error marshaling,201", err)
		}
		User.socket.WriteMessage(websocket.TextMessage, jsonBytes)

	}

}

// func (alert StartAlert) toJSON() []byte {
// 	jsonBytes, err := json.Marshal(alert)
// 	if err != nil {
// 		fmt.Println("Marshal error:", err)
// 		return nil
// 	}
// 	return jsonBytes
// }

// func (Lobby *Lobby) generateStartAlert() StartAlert {
// 	users := make([]UserInfo, 0, len(Lobby.Users))
// 	for _, user := range Lobby.Users {
// 		users = append(users, UserInfo{
// 			ID:     user.ID,
// 			Name:   user.Name,
// 			Pacman: user.pacman,
// 			Enemy:  user.Enemy,
// 			Host:   user.Host,
// 		})
// 	}
// 	return StartAlert{
// 		Type:  "StartAlert",
// 		ID:    Lobby.ID,
// 		Role:  "Pacman",
// 		Users: users,
// 	}
// }

type GameEndMessage struct {
	Type   string         `json:"type"`
	Winner string         `json:"winner"`
	Scores map[string]int `json:"scores"`
}

func (l *Lobby) handleGameEnd() {
	scores := make(map[string]int)

	pacmanScore := 0
	ghostScore := 0

	for _, user := range l.Users {
		scores[user.Name] = user.Score
		if user.pacman {
			pacmanScore = user.Score
		} else {
			ghostScore += user.Score
		}
	}

	gameEndMessage := GameEndMessage{
		Type:   "GameEnd",
		Scores: scores,
	}

	var ghostmean int
	if ghostScore != 0 {
		ghostmean = ghostScore / (len(l.Users) - 1)
	} else {
		ghostmean = 0
	}

	if pacmanScore > ghostmean {
		gameEndMessage.Winner = "pacman"
	} else {
		gameEndMessage.Winner = "ghosts"
	}

	jsonBytes, err := json.Marshal(gameEndMessage)
	if err != nil {
		fmt.Println("Error marshaling game end message:", err)
		return
	}

	l.broadcast <- jsonBytes
	fmt.Println("Game ended, winner:")
	Lobbies[l.ID] = nil
	delete(Lobbies, l.ID)
}

func (g *GameState) startGame(Lobby *Lobby, options GameOptions) (result bool) {

	fmt.Println(Lobby.GameState.PlayerPositions)

	gameOverTicker := time.NewTicker(time.Duration(options.Timeout) * time.Second)
	defer gameOverTicker.Stop()
	ticker := time.NewTicker(time.Duration(options.GameSpeed) * time.Millisecond)
	defer ticker.Stop()

	var PowerUpTicker *time.Ticker
	powerUpChan := make(chan []int)

	for {
		select {
		case <-gameOverTicker.C:
			defer Lobby.handleGameEnd()
			return true

		case <-PowerUpTicker.C:
			PowerUpTicker.Stop()
			PowerUpTicker = nil
			powerUpAlert := &PowerUpAlert{
				Type:   "powerUp",
				status: false,
			}
			jsonbytes, err := json.Marshal(powerUpAlert)
			if err != nil {
				fmt.Println("error marshaling 427", err)
			} else {
				Lobby.broadcast <- jsonbytes
			}

		case <-powerUpChan:
			if PowerUpTicker == nil {
				powerUpAlert := &PowerUpAlert{
					Type:   "powerUp",
					status: true,
				}
				jsonbytes, err := json.Marshal(powerUpAlert)
				if err != nil {
					fmt.Println("marshalingerr 425", err)
				} else {
					Lobby.broadcast <- jsonbytes
				}

				PowerUpTicker = time.NewTicker(time.Duration(options.Timeout) / 8 * time.Second)
			} else {
				PowerUpTicker.Stop()
				PowerUpTicker = time.NewTicker(time.Duration(options.Timeout) / 8 * time.Second)
			}

		case <-ticker.C:
			if g.gametick(Lobby) {
				return true
			}

			BoardUpdate := BoardUpdate{
				Type:   "BoardUpdate",
				ID:     Lobby.ID,
				Board:  g.Board.visualize(),
				Scores: make([]int, len(Lobby.Users)),
				Pacman: PacmanUpdate{
					TargetX: g.PlayerPositions[0][0],
					TargetY: g.PlayerPositions[0][1],
					Dir:     int(g.MoveState[0]),
				},
				Ghost: make([]GhostUpdate, 0),
			}
			for _, user := range Lobby.Users {
				if user.pacman {
					BoardUpdate.Scores[0] = user.Score
				} else {
					BoardUpdate.Scores[int(user.Enemy)] = user.Score
				}
			}

			for i, pos := range g.PlayerPositions[1:] {
				BoardUpdate.Ghost = append(BoardUpdate.Ghost, GhostUpdate{
					TargetX: pos[0],
					TargetY: pos[1],
					Dir:     int(g.MoveState[i+1]),
					Name:    int(i + 1),
				})
			}

			jsonBytes, err := json.Marshal(BoardUpdate)
			if err != nil {
				fmt.Println("Marshal error:", err)
				return false
			}

			Lobby.broadcast <- jsonBytes
		case <-g.ctx.Done():
			fmt.Println("Game stopped")
			return false
		}
	}

}
