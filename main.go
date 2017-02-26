package main

import (
	"net/http"
	"math/rand"
	"log"
)

var wsConnHub *WebsocketClientHubManager

func main() {
	wsConnHub = NewWebsocketClientHubManager()

	http.HandleFunc("/register", registerController)

	http.ListenAndServe(":3000", nil)
}

func registerController(w http.ResponseWriter, r *http.Request ) {
	conn, err := wsConnHub.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	client := NewWebsocketClient(rand.Int63(), conn)
	wsConnHub.AddClient(client)
}
