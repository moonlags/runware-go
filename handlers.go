package runware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) readHandler() {
	defer c.Close()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			c.reconnect <- true
			time.Sleep(10 * time.Second)
			continue
		}

		var msgData socketMessage
		_ = json.Unmarshal(msg, &msgData)

		if c.checkError(msgData) != nil {
			c.incomingMessages <- msg
			continue
		}

		if _, ok := msgData.Data[0]["pong"]; ok {
			continue
		}

		c.incomingMessages <- msg
	}
}

func (c *Client) heartbeatHandler() {
	ticker := time.NewTicker(100 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		if err := c.Send([]byte(`{"taskType":"ping" ,"ping": true}`)); err != nil {
			c.reconnect <- true
			time.Sleep(10 * time.Second)
		}
	}
}

func (c *Client) reconnectHandler() {
	for {
		<-c.reconnect

		log.Println("reconnecting to runware")
		c.Close()

		var err error
		c.Conn, _, err = websocket.DefaultDialer.Dial("wss://ws-api.runware.ai/v1", nil)
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		if c.Connect() != nil {
			time.Sleep(10 * time.Second)
			c.reconnect <- true
		}
	}
}
