package messenger

import (
	"encoding/json"
	"log"
	"time"
)

func (c *Client) handleNewMessage(mapJSON map[string]json.RawMessage) (*Message, error) {
	var chatID int
	var text string
	var reply int

	if _, ok := mapJSON["reply"]; ok {
		err := json.Unmarshal(mapJSON["reply"], &reply)
		if err != nil {
			c.logger.Error("unmarshal method ", err.Error()) //ОТЛОВИТЬ ОШИБКУ
			return nil, err
		}
	}

	err := json.Unmarshal(mapJSON["chat_id"], &chatID)
	if err != nil {
		c.logger.Error("unmarshal method ", err.Error()) //ОТЛОВИТЬ ОШИБКУ
		return nil, err
	}
	err = json.Unmarshal(mapJSON["text"], &text)
	if err != nil {
		c.logger.Error("unmarshal method ", err.Error())
		return nil, err
	}

	message := &Message{
		Sender: c.UUID,
		ChatID: chatID,
		Id:     0,
		Time:   time.Now().Unix(),
		Text:   text,
		Reply:  reply,
	}
	log.Println(message)

	return message, nil
}

func (c *Client) handleNewChat(mapJSON map[string]json.RawMessage) (*NewChatRes, error) {
	c.logger.Info("Starting new chat handler")

	var chatReq NewChatReq

	//var creator string
	var title string
	var isDirect bool
	var members []string

	//creatorJSON, ok := mapJSON["creator"]
	//if ok {
	//	err := json.Unmarshal(creatorJSON, &creator)
	//	if err != nil {
	//		c.logger.Error("unmarshal creator method ", err.Error())
	//		return nil, err
	//	}
	//}

	membersJSON, ok := mapJSON["members"]
	if ok {
		err := json.Unmarshal(membersJSON, &members)
		if err != nil {
			c.logger.Error("unmarshal title method ", err.Error())
			return nil, err
		}
	}

	titleJSON, ok := mapJSON["title"]
	if ok {
		err := json.Unmarshal(titleJSON, &title)
		if err != nil {
			c.logger.Error("unmarshal title method ", err.Error())
			return nil, err
		}
	}
	isDirectJSON, ok := mapJSON["is_direct"]
	if ok {
		err := json.Unmarshal(isDirectJSON, &isDirect)
		if err != nil {
			c.logger.Error("unmarshal is_direct method ", err.Error())
			return nil, err
		}
	}

	chatReq.Creator = c.UUID
	chatReq.Title = title
	chatReq.IsDirect = isDirect
	chatReq.Members = members

	chatRes, err := c.service.NewChat(&chatReq)
	if err != nil {
		return nil, err
	}

	c.logger.Debug(chatRes)

	return chatRes, nil
}

func (c *Client) handleReadMessage(mapJSON map[string]json.RawMessage) (*ReadMessage, error) {
	var chatID int
	var userID = c.UUID
	var msg int

	msgJSON, ok := mapJSON["last_read_msg"]
	if ok {
		err := json.Unmarshal(msgJSON, &msg)
		if err != nil {
			c.logger.Error("unmarshal last_read_msg method ", err.Error())
			return nil, err
		}
	}
	chatIdJSON, ok := mapJSON["chat_id"]
	if ok {
		err := json.Unmarshal(chatIdJSON, &chatID)
		if err != nil {
			c.logger.Error("unmarshal chat_id method ", err.Error())
			return nil, err
		}
	}

	readMessage := &ReadMessage{
		UserID:      userID,
		ChatID:      chatID,
		LastReadMsg: msg,
	}

	err := c.service.SaveReadStatus(readMessage)
	if err != nil {
		return nil, err
	}

	return readMessage, nil
}

func (c *Client) handleGetChatMessages(mapJSON map[string]json.RawMessage) ([]*Message, error) {
	var limit int
	var offset int
	var chatID int

	limitJSON, ok := mapJSON["limit"]
	if ok {
		err := json.Unmarshal(limitJSON, &limit)
		if err != nil {
			c.logger.Error("unmarshal limit method ", err.Error())
			return nil, err
		}
	}
	offsetJSON, ok := mapJSON["offset"]
	if ok {
		err := json.Unmarshal(offsetJSON, &offset)
		if err != nil {
			c.logger.Error("unmarshal offset method ", err.Error())
			return nil, err
		}
	}
	chatIdJSON, ok := mapJSON["chat_id"]
	if ok {
		err := json.Unmarshal(chatIdJSON, &chatID)
		if err != nil {
			c.logger.Error("unmarshal chat_id method ", err.Error())
			return nil, err
		}
	}

	userID := c.UUID

	return c.service.GetChatMessages(chatID, limit, offset, userID)
}

func (c *Client) handleDeleteMessage(mapJSON map[string]json.RawMessage) error {
	var messageID int

	messageJSON, ok := mapJSON["message"]
	if ok {
		err := json.Unmarshal(messageJSON, &messageID)
		if err != nil {
			c.logger.Error("unmarshal message method ", err.Error())
			return err
		}
	}

	err := c.service.DeleteMessage(messageID, c.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) handleEditMessage(mapJSON map[string]json.RawMessage) error {
	var messageID int
	var text string

	messageJSON, ok := mapJSON["message"]
	if ok {
		err := json.Unmarshal(messageJSON, &messageID)
		if err != nil {
			c.logger.Error("unmarshal message method ", err.Error())
			return err
		}
	}

	textJSON, ok := mapJSON["text"]
	if ok {
		err := json.Unmarshal(textJSON, &text)
		if err != nil {
			c.logger.Error("unmarshal message method ", err.Error())
			return err
		}
	}

	err := c.service.EditMessage(messageID, text, c.UUID)
	if err != nil {
		return err
	}

	return nil
}
