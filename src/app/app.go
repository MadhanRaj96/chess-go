package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MadhanRaj96/chess-go/src/game"
	"github.com/MadhanRaj96/chess-go/src/models"
	"github.com/MadhanRaj96/chess-go/src/ws"
	"github.com/gorilla/mux"
)

//App wraps server with game
type App struct {
	r *mux.Router
}

//Init the application
func (app *App) Init() {
	log.Println("Initializing app")
	app.r = mux.NewRouter().StrictSlash(true)
	app.initializeRoutes()
	game.Init()
}

func (app *App) initializeRoutes() {
	log.Println("Initializing Routes")
	app.r.HandleFunc("/{uid:[0-9]+}", handleConnections)
	app.r.HandleFunc("/gameId/{id:[0-9]+}", gameRequestHandler).Methods("GET")
}

//Run starts the server
func (app *App) Run() {
	log.Println("Running server on 127.0.0.1:8000")
	var wait time.Duration
	srv := &http.Server{
		Addr: "127.0.0.1:8000",

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      app.r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["uid"]
	/*
		user := game.GetUser(userID)
		if user == nil {
			log.Fatal("Invalid USER ID")
			return
		}
	*/
	log.Printf("Upgrading %s connection to a WS", userID)
	s, err := ws.Upgrade(w, r)
	if err != nil {
		log.Fatal(err)
		return
	}
	user := game.CreateUser(userID)
	if user == nil {
		log.Fatal("Unable to create User")
		return
	}
	user.Conn = s

	/*finding opponent*/
	err = game.RegisterUser(user)
	if err != nil {
		log.Fatal("Unable to find opponent")
		return
	}

	go ws.Worker(s, user)

	resp := models.GameResp{}

	resp.GameID = *user.GameID
	resp.Color = user.Color

	fmt.Println(resp)

	user.Conn.WriteJSON(resp)
}

func gameRequestHandler(w http.ResponseWriter, r *http.Request) {
	/*
		vars := mux.Vars(r)
		userID := vars["id"]

		log.Printf("Recieved game request from user: %s", userID)

		user, err := game.RegisterUser(userID)

		resp := m.GameResp{}

		if err != nil {
			log.Println("No player found")
			resp.GameID = ""
			resp.Color = ""
		} else {
			resp.GameID = *user.GameID
			resp.Color = user.Color
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	*/
}
