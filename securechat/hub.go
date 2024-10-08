package main

var ZERO_BYTE = []byte("0")
var ONE_BYTE = []byte("1")

type Message struct {
	data   []byte
	sender *Client
}

type Room struct {
	client1 *Client
	client2 *Client
}

type Hub struct {
	rooms      map[string]*Room
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if room, ok := h.rooms[client.roomID]; !ok {
				h.rooms[client.roomID] = newRoom(client)
			} else {
				room.Add(client)
			}
		case client := <-h.unregister:
			if room, ok := h.rooms[client.roomID]; !ok {
				// danger
			} else {
				room.Remove(client)
			}
		case message := <-h.broadcast:
			client := message.sender
			if room, ok := h.rooms[client.roomID]; !ok {
				// danger
			} else {
				room.Broadcast(client, message)
			}
		}
	}
}
