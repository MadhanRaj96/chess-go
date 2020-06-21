package ws

import (
	"log"
	"net/http"

	"github.com/MadhanRaj96/chess-go/src/game"
	"github.com/MadhanRaj96/chess-go/src/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Upgrade a http connnection to a websocket
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, err
}

//Worker reads message from a player and sends it to another player
func Worker(conn *websocket.Conn, user *models.User) {
	game := game.GetGameByID(*user.GameID)

	player := game.Player1

	if game.Player1 == user {
		player = game.Player2
	}

	for {

		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("JSON Parse error info: %#v", err)
			//delete(clients, ws)
			break
		}
		log.Printf("message %s", string(message))
		/*send the data received to player 2*/
		if player == nil {
			log.Println("nil player")
		} else {
			log.Printf("sending message to %s", player.UserID)
			player.Conn.WriteMessage(mt, []byte(message))
		}
	}
	defer conn.Close()
}
