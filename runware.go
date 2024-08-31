package runware

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn             *websocket.Conn
	incomingMessages chan []byte
	reconnect        chan bool
	ApiKey           string

	mu sync.Mutex
}

type socketMessage struct {
	Data   []map[string]interface{} `json:"data"`
	Errors []map[string]interface{} `json:"errors"`
}

func New(ctx context.Context, apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey is required")
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, "wss://ws-api.runware.ai/v1", nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		ApiKey:           apiKey,
		Conn:             conn,
		incomingMessages: make(chan []byte),
		reconnect:        make(chan bool),
	}

	go client.readHandler()
	go client.reconnectHandler()
	go client.heartbeatHandler()

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}

func (c *Client) Send(msg []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *Client) checkError(msg socketMessage) error {
	if len(msg.Errors) < 1 {
		return nil
	}

	if isErr, ok := msg.Errors[0]["error"]; ok && isErr.(bool) {
		err := msg.Errors[0]["errorMessage"]
		return fmt.Errorf("runware error: %v", err)
	}

	return nil
}

func (c *Client) Connect() error {
	sendData, err := json.Marshal([]*ConnectRequestData{NewConnectRequestData(c.ApiKey)})
	if err != nil {
		return err
	}

	if err := c.Send(sendData); err != nil {
		return err
	}

	for msg := range c.incomingMessages {
		var msgData socketMessage
		if err := json.Unmarshal(msg, &msgData); err != nil {
			return err
		}

		if err := c.checkError(msgData); err != nil {
			return err
		}

		if msgData.Data[0]["taskType"].(string) != "authentication" {
			c.incomingMessages <- msg
			continue
		}

		break
	}
	return nil
}
