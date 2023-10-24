package messenger

import (
	"WSChats/pkg/logger"
	"encoding/json"
	"github.com/gorilla/websocket"
)

const (
	SendMessageMethod = "send_message"
	CreateChatMethod  = "new_chat"
	ReadMessageMethod = "read_message"
	GetMessagesMethod = "get_messages"
)

type Client struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Conn     websocket.Conn
	Messages chan *Message
	Errors   chan *error
	service  Service
	logger   logger.Logger
}

func newClient(username, uuid string, conn *websocket.Conn, logger logger.Logger, s *Service) *Client {
	return &Client{
		Username: username,
		UUID:     uuid,
		Conn:     *conn,
		Messages: make(chan *Message),
		//ReadMessages: make(chan *),
		logger:  logger,
		service: *s,
	}
}

func (c *Client) SendRes() {

	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message := <-c.Messages:
			c.Conn.WriteJSON(message)
		case err := <-c.Errors:
			c.Conn.WriteJSON(err)
		}
	}
}

func (c *Client) HandleReq(manager *Manager) {
	defer func() {
		c.Conn.Close()

		manager.Offline <- c
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("websocket conn ", err.Error())
			}
			break
		}

		mapJSON := make(map[string]json.RawMessage)
		err = json.Unmarshal(p, &mapJSON)
		if err != nil {
			c.logger.Error("unmarshal method ", err.Error())
			break
		}

		methodJSON := mapJSON["method"]
		var methodString string
		c.logger.Info("map method", string(methodJSON))

		err = json.Unmarshal(methodJSON, &methodString)
		if err != nil {
			c.logger.Error("unmarshal method ", err.Error())
			continue
		}

		switch methodString {
		case SendMessageMethod:
			c.logger.Info(SendMessageMethod)

			message, err := c.handleNewMessage(mapJSON)
			if err != nil {
				c.logger.Error(err.Error())
			}

			manager.Broadcast <- message

		case ReadMessageMethod:
			c.logger.Info(ReadMessageMethod)

		case GetMessagesMethod:
			c.logger.Info(GetMessagesMethod)

		//case CreateChatMethod:
		//	c.logger.Info(CreateChatMethod)
		//
		//	new c.handleNewChat(mapJSON)
		default:
			c.logger.Info("error", err)
		}

	}
}
