package ws

type Message struct {
	Sender   string `json:"author"`
	Receiver string `json:"receiver"`
	Id       int    `json:"id"`
	Time     int64  `json:"time"`
	Text     string `json:"text"`
}
