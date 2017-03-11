package webpusher

import (
	"github.com/gorilla/websocket"
	"errors"
	"sync"
	"time"
	"log"
	"net/http"
)

type WebsocketClient struct {
	id   int64
	conn *websocket.Conn
	sync.Mutex
}

func NewWebsocketClient(id int64, w http.ResponseWriter, r *http.Request) (*WebsocketClient, error) {

	wsUpgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if (err != nil) {
		return nil, err
	}

	client := new(WebsocketClient)
	client.conn = conn
	client.id = id

	client.conn.SetReadLimit(512)
	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil
	})
	go client.readLoop()
	go client.sendPings(10 * time.Second)
	return client, nil
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

func (this *WebsocketClient) sendPing() error {
	this.Lock(); defer this.Unlock()
	this.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return this.conn.WriteMessage(websocket.PingMessage, []byte{});
}

func (this *WebsocketClient) sendPings(d time.Duration) {
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
