package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/storage"
	"github.com/umardev500/messaging-api/types"
)

type chatHandler struct {
	chatService domain.ChatService
}

func NewChatHandler(chatService domain.ChatService) domain.ChatHandler {
	return &chatHandler{
		chatService: chatService,
	}
}

type Online struct {
	Conn *websocket.Conn
}

var onlines = make(map[string]*Online)
var rooms = make(map[string]map[string]*types.Client)
var listMu, chatMu sync.Mutex

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

			broadcastData := types.Broadcast{
				Sender:  userId,
				Room:    room,
				Clients: rooms[room],
				Message: string(msg),
			}
			go ch.broadcastMessage(broadcastData)
		}
	})
}

func (ch *chatHandler) broadcastMessage(msg types.Broadcast) {
	chatMu.Lock()
	defer chatMu.Unlock()

	log.Info().Msgf("broadcasting message: %s", msg.Message)
	sender := msg.Clients[msg.Sender]

	// @Todo store message to the database
	//
	//
	msgId := uuid.New().String()
	_ = types.InputNewMessage{
		Id:     msgId,
		Room:   &msg.Room,
		UserId: msg.Sender,
		Text:   &msg.Message,
	}

	// Broadcasting... we can deal what we want?
	// Does we need input to the database first or send websocket data first
	// If want to make input to the database first just make line of code is sync and make pushing websocket logic is after input data
	for _, client := range msg.Clients {
		// Broadcast message with goroutine
		go func() {
			client.Mu.Lock()
			isSender := client == sender
			if !isSender {
				// This will happend if antoher user is online
				// Otherwhise this block code will proccessing
				// Because we have not live users to chat
				client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Message))
				// @Todo handle error
				//
			}
			client.Mu.Unlock()
		}()
	}

	// @Todo push incoming chat to the chat list
	// We get participant data from redis cache
	// And then we push new chat to the chat list
	// First we matching does list of participants is exist in the onlines variable
	// If yes then we push it otherwise skip it
	// And skip to push to sender chat list
	b, err := storage.Redis.Get(msg.Room)
	fmt.Println(b)
	fmt.Println(err)
}

func (ch *chatHandler) PushNewChat(c *fiber.Ctx) error {
	var payload types.PushNewChatPayload
	if err := c.BodyParser(&payload); err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userId = "1"
	payload.UserId = userId
	ch.chatService.PushNewChat(ctx, payload)

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
