package auth

type session struct {
	uuid         string `json:"uuid"`
	refreshToken string `json:"refreshToken"`
	accessToken  string `json:"accessToken"`
}
