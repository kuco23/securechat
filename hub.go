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

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// rooms
	rooms map[string]*Room
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) sendMessage(message Message, target *Client) {
	select {
	case target.send <- message:
	default:
		close(target.send)
		delete(h.clients, target)
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			room := h.rooms[client.roomID]
			target := room.client2
			if room.client1 == nil || client.id == room.client1.id {
				room.client1 = client
			} else {
				room.client2 = client
				target = room.client1
			}
			if target != nil {
				h.sendMessage(Message{data: ZERO_BYTE}, target)
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			room := h.rooms[message.sender.roomID]
			target := message.sender
			if room.client1 == target {
				target = room.client2
			}
			if target == nil {
				delete(h.clients, message.sender)
				close(message.sender.send)
			} else {
				h.sendMessage(message, target)
			}
		}
	}
}
