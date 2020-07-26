package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	app.r.Path("/validate").
		HandlerFunc(validateGame).
		Name("validateGame")

	app.r.HandleFunc("/gameId/{uid:[a-zA-Z0-9]+}", gameRequest).Methods("GET")

	app.r.Path("/game/{[a-zA-Z0-9]+}").
		Queries("userId", "{[0-9]+}").
		HandlerFunc(startGame).
		Name("startGame")

	app.r.HandleFunc("/{uid:[0-9]+}", playOnline)

	app.r.Path("/room").
		HandlerFunc(playWithFriends).
		Name("playWithFriends")

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
	log.Printf("inside start game")
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

	if user.C == models.WHITE {
		resp := models.GameResp{}
		resp.Type = "ready"
		user.Conn.WriteJSON(resp)
	}

	user.Conn = s
	log.Printf("Updated user's websocket connection")

	go ws.Worker(s, user)
}

func playOnline(w http.ResponseWriter, r *http.Request) {
	log.Printf("inside play online")
	vars := mux.Vars(r)
	userID := vars["uid"]

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
	vars := mux.Vars(r)
	userID := vars["uid"]

	g := game.CreateGame()

	user := game.CreateUser(userID)

	if user == nil || g == nil {
		log.Fatal("error in USER/GAME creation")
	}

	game.AddPlayer(g, user)

	resp := make(map[string]string)
	resp["color"] = user.C.String()
	resp["gameId"] = g.GameID
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
	log.Printf("inside Play with friends")

	gameID := r.FormValue("gameId")
	userID := r.FormValue("userId")
	g, ok := game.GetGameByID(gameID)

	log.Printf("userId: %s gameID: %s", userID, gameID)

	s, err := ws.Upgrade(w, r)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Upgrading %s connection to a WS", userID)
	user := game.GetUser(userID)
	if user == nil {
		log.Fatal("Invalid USER ID")
	}

	user.Conn = s

	for g.State != models.RUNNING {

	}
	resp := models.GameResp{}

	if ok == false {
		log.Fatal("Invalid GameID")
		return
	}
	if g.Player1 == user {
		resp.Opponent = g.Player2.UserID
	} else {
		resp.Opponent = g.Player1.UserID
	}

	user.Conn.WriteJSON(resp)
}

func validateGame(w http.ResponseWriter, r *http.Request) {

	log.Printf("inside validate game")

	gameID := r.FormValue("gameId")
	userID := r.FormValue("userId")

	g, ok := game.GetGameByID(gameID)
	user := game.CreateUser(userID)

	if user == nil {
		log.Fatal("error in USER creation")
	}

	game.AddPlayer(g, user)

	resp := make(map[string]string)

	resp["color"] = user.C.String()
	resp["valid"] = strconv.FormatBool(ok)
	log.Printf("Game ID: %s validation: %s", gameID, strconv.FormatBool(ok))
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
