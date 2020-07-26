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
	defer conn.Close()

	g, _ := game.GetGameByID(*user.GameID)

	defer game.DeleteGame(g)
	player := g.Player1

	if g.Player1 == user {
		player = g.Player2
	}

	if player == nil {
		return
	}

	for {

		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("JSON Parse error info: %#v", err)
			//delete(clients, ws)
			resp := models.GameResp{}
			resp.Type = "opponentLeft"
			player.Conn.WriteJSON(resp)
			game.DeleteUser(user)
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

}
