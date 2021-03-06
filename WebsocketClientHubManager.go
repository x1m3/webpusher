package webpusher

import (
	"sync"
	"log"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebsocketClientHubManager struct {
	sync.RWMutex
	connectionList map[int64]*WebsocketClient
}

func NewWebsocketClientHubManager() *WebsocketClientHubManager {
	hub := new(WebsocketClientHubManager)
	hub.connectionList = make(map[int64]*WebsocketClient)
	return hub
}

func (this *WebsocketClientHubManager) AddClient(client *WebsocketClient) {
	this.Lock()
	this.connectionList[client.Id()] = client
	this.Unlock()
}

func (this *WebsocketClientHubManager) RemoveClient(clientId int64) {
	this.Lock()
	delete(this.connectionList, clientId)
	this.Unlock()
}

func (this *WebsocketClientHubManager) Send(clientId int64, message []byte) {
	this.RLock()
	client := this.connectionList[clientId]
	err := client.Send(message)
	if (err != nil) {
		log.Print(err)
		go this.RemoveClient(clientId)
	}
	this.RUnlock()
}

func (this *WebsocketClientHubManager) Broadcast(message []byte) {
	this.RLock()
	for id, _ := range this.connectionList {
		go this.Send(id, message)
	}
	this.RUnlock()
}