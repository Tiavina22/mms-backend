package websocket

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients (mapped by user ID)
	clients map[uuid.UUID]*Client

	// Inbound messages from the clients
	broadcast chan []byte

	// Direct messages to specific users
	directMessage chan *Message

	// Group messages
	groupMessage chan *Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Group memberships (groupID -> []userID)
	groups map[uuid.UUID][]uuid.UUID
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:     make(chan []byte),
		directMessage: make(chan *Message),
		groupMessage:  make(chan *Message),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make(map[uuid.UUID]*Client),
		groups:        make(map[uuid.UUID][]uuid.UUID),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.UserID] = client
			log.Printf("Client connected: %s (UserID: %s)", client.Username, client.UserID)

			// Send user_joined event to all clients
			joinedMsg := Message{
				Type: "user_joined",
				Data: map[string]interface{}{
					"user_id":  client.UserID,
					"username": client.Username,
				},
			}
			if data, err := json.Marshal(joinedMsg); err == nil {
				h.BroadcastToAll(data)
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				log.Printf("Client disconnected: %s (UserID: %s)", client.Username, client.UserID)

				// Send user_left event to all clients
				leftMsg := Message{
					Type: "user_left",
					Data: map[string]interface{}{
						"user_id":  client.UserID,
						"username": client.Username,
					},
				}
				if data, err := json.Marshal(leftMsg); err == nil {
					h.BroadcastToAll(data)
				}
			}

		case message := <-h.broadcast:
			// Broadcast to all connected clients
			h.BroadcastToAll(message)

		case message := <-h.directMessage:
			// Send message to specific user
			if message.ReceiverID != uuid.Nil {
				h.SendToUser(message.ReceiverID, message)
			}

		case message := <-h.groupMessage:
			// Broadcast to all group members
			if message.GroupID != uuid.Nil {
				h.BroadcastToGroup(message.GroupID, message)
			}
		}
	}
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(message []byte) {
	for _, client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client.UserID)
		}
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID uuid.UUID, message *Message) {
	client, ok := h.clients[userID]
	if !ok {
		log.Printf("User %s is not connected", userID)
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case client.Send <- data:
	default:
		close(client.Send)
		delete(h.clients, userID)
	}
}

// BroadcastToGroup sends a message to all members of a group
func (h *Hub) BroadcastToGroup(groupID uuid.UUID, message *Message) {
	// Get group members
	members, ok := h.groups[groupID]
	if !ok {
		log.Printf("Group %s not found in hub", groupID)
		// In a real implementation, you would fetch members from the database
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	// Send to all group members
	for _, memberID := range members {
		if client, ok := h.clients[memberID]; ok {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, memberID)
			}
		}
	}
}

// AddUserToGroup adds a user to a group for message routing
func (h *Hub) AddUserToGroup(groupID, userID uuid.UUID) {
	if members, ok := h.groups[groupID]; ok {
		// Check if user is already in the group
		for _, id := range members {
			if id == userID {
				return
			}
		}
		h.groups[groupID] = append(members, userID)
	} else {
		h.groups[groupID] = []uuid.UUID{userID}
	}
}

// RemoveUserFromGroup removes a user from a group
func (h *Hub) RemoveUserFromGroup(groupID, userID uuid.UUID) {
	if members, ok := h.groups[groupID]; ok {
		newMembers := make([]uuid.UUID, 0)
		for _, id := range members {
			if id != userID {
				newMembers = append(newMembers, id)
			}
		}
		if len(newMembers) > 0 {
			h.groups[groupID] = newMembers
		} else {
			delete(h.groups, groupID)
		}
	}
}

// GetOnlineUsers returns a list of online user IDs
func (h *Hub) GetOnlineUsers() []uuid.UUID {
	users := make([]uuid.UUID, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline checks if a user is online
func (h *Hub) IsUserOnline(userID uuid.UUID) bool {
	_, ok := h.clients[userID]
	return ok
}

