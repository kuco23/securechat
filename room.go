package main

func sendMessage(message Message, target *Client) {
	select {
	case target.send <- message:
	default:
		close(target.send)
	}
}

func newRoom(client *Client) *Room {
	return &Room{client, nil}
}

func (r *Room) add(client *Client) {
	if r.client1 == nil || r.client1.id == client.id {
		if r.client1 != nil {
			close(r.client1.send)
		}
		r.client1 = client
		if r.client2 != nil {
			sendMessage(Message{data: ZERO_BYTE}, r.client2)
		}
	} else if r.client2 == nil || r.client2.id == client.id {
		if r.client2 != nil {
			close(r.client2.send)
		}
		r.client2 = client
		if r.client1 != nil {
			sendMessage(Message{data: ZERO_BYTE}, r.client1)
		}
	} else {
		// danger
		sendMessage(Message{data: ONE_BYTE}, r.client1)
		sendMessage(Message{data: ONE_BYTE}, r.client2)
		close(client.send)
	}
}

func (r *Room) remove(client *Client) {
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

func (r *Room) broadcast(sender *Client, message Message) {
	if r.client1 == sender {
		sendMessage(message, r.client2)
	} else if r.client2 == sender {
		sendMessage(message, r.client1)
	} else {
		// danger
		close(sender.send)
	}
}
