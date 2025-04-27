package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	conn    *websocket.Conn
	send    chan Message
	hub     Hub
	userID  uint
	chatIDs map[uint]bool
	logger  *zap.SugaredLogger
}

func NewClient(conn *websocket.Conn, send chan Message, hub Hub, userID uint, chatIDs map[uint]bool, logger *zap.SugaredLogger) *Client {
	return &Client{conn: conn, send: send, hub: hub, userID: userID, chatIDs: chatIDs, logger: logger}
}

type IncomingMessage struct {
	ChatID  uint   `json:"chat_id"`
	Content string `json:"content"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.logger.Infow("WebSocket connection closed", "clientID", c.userID)
				return
			}
			c.logger.Errorw("Error reading message", "error", err, "clientID", c.userID)
			return
		}

		var msgStruct IncomingMessage
		if err := json.Unmarshal(msg, &msgStruct); err != nil {
			c.logger.Errorw("Error unmarshalling message", "error", err)
			return
		}

		c.hub.BroadcastMessage(Message{
			IncomingMessage: msgStruct,
			SenderID:        c.userID,
		})
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(time.Second * 60)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		c.logger.Infow("WebSocket connection closed", "clientID", c.userID)
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.logger.Infow("Channel closed, sending close message", "clientID", c.userID)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				c.logger.Errorw("Error sending message", "error", err, "clientID", c.userID)
				return
			}
			c.logger.Infow("Message sent", "senderID", message.SenderID, "chatID", message.ChatID, "content", message.Content)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Errorw("Error sending ping", "error", err, "clientID", c.userID)
				return
			}
			c.logger.Infow("Ping sent", "clientID", c.userID)
		}
	}
}
