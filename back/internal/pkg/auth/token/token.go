package controller

import "time"

type Token struct {
	Token        string `json:"token"`
	RefteshToken string `json:"refresh_token"`
	ExpiredAt    int64  `json:"expired_at"`
}

func (t *Token) Expired() bool {
	return t.ExpiredAt < time.Now().Unix()
}
