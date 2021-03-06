package game

import (
	"errors"
	"log"

	"github.com/MadhanRaj96/chess-go/src/models"
	"github.com/MadhanRaj96/chess-go/src/utils"
)

var gameQ gameQueue
var userQ userQueue
var readyQ readyQueue

//Init the game
func Init() {
	gameQ = gameQueue{m: make(map[string]*models.Game)}
	userQ = userQueue{m: make(map[string]*models.User)}
	go matchmaker()
}

func (r *readyQueue) pop() *userChan {
	if len(r.q) > 0 {
		u := r.q[0]
		r.q = r.q[1:]
		return &u
	}
	return nil
}

//CreateUser creates a new User
func CreateUser(u string) *models.User {
	user := models.User{UserID: u}
	userQ.mux.Lock()
	userQ.m[u] = &user
	userQ.mux.Unlock()
	return &user
}

//DeleteUser deletes an existing User
func DeleteUser(u *models.User) {
	userQ.mux.Lock()
	delete(userQ.m, u.UserID)
	userQ.mux.Unlock()
	g, _ := GetGameByID(*u.GameID)
	g.Mux.Lock()
	if g.Player1 == u {
		g.Player1 = nil
	} else {
		g.Player2 = nil
	}
	g.Mux.Unlock()
	*u = models.User{}
}

//CreateGame creates a new game
func CreateGame() *models.Game {
	gameID := utils.GenerateGameID()
	game := models.Game{GameID: gameID, Player1: nil, Player2: nil, State: models.NEW}

	gameQ.mux.Lock()
	gameQ.m[gameID] = &game
	gameQ.mux.Unlock()

	return &game
}

//DeleteGame deletes a game
func DeleteGame(game *models.Game) {
	if game.Player1 == nil && game.Player2 == nil {
		gameQ.mux.Lock()
		delete(gameQ.m, game.GameID)
		gameQ.mux.Unlock()
		*game = models.Game{}
	}
}

//AddPlayer adds a player to the game
func AddPlayer(game *models.Game, user *models.User) {
	game.Mux.Lock()
	defer game.Mux.Unlock()
	if game.State == models.NEW {
		game.Player1 = user
		game.State = models.WAITING
		user.Mux.Lock()
		user.GameID = &game.GameID
		user.C = utils.GetColor()
		user.Mux.Unlock()
	} else if game.State == models.WAITING {
		game.Player2 = user
		game.State = models.RUNNING
		user.Mux.Lock()
		user.GameID = &game.GameID
		if game.Player1.C == models.WHITE {
			user.C = models.BLACK
		} else {
			user.C = models.WHITE
		}
		user.Mux.Unlock()
	}

}

/*
func createGame(user1 *models.User, user2 *models.User) {

	gameID := utils.GenerateGameID()
	game := models.Game{GameID: gameID, Player1: user1, Player2: user2}

	user1.Mux.Lock()
	user1.GameID = &gameID
	user1.Color = utils.GetColor()
	user1.Mux.Unlock()

	user2.Mux.Lock()
	user2.GameID = &gameID
	if user1.Color == "white" {
		user2.Color = "black"
	} else {
		user2.Color = "white"
	}

	user2.Mux.Unlock()

	gameQ.mux.Lock()
	gameQ.m[gameID] = &game
	gameQ.mux.Unlock()

	log.Printf("Creating game %s player 1: %s player 2:%s", gameID, user1.UserID, user2.UserID)

}
*/

//RegisterUser registers user with the matchmaker
func RegisterUser(user *models.User) error {

	err := errors.New("No player found")

	readyQ.mux.Lock()
	ch := make(chan bool)
	uc := userChan{u: user, ch: ch}
	readyQ.q = append(readyQ.q, uc)
	readyQ.mux.Unlock()

	select {
	case res := <-ch:
		if res {
			return nil
		}
		return err
	}
}

//GetUser gets user by userID
func GetUser(u string) *models.User {
	user := userQ.m[u]
	return user
}

//GetGameByID returns game by ID
func GetGameByID(g string) (*models.Game, bool) {
	game, ok := gameQ.m[g]
	return game, ok
}

func matchmaker() {
	for {
		if len(readyQ.q) >= 2 {
			user1 := readyQ.pop()
			user2 := readyQ.pop()
			if user1 == nil || user2 == nil {
				log.Fatal("Internal Error: unable to retrieve user.")
				return
			}
			game := CreateGame()
			AddPlayer(game, user1.u)
			AddPlayer(game, user2.u)
			user1.ch <- true
			user2.ch <- true
		}
	}
}
