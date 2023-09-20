package DTO

type CreateUserReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRes struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type GetUserByEmailReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserByEmailRes struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
