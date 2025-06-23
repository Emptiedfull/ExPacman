package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type UserMessage struct {
	Type      string `json:"type"`
	UserID    string `json:"id"`
	Content   string `json:"content"`
	Direction int    `json:"direction"`
}

type UserAlert struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Message string `json:"name"`
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
}

type BoardUpdate struct {
	Type   string        `json:"type"`
	ID     string        `json:"id"`
	Board  []string      `json:"board"`
	Scores [5]int        `json:"scores"`
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

	UserInfoUpdate := UserInfoUpdate{
		Type:  "UserInfoUpdate",
		ID:    lobby.ID,
		Users: make([]UserInfo, 0, len(lobby.Users)),
	}
	for _, u := range lobby.Users {
		UserInfoUpdate.Users = append(UserInfoUpdate.Users, UserInfo{
			ID:     u.ID,
			Name:   u.Name,
			Pacman: u.pacman,
			Enemy:  u.Enemy,
			Host:   u.Host,
		})
	}

	jsonbyytes, error := json.Marshal(UserInfoUpdate)
	if error != nil {
		fmt.Println("Marshal error:", err)

	}
	lobby.broadcast <- jsonbyytes

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
			fmt.Println("Received MoveState message:", msg)
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

			go lobby.GameState.startGame(lobby)
		}

	}

	fmt.Println("Message Disconnectd:", user.Name)
	return false
}

func (Lobby *Lobby) AssignRoles() {
	pacman := randHost(Lobby.Users)
	pacman.pacman = true
	i := 1
	for _, v := range Lobby.Users {
		if v != pacman {
			v.pacman = false
			v.Enemy = Enemy(i)
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

func (alert StartAlert) toJSON() []byte {
	jsonBytes, err := json.Marshal(alert)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return nil
	}
	return jsonBytes
}

func (Lobby *Lobby) generateStartAlert() StartAlert {
	users := make([]UserInfo, 0, len(Lobby.Users))
	for _, user := range Lobby.Users {
		users = append(users, UserInfo{
			ID:     user.ID,
			Name:   user.Name,
			Pacman: user.pacman,
			Enemy:  user.Enemy,
			Host:   user.Host,
		})
	}
	return StartAlert{
		Type:  "StartAlert",
		ID:    Lobby.ID,
		Role:  "Pacman",
		Users: users,
	}
}

func (g *GameState) startGame(Lobby *Lobby) (result bool) {

	fmt.Println(Lobby.GameState.PlayerPositions)
	// ticker := time.NewTicker(1 * time.Second)
	ticker := time.NewTicker(250 * time.Millisecond) // 10 ticks per second
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if g.gametick() {
				return true
			}

			BoardUpdate := BoardUpdate{
				Type:   "BoardUpdate",
				ID:     Lobby.ID,
				Board:  g.Board.visualize(),
				Scores: g.Scores,
				Pacman: PacmanUpdate{
					TargetX: g.PlayerPositions[0][0],
					TargetY: g.PlayerPositions[0][1],
					Dir:     int(g.MoveState[0]),
				},
				Ghost: make([]GhostUpdate, 0),
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
