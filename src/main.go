package main

import (
	"github.com/MadhanRaj96/chess-go/src/app"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

func main() {
	app := app.App{}

	app.Init()
	app.Run()

}
