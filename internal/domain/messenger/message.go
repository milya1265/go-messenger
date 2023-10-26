package messenger

type Message struct {
	Id     int    `json:"id"`
	ChatID int    `json:"chat"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	Time   int64  `json:"time"`
	Reply  int    `json:"Reply"`
}

type ReadMessage struct {
	UserID      string `json:"user_id"`
	ChatID      int    `json:"chat_id"`
	LastReadMsg int    `json:"last_read_msg"`
}
