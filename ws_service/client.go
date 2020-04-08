package ws_service

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

type Client struct {
	UserId    string
	Timestamp int64
	conn      *websocket.Conn
}

func (c *Client) sendMessage(userId, message string) {
	if c.conn != nil {
		c.conn.Write([]byte(message))
	}
}

func (c *Client) heartbeat() error {
	millis := time.Now().UnixNano() / 1e6
	heartbeat := struct {
		Heartbeat int64 `json:heartbeat`
	}{millis}
	bytes, _ := json.Marshal(heartbeat)
	_, err := c.conn.Write(bytes)
	return err
}

func (c *Client) Listen() {
	for range time.Tick(ListenInterval * time.Second) {
		err := c.heartbeat()
		if err != nil {
			log.Println("heartbeat error")
			return
		}
	}
}
