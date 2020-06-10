package models

import (
	"github.com/gorilla/websocket"
	"sync"
)

//Game type
type Game struct {
	GameID  string
	Player1 *User
	Player2 *User
	mux 	 sync.RWMutex
}

//User struct to hold user details
type User struct {
	UserID string
	Color  *string
	GameID *string
	Conn   *websocket.Conn
	mux 	 sync.RWMutex
}

//GameReq Request
type GameReq struct {
	Message string `json:"type"`
	UserID  string `json:"userId"`
}

//GameResp response
type GameResp struct {
	GameID  string `json:"gameId"`
	Message string `json:"type"`
	Color   string `json:"color"`
}

func (u User) GetGame() Game {
	var g Game
	return g
}
