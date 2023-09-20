package DTO

type AuthUserReq struct {
	UUID string `json:"uuid"`
}

type AuthUserRes struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
