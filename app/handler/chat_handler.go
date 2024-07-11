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

func init() {
	// Broadcast messages to all clients in rooms
	go func() {
		for {
			msg := <-boradcast
			chatMu.Lock()
			log.Info().Msgf("broadcasting message: %s", msg.Message)

			// msg.Clients is list of clients in room
			// so we need to send message to all clients except sender
			for userId, client := range msg.Clients {

				isSender := client == msg.Sender
				if !isSender {
					fmt.Println(msg.Room)
					if online, ok := onlines[userId]; ok {
						fmt.Println("user is online and ready to receive message")
						// Push new message to chat list of online user
						pushData := map[string]string{
							"message": msg.Message,
							"room":    msg.Room,
						}
						online.Conn.WriteJSON(pushData)
					}
					client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Message))
				}
			}

			chatMu.Unlock()
		}
	}()
}

func (ch *chatHandler) PushNewChat(c *fiber.Ctx) error {
	onlines["1"].Conn.WriteMessage(websocket.TextMessage, []byte("hi from push"))

	return c.JSON(len(onlines))
}
