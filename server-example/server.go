package main

import (
	"net/http"
	"math/rand"
	wb "github.com/x1m3/webpusher"
	"log"
)

var wsConnHub *wb.WebsocketClientHubManager

func main() {
	wsConnHub = wb.NewWebsocketClientHubManager()

	http.HandleFunc("/register", registerController)

	http.ListenAndServe(":3000", nil)
}

func registerController(w http.ResponseWriter, r *http.Request) {

	client, err := wb.NewWebsocketClient(rand.Int63(), w, r)
	if (err != nil) {
		log.Print(err)
	}
	wsConnHub.AddClient(client)
}
