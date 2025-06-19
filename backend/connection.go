package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type NameResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GenMessage struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	path := r.URL.Path
	if len(path) >= 4 && path[:4] == "/ws/" {
		path = path[4:]
	}

	if path == "" || path == "/" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "WebSocket path cannot be empty")
		return
	}

	fmt.Println("LobbyId", path)
	fmt.Println(Lobbies)

	lobby := Lobbies[path]
	if lobby == nil {
		fmt.Println("Lobby not found:", path)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Lobby not found")
		return
	}

	if len(lobby.Users) >= 5 {
		fmt.Println("Lobby is full:", path)
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintln(w, "Lobby is at max capacity")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var nameResponse NameResponse
	if err := conn.ReadJSON(&nameResponse); err != nil {
		fmt.Println("ReadJSON error:", err)
		return
	}
	if nameResponse.ID == "" || nameResponse.Name == "" {
		fmt.Println("Invalid name response:", nameResponse)
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid name response"))
		return
	}
	if nameResponse.ID != lobby.ID {
		fmt.Println("ID mismatch:", nameResponse.ID, "expected:", lobby.ID)
		conn.WriteMessage(websocket.TextMessage, []byte("ID mismatch"))
		return
	}

	addUser(lobby, conn, nameResponse.Name)

}

type User struct {
	ID     string
	Name   string
	socket *websocket.Conn
	Host   bool
	pacman bool
	Enemy  Enemy
}

type Lobby struct {
	Name      string
	ID        string
	Users     map[string]*User
	broadcast chan []byte
	GameState *GameState
	host      *User
}

func randID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}

func createLobby() *Lobby {
	ID := randID()
	lobby := &Lobby{
		ID:        strconv.Itoa(ID),
		Users:     make(map[string]*User),
		broadcast: make(chan []byte, 10),
	}
	Lobbies[lobby.ID] = lobby
	go broadcast(lobby)
	return lobby
}

func broadcast(Lobby *Lobby) {
	for {
		msg := <-Lobby.broadcast
		for _, user := range Lobby.Users {
			if user.socket != nil {
				err := user.socket.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Println("WriteJSON error:", err)
				}
			}
		}
	}
}

func addUser(lobby *Lobby, conn *websocket.Conn, name string) {
	user := &User{
		ID:     strconv.Itoa(randID()),
		Name:   name,
		socket: conn,
	}
	lobby.Users[user.ID] = user
	if lobby.host == nil {
		lobby.host = user
		user.Host = true

	}

	err := listenForMessages(conn, user, lobby)
	if err {
		delete(lobby.Users, user.ID)
		if len(lobby.Users) == 0 {
			Lobbies[lobby.ID] = nil
			delete(Lobbies, lobby.ID)
		} else {
			lobby.host = randHost(lobby.Users)
			fmt.Printf("New host for lobby %s is %s\n", lobby.ID, lobby.host.Name)
		}

	}

}

func randHost(Users map[string]*User) *User {
	keys := make([]string, 0, len(Users))
	for k := range Users {
		keys = append(keys, k)
	}
	rand.Seed(time.Now().UnixNano())
	randomKey := keys[rand.Intn(len(keys))]
	return Users[randomKey]
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func setUpServer() {
	http.Handle("/lobbies", enableCORS(http.HandlerFunc(getLobbies)))
	http.Handle("/ws/", enableCORS(http.HandlerFunc(wsHandler)))
	fmt.Println("WebSocket server started n :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}
