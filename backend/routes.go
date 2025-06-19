package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LobbyResponse struct {
	Name    string `json:"name"`
	Players int    `json:"players"`
	ID      string `json:"id"`
}

func getLobbies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	lobbies := make([]*LobbyResponse, 0, len(Lobbies))
	for _, lobby := range Lobbies {
		lobbies = append(lobbies, &LobbyResponse{
			Name:    lobby.ID,
			Players: len(lobby.Users),
			ID:      lobby.ID,
		})
	}

	jsonBytes, err := json.Marshal(lobbies)
	if err != nil {
		http.Error(w, "Error marshaling lobbies", http.StatusInternalServerError)
		return
	}
	fmt.Println("Returning lobbies:", string(jsonBytes))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
