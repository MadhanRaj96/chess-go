package game

import (
	"sync"

	"github.com/MadhanRaj96/chess-go/src/models"
)

type gameQueue struct {
	m   map[string]*models.Game
	mux sync.RWMutex
}
type userQueue struct {
	m   map[string]*models.User
	mux sync.RWMutex
}
type userChan struct {
	u  *models.User
	ch chan bool
}
type readyQueue struct {
	q   []userChan
	mux sync.RWMutex
}
