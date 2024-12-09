package models

type JSONRequestOTP struct {
	Request struct {
		Email string `json:"email"`
	} `json:"request"`
}
