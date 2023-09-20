package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type Client struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Conn     websocket.Conn
	Messages chan *Message
}

func (c *Client) SendWS() {

	defer func() {
		logrus.Error("КОНЕЦ send")
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Messages
		if ok == false {
			return
		}
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ListenWS(chat *Messenger) {
	defer func() {
		logrus.Error("КОНЕЦ listen")
		c.Conn.Close()

		chat.Offline <- c
	}()

	for {
		message := &Message{
			Sender:   c.UUID,
			Receiver: "",
			Id:       0,
			Time:     time.Now().Unix(),
			Text:     "",
		}

		err := c.Conn.ReadJSON(message)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Error("websocket conn", err.Error())
			}
			break
		}

		log.Println("message after ", message)

		chat.Broadcast <- message
	}
}
