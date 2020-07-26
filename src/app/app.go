package app

import (
	"context"
	"encoding/json"
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
	app.r.HandleFunc("/{uid:[0-9]+}", playOnline)

	app.r.HandleFunc("/gameId/{[a-zA-Z0-9]+}", gameRequest).Methods("GET")
	app.r.Path("/game/{[a-zA-Z0-9]+}").
		Queries("userId", "{[0-9]+}").
		HandlerFunc(startGame).
		Name("startGame")

	app.r.Path("/{uid:[-a-zA-z0-9]+}").
		Queries("gameId", "{[a-zA-Z0-9]+}").
		HandlerFunc(playWithFriends).
		Name("playWithFriends")

	app.r.Path("/validate").
		Queries("gameId", "{[a-zA-Z0-9]+}").
		HandlerFunc(validateGame).
		Name("validateGame")

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

func startGame(w http.ResponseWriter, r *http.Request) {
	s, err := ws.Upgrade(w, r)
	if err != nil {
		log.Fatal(err)
		return
	}

	userID := r.FormValue("userId")

	user := game.GetUser(userID)

	if user == nil {
		log.Fatal("invalid user id")
		return
	}

	user.Conn = s
	log.Printf("Updated user's websocket connection")

	go ws.Worker(s, user)
}

func playOnline(w http.ResponseWriter, r *http.Request) {
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
	resp.Color = user.C.String()

	fmt.Println(resp)

	user.Conn.WriteJSON(resp)
}

func gameRequest(w http.ResponseWriter, r *http.Request) {

	log.Printf("Recieved a new game request")

	game := game.CreateGame()

	resp := make(map[string]string)
	resp["gameId"] = game.GameID
	JSONResponse(w, http.StatusOK, resp)
	/*
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		} else {
			w.Write(resp)
		}
	*/

}

func playWithFriends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["uid"]

	gameID := r.FormValue("gameId")
	g, ok := game.GetGameByID(gameID)

	resp := make(map[string]string)
	if ok == false {
		log.Fatal("Invalid GameID")
		JSONResponse(w, http.StatusUnprocessableEntity, resp)
	}

	s, err := ws.Upgrade(w, r)
	if err != nil {
		log.Fatal(err)
		JSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	log.Printf("Upgrading %s connection to a WS", userID)
	user := game.CreateUser(userID)
	if user == nil {
		log.Fatal("Unable to create User")
		JSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	user.Conn = s
	game.AddPlayer(g, user)
	for g.State != models.RUNNING {

	}

	if g.Player1 == user {
		resp["color"] = user.C.String()
		resp["opponent"] = g.Player2.UserID
	}
	JSONResponse(w, http.StatusOK, resp)
}

func validateGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.FormValue("gameId")
	_, ok := game.GetGameByID(gameID)
	resp := make(map[string]bool)

	resp["valid"] = ok
	JSONResponse(w, http.StatusOK, resp)
}

//JSONResponse returns a JSON response
func JSONResponse(w http.ResponseWriter, code int, output interface{}) {
	// Convert our interface to JSON
	response, _ := json.Marshal(output)
	// Set the content type to json for browsers
	w.Header().Set("Content-Type", "application/json")
	// Our response code
	w.WriteHeader(code)

	w.Write(response)
}
