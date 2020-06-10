package game

import (
	"errors"
	"sync"
	"time"

	"github.com/MadhanRaj96/chess-go/src/models"
)

//Engine maintains current games and ready players queue.

type gameQueue struct {
	q   []string
	mux sync.RWMutex
}

type userQueue struct {
	q   []models.User
	mux sync.RWMutex
}

type readyQueue struct {
	q   []userChan
	mux sync.RWMutex
}
type userChan struct {
	u *models.User
	c chan bool
}

var readyQ readyQueue
var userQ userQueue

func matchMaker() {
	for {

	}
}

/*no ready queue change user queue to map -> userID : Object*/
func dequeueUser(user *models.User) {

}

func createUser(u string) *models.User {
	user := models.User{UserID: u}
	userQ.q = append(userQ.q, user)
	return &user
}

//RegisterUser registers the user to the ready queue
func RegisterUser(u string) (*models.User, error) {
	//e.ready = append(e.ready, u)
	user := createUser(u)
	ch := make(chan bool)
	uc := userChan{u: user, c: ch}
	err := errors.New("No player found")
	readyQ.mux.Lock()
	readyQ.q = append(readyQ.q, uc)
	readyQ.mux.Unlock()

	defer dequeueUser(user)
	select {
	case res := <-ch:
		if res == true {

			return user, nil
		}
		return user, err

	case <-time.After(10 * time.Second):
		return user, err
	}

}
