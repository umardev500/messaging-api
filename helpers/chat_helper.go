package helpers

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/storage"
	"github.com/umardev500/messaging-api/types"
)

var mu sync.Mutex

func BroadcastChat(msg types.Broadcast) {
	mu.Lock()
	defer mu.Unlock()
	theRoom := types.Rooms[msg.Room]

	for clientId, client := range theRoom {
		go func(cl *types.Client) {
			cl.Mu.Lock()
			isSender := clientId == msg.Sender
			if !isSender {
				// This will happend if antoher user is online
				// Otherwhise this block code will proccessing
				// Because we have not live users to chat
				if err := cl.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Message)); err != nil {
					log.Error().Msgf("failed to write message %v", err)
					// Handle error: remove client from the room
					delete(types.Onlines, clientId)
				}
				// @Todo handle error
				//
			}
			cl.Mu.Unlock()
		}(client)
	}
}

func BroadcastChatList(msg types.BroadcastChatList) {
	mu.Lock()
	defer mu.Unlock()

	jsonBytes, err := storage.Redis.Get(msg.Room)
	if err != nil {
		log.Error().Msgf("failed to get room data from redis cache")
		return
	}

	fmt.Println(jsonBytes)
	if len(jsonBytes) == 0 {
		log.Error().Msgf("no room data found in redis cache")
		return
	}

	var userIds []string
	if err := json.Unmarshal(jsonBytes, &userIds); err != nil {
		log.Error().Msgf("failed marshaling room data | err: %v", err)
		return
	}

	for _, userId := range userIds {
		var localMu sync.Mutex
		go func(userId string) {
			localMu.Lock()
			online, ok := types.Onlines[userId]
			if !ok {
				return
			}
			conn := online.Conn
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Error().Msgf("failed to write message %v", err)
			}
			localMu.Unlock()
		}(userId)
	}
}
