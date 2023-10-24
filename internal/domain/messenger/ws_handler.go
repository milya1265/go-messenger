package messenger

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

func (c *Client) handleNewMessage(mapJSON map[string]json.RawMessage) (*Message, error) {
	var chatID int
	var text string

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
	}
	log.Println(message)

	return message, nil
}

func (c *Client) handleNewChat(mapJSON map[string]json.RawMessage) (*NewChatRes, error) {
	c.logger.Info("Starting new chat handler")

	var chatReq NewChatReq

	creatorJSON, ok1 := mapJSON["creator"]
	titleJSON, ok2 := mapJSON["title"]
	isDirectJSON, ok3 := mapJSON["is_direct"]

	var creator string
	var title string
	var isDirect bool

	if ok1 && ok2 && ok3 {
		err := json.Unmarshal(creatorJSON, &creator)
		if err != nil {
			c.logger.Error("unmarshal creator method ", err.Error())
			return nil, err
		}

		err = json.Unmarshal(titleJSON, &title)
		if err != nil {
			c.logger.Error("unmarshal title method ", err.Error())
			return nil, err
		}

		err = json.Unmarshal(isDirectJSON, &isDirect)
		if err != nil {
			c.logger.Error("unmarshal is_direct method ", err.Error())
			return nil, err
		}
	} else {
		return nil, errors.New("bad request")
	}

	chatReq.Creator = c.UUID
	chatReq.Title = title
	chatReq.IsDirect = isDirect

	c.logger.Debug(chatReq)

	chatRes, err := c.service.NewChat(&chatReq)
	if err != nil {
		return nil, err
	}

	return chatRes, nil
}

/*func (c *Client) handleNewChat(mapJSON map[string]json.RawMessage) (*NewChatReq, error) {
	creatorJSON, ok1 := mapJSON["creator"]
	titleJSON, ok2 := mapJSON["title"]
	isDirectJSON, ok3 := mapJSON["is_direct"]
	var creator string
	var title string
	var isDirect bool
	if ok1 && ok2 && ok3 {
		err := json.Unmarshal(creatorJSON, &creator)
		if err != nil {
			c.logger.Error("unmarshal creator method ", err.Error())
			return nil, err
		}

		err = json.Unmarshal(titleJSON, &title)
		if err != nil {
			c.logger.Error("unmarshal title method ", err.Error())
			return nil, err
		}

		err = json.Unmarshal(isDirectJSON, &isDirect)
		if err != nil {
			c.logger.Error("unmarshal is_direct method ", err.Error())
			return nil, err
		}
	} else {
		return nil, errors.New("bad request")
	}

	membersJSON := mapJSON["members"]
	var members []string

	err := json.Unmarshal(membersJSON, &members)
	if err != nil {
		log.Println("unmarshal method client 177 ", err.Error())
		return nil, err
	}

	ch := NewChatReq{
		Creator:  creator,
		Title:    title,
		IsDirect: isDirect,
		Members:  members,
	}
	return &ch, nil
}
*/
