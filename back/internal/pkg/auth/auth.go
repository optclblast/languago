package auth

type Token struct {
	Token        string `json:"token"`
	RefteshToken string `json:"refresh_token"`
}
