package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

    "github.com/MadhanRaj96/chess-go/src/pkg/websocket"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

var linker = make(chan gameChannel)

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func getColor() string {
	r := randomInt(1, 9)
	if r%2 == 0 {
		return "white"
	}
	return "black"
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("Upgrading connection to a WS")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	/*global client storage*/
	//clients[ws] = true
	for {
		var msg GameReq
		var channel gameChannel

		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("JSON Parse error info: %#v", err)
			delete(clients, ws)
			log.Println("Closing Connection")
			break
		}
		log.Printf("Received Message from user: %s", msg.UserID)

		channel.ws = ws
		channel.resp = &GameResp{
			Message: "gameId",
			GameID:  msg.UserID + randomString(6),
			Color:   getColor(),
		}

		log.Printf("Sending Message to User: %s", msg.UserID)
		linker <- channel

	}
	defer ws.Close()
}

func handleMessages() {
	for {
		channel := <-linker

		channel.ws.WriteJSON(channel.resp)
	}
}

func main() {
	fmt.Println("inside main")
	var wait time.Duration
	rand.Seed(time.Now().UnixNano())
	r := mux.NewRouter()
	//r.HandleFunc("/", homeHandler)
	r.HandleFunc("/", handleConnections)
	go handleMessages()

	srv := &http.Server{
		Addr: "127.0.0.1:8000",

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
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
