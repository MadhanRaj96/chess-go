package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type state int

const (
	NEW state = iota
	WAITING
	RUNNING
	FINISH
)

func (s state) String() string {
	return [...]string{"NEW", "WAITING", "RUNNING", "FINISH"}[s]
}

type Color int

const (
	BLACK Color = iota
	WHITE
)

func (c Color) String() string {
	return [...]string{"black", "white"}[c]
}

//Game type
type Game struct {
	GameID  string
	Player1 *User
	Player2 *User
	State   state
	Mux     sync.RWMutex
}

//User struct to hold user details
type User struct {
	UserID string
	C      Color
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
