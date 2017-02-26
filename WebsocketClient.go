package main

import (
	"github.com/gorilla/websocket"
	"errors"
	"sync"
	"time"
	"log"
)

type WebsocketClient struct {
	id   int64
	conn *websocket.Conn
	sync.Mutex
}

func NewWebsocketClient(id int64, conn *websocket.Conn) *WebsocketClient {
	client := new(WebsocketClient)
	client.id = id
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil
	})
	go client.readLoop()
	go client.sendPings(10 * time.Second)
	client.conn = conn
	return client
}

func (this *WebsocketClient) Send(message []byte) error {
	this.Lock(); defer this.Unlock()
	err := this.conn.WriteMessage(websocket.TextMessage, message)
	if (err != nil) {
		return errors.New("Connection lost")
	}
	return nil
}

func (this *WebsocketClient) Id() int64 {
	return this.id
}

func (this * WebsocketClient) sendPing() error {
	this.Lock(); defer this.Unlock()
	this.conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	return this.conn.WriteMessage(websocket.PingMessage, []byte{});
}

func (this * WebsocketClient) sendPings(d time.Duration) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	defer this.conn.Close()
	for {
		select {
		case <-ticker.C:
			if err := this.sendPing(); err != nil {
				log.Print(err)
				return
			}
		}
	}
}

func (this *WebsocketClient) readLoop() {
	for {
		if _, _, err := this.conn.NextReader(); err != nil {
			this.conn.Close()
			return
		}
	}
}
