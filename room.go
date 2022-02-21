package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	//他のクライアントに転送するメッセージを保持するチャネル
	forward chan []byte
	//チャットルームに参加希望のクライアントのチャネル
	join chan *client
	//チャットルームから退出希望のクライアントのチャネル
	leave chan *client
	//clientに在籍している全てのクライアントを保持
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <- r.join:
			r.clients[client] = true
		case client := <- r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <- r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
					//messageの送信
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
	}
	client := &client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() {r.leave <- client} ()
	go client.write()
	client.read()
}