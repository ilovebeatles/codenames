package hub

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"nhooyr.io/websocket"
)

type Client struct {
	conn      *websocket.Conn
	hub       *Hub
	roomID    string
	sessionID string
	playerID  string
	send      chan []byte
}

func NewClient(conn *websocket.Conn, hub *Hub, roomID, sessionID, playerID string) *Client {
	return &Client{
		conn:      conn,
		hub:       hub,
		roomID:    roomID,
		sessionID: sessionID,
		playerID:  playerID,
		send:      make(chan []byte, 64),
	}
}

func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) != -1 {
				log.Printf("ws closed: %v", err)
			}
			return
		}

		var msg IncomingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			c.SendError("invalid message format")
			continue
		}

		c.hub.HandleMessage(ctx, c, msg)
	}
}

func (c *Client) WritePump(ctx context.Context) {
	defer c.conn.Close(websocket.StatusNormalClosure, "")

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) Send(msg OutgoingMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	select {
	case c.send <- data:
	default:
		log.Printf("client send buffer full, dropping message")
	}
}

func (c *Client) SendError(errMsg string) {
	c.Send(OutgoingMessage{Type: MsgError, Error: errMsg})
}
