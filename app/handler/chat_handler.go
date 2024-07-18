package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
	"github.com/umardev500/messaging-api/utils"
)

type chatHandler struct {
	chatService domain.ChatService
}

func NewChatHandler(chatService domain.ChatService) domain.ChatHandler {
	return &chatHandler{
		chatService: chatService,
	}
}

var listMu, chatMu sync.Mutex

func (ch *chatHandler) GetChatList(c *fiber.Ctx) error {
	var date = c.Query("date")
	var ticket = uuid.New().String()

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, types.ProcIdKey, ticket)

	handler, err := utils.ValidateDateResp(ctx, date, c)
	if err != nil {
		log.Error().Msgf("date validation failed | err: %v | ticket: %s", err, ticket)
		return handler
	}

	userId, err := utils.GetUserIdFromLocals(ctx)
	if err != nil {
		log.Error().Msgf("failed to get user id from locals | err: %v | ticket: %s", err, ticket)
		return c.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Code:    fiber.StatusInternalServerError,
			Message: fiber.ErrInternalServerError.Message,
		})
	}

	resp := ch.chatService.GetChatList(ctx, types.GetChatListParam{
		UserId: userId,
		Date:   date,
	})

	return c.Status(resp.Code).JSON(resp)
}

// WsChatList is websocket handler to get realtime chat list
func (ch *chatHandler) WsChatList() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		tokenString := c.Query("token")
		online := &types.Online{
			Conn: c,
		}

		resp, err := ch.chatService.GetClaims(tokenString)
		if err != nil {
			if err := c.Conn.WriteJSON(resp); err != nil {
				log.Error().Msgf("failed to write ws message | err: %v | ticket: %s", err, resp.Ticket)
			}

			return
		}
		userData := resp.Data.(jwt.MapClaims)["user"].(map[string]interface{})
		userId := userData["id"].(string)

		// Add user to online group
		listMu.Lock()
		if _, ok := types.Onlines[userId]; !ok {
			types.Onlines[userId] = &types.Online{}
		}
		types.Onlines[userId] = online
		listMu.Unlock()

		// Remove disconnected user
		defer func() {
			listMu.Lock()
			delete(types.Onlines, userId)
			listMu.Unlock()
			c.Close()
		}()

		// @Todo
		// We do fetch chat list of connected user
		//
		fmt.Println(types.Onlines)

		// Listen for incoming message
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Error().Msgf("failed to read message %v", err)
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
		tokenString := c.Query("token")
		client := &types.Client{
			Conn: c,
		}

		resp, err := ch.chatService.GetClaims(tokenString)
		if err != nil {
			if err := c.Conn.WriteJSON(resp); err != nil {
				log.Error().Msgf("failed to write ws message | err: %v | ticket: %s", err, resp.Ticket)
			}

			return
		}
		userData := resp.Data.(jwt.MapClaims)["user"].(map[string]interface{})
		userId := userData["id"].(string)

		// Appends user to types.Rooms
		chatMu.Lock()
		if _, ok := types.Rooms[room]; !ok {
			types.Rooms[room] = make(map[string]*types.Client)
		}
		types.Rooms[room][userId] = client
		chatMu.Unlock()

		// Remove user from types.Rooms
		defer func() {
			chatMu.Lock()
			delete(types.Rooms[room], userId)
			chatMu.Unlock()
			c.Close()
		}()

		// Listen for message
		for {
			_, _, err := c.Conn.ReadMessage()
			if err != nil {
				log.Error().Err(err).Msg("error reading message")
				return
			}
		}
	})
}

func (ch *chatHandler) CreateNewChat(c *fiber.Ctx) error {
	// Parsing
	var payload types.CreateNewChatPayload
	if err := c.BodyParser(&payload); err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Proccessing intialization of new chat
	var userId = "78901234-5678-9012-3456-789012345678"
	payload.UserId = userId
	resp, err := ch.chatService.CreateNewChat(ctx, payload)
	if err != nil {
		return c.Status(resp.Code).JSON(resp)
	}
	payload.Participants = append(payload.Participants, userId)

	// Push new chat list to the user online matched
	for _, participant := range payload.Participants {
		var localMu sync.Mutex
		go func(participant string, resp types.Response) {
			localMu.Lock()
			if online, ok := types.Onlines[participant]; ok {
				log.Debug().Msg("Participant is online ready to push new chat list")

				err := online.Conn.WriteJSON(resp.Data)
				if err != nil {
					log.Error().Msgf("Failed to push to user chat list id: %s", participant)
					return
				}
			}
			localMu.Unlock()
		}(participant, resp)
	}

	return c.Status(resp.Code).JSON(resp)
}
