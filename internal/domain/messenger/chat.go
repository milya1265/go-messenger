package messenger

type Chat struct {
	Id       int      `json:"id"`
	Creator  string   `json:"creator"`
	Title    string   `json:"title"`
	Members  []string `json:"members"`
	IsDirect bool     `json:"is_direct"`
}

type NewChatReq struct {
	Title    string   `json:"title"`
	Creator  string   `json:"creator"`
	Members  []string `json:"members"`
	IsDirect bool     `json:"is_direct"`
}

type NewChatRes struct {
	Id       int      `json:"id"`
	Title    string   `json:"title"`
	Creator  string   `json:"creator"`
	Members  []string `json:"members"`
	IsDirect bool     `json:"is_direct"`
}
