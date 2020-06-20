package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

//Game type
type Game struct {
	GameID  string
	Player1 *User
	Player2 *User
	mux     sync.RWMutex
}

//User struct to hold user details
type User struct {
	UserID string
	Color  string
	GameID *string
	Conn   *websocket.Conn
	Mux    sync.RWMutex
}

//GameReq Request
type GameReq struct {
	Message string `json:"type"`
	UserID  string `json:"userId"`
}

//GameResp response
type GameResp struct {
	GameID string `json:"gameId"`
	Color  string `json:"color"`
}

//Move represents a chess Move
type Move struct {
	Type   string `json:"type"`
	GameID string `json:"gameId"`
	State  string `json:"state"`
	Color  string `json:"color"`
}

//GetGame returns game
func (u *User) GetGame() *Game {
	var g Game
	return &g
}
