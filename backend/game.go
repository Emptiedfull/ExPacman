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
	Direction int    `json:"direction1"`
}

type UserAlert struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Message string `json:"name"`
}
type StartAlert struct {
	Type  string     `json:"type"`
	ID    string     `json:"id"`
	Role  string     `json:"role"`
	Users []UserInfo `json:"users"`
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
}

type BoardUpdate struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Board  string `json:"board"`
	Scores [5]int `json:"scores"`
}

type UserError struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Error string `json:"error"`
}

func listenForMessages(conn *websocket.Conn, user *User, lobby *Lobby) (err bool) {
	// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Welcome %s to lobby %s", user.Name, lobby.ID)))
	GenMessage := GenMessage{
		Type:    "Message",
		ID:      lobby.ID,
		Message: fmt.Sprintf("User %s has joined the lobby", user.Name),
		Sender:  "Server",
	}
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
			if x <= 4 {
				lobby.GameState.MoveState[x+1] = move(msg.Direction)
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

			lobby.GameState = InitializeGameState()
			lobby.AssignRoles()

			lobby.GameState.startGame(lobby)
		}

	}
}

func (Lobby *Lobby) AssignRoles() {
	pacman := randHost(Lobby.Users)
	pacman.pacman = true
	i := 0
	for _, v := range Lobby.Users {
		if v != pacman {
			v.pacman = false
			v.Enemy = Enemy(i)
		}
		i++
	}

	Lobby.broadcast <- Lobby.generateStartAlert().toJSON()
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
	ticker := time.NewTicker(1 * time.Second)
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
