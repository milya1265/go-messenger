package messenger

import (
	"WSChats/pkg/logger"
	"encoding/json"
	"github.com/gorilla/websocket"
)

const (
	SendMessageMethod   = "send_message"
	CreateChatMethod    = "new_chat"
	ReadMessageMethod   = "read_message"
	GetMessagesMethod   = "get_messages"
	DeleteMessageMethod = "delete_message"
)

type Client struct {
	Username     string `json:"username"`
	UUID         string `json:"uuid"`
	Conn         websocket.Conn
	Events       chan *[]byte
	Messages     chan *Message
	ReadMessages chan *ReadMessage
	Errors       chan *error
	service      Service
	logger       logger.Logger
}

func newClient(username, uuid string, conn *websocket.Conn, logger logger.Logger, s *Service) *Client {
	return &Client{
		Username:     username,
		UUID:         uuid,
		Conn:         *conn,
		Messages:     make(chan *Message),
		ReadMessages: make(chan *ReadMessage),
		Events:       make(chan *[]byte),
		Errors:       make(chan *error),
		logger:       logger,
		service:      *s,
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
		case readMessage := <-c.ReadMessages:

			c.Conn.WriteJSON(readMessage)
		case event := <-c.Events:
			c.logger.Debug("try to send ", event)
			c.Conn.WriteMessage(1, *event)
		case err := <-c.Errors:
			c.Conn.WriteJSON((*err).Error())
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
			c.Errors <- &err
			continue
		}

		methodJSON := mapJSON["method"]
		var methodString string

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
				c.Errors <- &err
				continue
			}

			manager.Broadcast <- message

		case ReadMessageMethod:
			c.logger.Info(ReadMessageMethod)
			readMsg, err := c.handleReadMessage(mapJSON)
			if err != nil {
				c.Errors <- &err
				continue
			}
			manager.BroadcastReadMsg <- readMsg
		case GetMessagesMethod:
			c.logger.Info(GetMessagesMethod)

			messages, err := c.handleGetChatMessages(mapJSON)
			if err != nil {
				c.Errors <- &err
				continue
			}
			response, err := json.Marshal(messages)
			if err != nil {
				c.Errors <- &err
				continue
			}

			c.Events <- &response

		case CreateChatMethod:
			c.logger.Info(CreateChatMethod)

			newChat, err := c.handleNewChat(mapJSON)
			if err != nil {
				c.Errors <- &err
				continue
			}

			response, err := json.Marshal(newChat)
			if err != nil {
				c.Errors <- &err
				continue
			}

			c.Events <- &response
		}

	}
}
