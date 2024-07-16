package helpers

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/storage"
	"github.com/umardev500/messaging-api/types"
)

var mu sync.Mutex
var clMu sync.Mutex

func BroadcastChat(msg types.Broadcast, wg *sync.WaitGroup) {
	mu.Lock()
	defer mu.Unlock()
	if wg != nil {
		defer wg.Done()
	}

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

func BroadcastChatList(ctx context.Context, msg types.BroadcastChatList, wg *sync.WaitGroup) {
	clMu.Lock()
	defer clMu.Unlock()
	if wg != nil {
		defer wg.Done()
	}

	jsonData, err := storage.Redis.JSONGet(ctx, msg.Room, "$").Result()
	if err != nil {
		if err == redis.Nil {
			log.Error().Msgf("no room data found in redis cache")
			return
		}

		log.Error().Msgf("failed to get room data from redis cache")
		return
	}

	// Unmarshal it
	var participantsData []types.RedisParticipants
	if err := json.Unmarshal([]byte(jsonData), &participantsData); err != nil {
		log.Error().Msgf("failed to unmarshal room data | err: %v", err)
		return
	}

	for userId := range participantsData[0] {
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
				delete(types.Onlines, userId)
			}
			localMu.Unlock()
		}(userId)
	}
}
