package websocket

import (
	"log"
	"net/http"
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

	for {
		var msg models.GameReq
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("JSON Parse error info: %#v", err)
			//delete(clients, ws)
			break
		}
		/*send the data recieved to player 2*/
		game := user.GetGame()

        /*
		if err != nil {
			log.Fatalf("User %s Game not found", user.UserID)
		}
        */
		ws := game.Player1
		if game.Player1 == user.Conn {
			ws = game.Player2
		}
		ws.WriteJSON(msg)
	}
	defer conn.Close()
}
