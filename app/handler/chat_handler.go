package handler

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type chatHandler struct{}

func NewChatHandler() domain.ChatHandler {
	return &chatHandler{}
}

type Online struct {
	Conn *websocket.Conn
}

var onlines = make(map[string]*Online)
var rooms = make(map[string]map[string]*types.Client)
var listMu, chatMu sync.Mutex
var boradcast = make(chan types.Broadcast)

// WsChatList is websocket handler to get realtime chat list
func (ch *chatHandler) WsChatList() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		userId := c.Query("userid")
		online := &Online{
			Conn: c,
		}

		// Add user to online group
		listMu.Lock()
		if _, ok := onlines[userId]; !ok {
			onlines[userId] = &Online{}
		}
		onlines[userId] = online
		listMu.Unlock()

		// Remove disconnected user
		defer func() {
			listMu.Lock()
			delete(onlines, userId)
			listMu.Unlock()
			c.Close()
		}()
		fmt.Println(onlines)

		// Listen for incoming message
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(msg))
		}
	})
}

// WsChat is websocket handler for live chatting
func (ch *chatHandler) WsChat() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		room := c.Params("room")
		userId := c.Query("userid")
		client := &types.Client{
			Conn: c,
		}

		// Appends user to rooms
		chatMu.Lock()
		if _, ok := rooms[room]; !ok {
			rooms[room] = make(map[string]*types.Client)
		}
		rooms[room][userId] = client
		chatMu.Unlock()

		// Remove user from rooms
		defer func() {
			chatMu.Lock()
			delete(rooms[room], userId)
			chatMu.Unlock()
			c.Close()
		}()

		// Listen for message
		for {
			_, msg, err := c.Conn.ReadMessage()
			if err != nil {
				log.Error().Err(err).Msg("error reading message")
				return
			}

			boradcast <- types.Broadcast{
				Sender:  client,
				Room:    room,
				Clients: rooms[room],
				Message: string(msg),
			}
		}
	})
}

// This is method for listening broadcast message
func init() {
	// Broadcast messages to all clients in rooms
	go func() {
		for {
			msg := <-boradcast
			chatMu.Lock()
			log.Info().Msgf("broadcasting message: %s", msg.Message)

			// msg.Clients is list of clients in room
			// so we need to send message to all clients except sender
			for _, client := range msg.Clients {
				client.Mu.Lock()
				isSender := client == msg.Sender
				if !isSender {
					// This will happend if antoher user is online
					// Otherwhise this block code will proccessing
					// Because we have not live users to chat
					client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Message))
				}
				client.Mu.Unlock()
			}

			// Do get all clients on room participants from redis cache
			// And then do matching that user ids to onlines users variable
			// If user online found just push that new message to the room or chat list
			// As higlight message or with counter message
			// @Todo

			// Check for user is online
			// if online, ok := onlines[userId]; ok {
			// 	fmt.Println("user is online and ready to receive message")
			// 	// Push new message to chat list of online user
			// 	pushData := map[string]string{
			// 		"message": msg.Message,
			// 		"room":    msg.Room,
			// 	}
			// 	online.Conn.WriteJSON(pushData)
			// }

			chatMu.Unlock()
		}
	}()
}

func (ch *chatHandler) PushNewChat(c *fiber.Ctx) error {
	var payload types.PushNewChatPayload
	if err := c.BodyParser(&payload); err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var totalOnline = len(onlines)

	data := map[string]interface{}{
		"online":  totalOnline,
		"payload": payload,
	}

	// Push new chat list to the user online matched
	for _, participant := range payload.Participants {
		if online, ok := onlines[participant]; ok {
			log.Debug().Msg("Participant is online ready to push new chat list")

			pushData := map[string]string{
				"message": *payload.Message.Text,
				"room":    "1000",
			}
			err := online.Conn.WriteJSON(pushData)
			if err != nil {
				return c.JSON("failed to push new chat list to the users")
			}
		}
	}

	// @Todo database storing

	return c.JSON(data)
}
