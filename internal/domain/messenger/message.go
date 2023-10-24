package messenger

type Message struct {
	Id     int    `json:"id"`
	ChatID int    `json:"chat"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	Time   int64  `json:"time"`
	Reply  int    `json:"Reply"`
}

type NewMessageReq struct {
	ChatID string `json:"chat"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	Reply  int    `json:"Reply"`
}

type NewMessageRes struct {
	Id     int    `json:"id"`
	ChatID string `json:"chat"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	Time   int64  `json:"time"`
	Reply  int    `json:"Reply"`
}
