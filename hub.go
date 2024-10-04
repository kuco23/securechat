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

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if room, ok := h.rooms[client.roomID]; !ok {
				h.rooms[client.roomID] = newRoom(client)
			} else {
				room.add(client)
			}
		case client := <-h.unregister:
			if room, ok := h.rooms[client.roomID]; !ok {
				// danger
				close(client.send)
			} else {
				room.remove(client)
			}
		case message := <-h.broadcast:
			client := message.sender
			if room, ok := h.rooms[client.roomID]; !ok {
				// danger
				close(client.send)
			} else {
				room.broadcast(client, message)
			}
		}
	}
}
