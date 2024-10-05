package main

func newRoom(client *Client) *Room {
	return &Room{client, nil}
}

func (r *Room) Add(client *Client) {
	// check if one of the room clients reconnected
	idx := 0
	if client.id == r.client1.id {
		idx = 1
	} else if client.id == r.client2.id {
		idx = 2
	}
	// handle reassignment
	if idx == 1 || r.client1 == nil && idx < 2 {
		if r.client1 != nil {
			close(r.client1.send)
		}
		r.client1 = client
		if r.client2 != nil {
			r.sendConnected()
		}
	} else if idx == 2 || r.client2 == nil {
		if r.client2 != nil {
			close(r.client2.send)
		}
		r.client2 = client
		if r.client1 != nil {
			r.sendConnected()
		}
	} else {
		// danger
		close(client.send)
	}
}

func (r *Room) Remove(client *Client) {
	if r.client1 == client {
		r.client1 = nil
		if r.client2 != nil {
			sendMessage(Message{data: ONE_BYTE}, r.client2)
		}
	} else if r.client2 == client {
		r.client2 = nil
		if r.client1 != nil {
			sendMessage(Message{data: ONE_BYTE}, r.client1)
		}
	} else {
		// danger
		sendMessage(Message{data: ONE_BYTE}, client)
	}
	close(client.send)
}

func (r *Room) Broadcast(sender *Client, message Message) {
	if r.client1 == sender {
		if r.client2 != nil {
			sendMessage(message, r.client2)
		}
	} else if r.client2 == sender {
		if r.client1 != nil {
			sendMessage(message, r.client1)
		}
	} else {
		// danger
		close(sender.send)
	}
}

func (r *Room) sendConnected() {
	sendMessage(Message{data: ZERO_BYTE}, r.client1)
	sendMessage(Message{data: ZERO_BYTE}, r.client2)
}

func sendMessage(message Message, target *Client) {
	select {
	case target.send <- message:
	default:
		close(target.send)
	}
}
