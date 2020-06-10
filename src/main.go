package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MadhanRaj96/chess-go/src/game"
	m "github.com/MadhanRaj96/chess-go/src/models"
	"github.com/MadhanRaj96/chess-go/src/ws"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)



func main() {
	fmt.Println("inside main")
	var wait time.Duration
	rand.Seed(time.Now().UnixNano())
	r := mux.NewRouter()
	//r.HandleFunc("/", homeHandler)
	r.HandleFunc("/", handleConnections)
	r.HandleFunc("/gameId/id:{[0-9]+}", gameRequestHandler)
	//go handleMessages()

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
