package game

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/MadhanRaj96/chess-go/src/models"
	"github.com/MadhanRaj96/chess-go/src/utils"
)

//Engine maintains current games and ready players queue.

type gameQueue struct {
	m   map[string]*models.Game
	mux sync.RWMutex
}

var gameQ = &gameQueue{m: make(map[string]*models.Game)}

type userQueue struct {
	m   map[string]*models.User
	mux sync.RWMutex
}

var userQ = &userQueue{m: make(map[string]*models.User)}

type userChan struct {
	u  *models.User
	ch chan bool
}
type readyQueue struct {
	q   []userChan
	mux sync.RWMutex
}

var readyQ readyQueue

func (r *readyQueue) pop() *userChan {

	if len(r.q) > 0 {
		u := r.q[0]
		r.q = r.q[1:]
		return &u
	}
	return nil
}

func createUser(u string) *models.User {
	user := models.User{UserID: u}
	userQ.mux.Lock()
	userQ.m[u] = &user
	userQ.mux.Unlock()
	return &user
}

func deleteUser(u *models.User) {
	userQ.mux.Lock()
	delete(userQ.m, u.UserID)
	userQ.mux.Unlock()
}

func createGame(user1 *models.User, user2 *models.User) {

	gameID := user1.UserID + utils.RandomString(6)
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

	log.Printf("Creating game %s", gameID)

}

//RegisterUser registers the user to the ready queue
func RegisterUser(u string) (*models.User, error) {
	//e.ready = append(e.ready, u)
	user := createUser(u)
	err := errors.New("No player found")

	if len(readyQ.q) == 0 {
		readyQ.mux.Lock()
		ch := make(chan bool)
		uc := userChan{u: user, ch: ch}
		readyQ.q = append(readyQ.q, uc)
		readyQ.mux.Unlock()

		select {
		case <-ch:
			return user, nil

		case <-time.After(10 * time.Second):
			_ = readyQ.pop()
			deleteUser(user)
			return user, err
		}
	} else {
		user1 := readyQ.pop()
		if user1 != nil {
			createGame(user, user1.u)
			user1.ch <- true
			return user, nil
		}
		deleteUser(user)
		return user, err

	}

}
