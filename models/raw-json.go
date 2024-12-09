package models

type JSONRequestOTP struct {
	Request struct {
		Email string `json:"email"`
	} `json:"request"`
}

type JSONValidasiOTP struct {
	Request struct {
		Email    string `json:"email"`
		Otp      string `json:"otp"`
		Hp       string `json:"hp"`
		Password string `json:"password"`
	} `json:"request"`
}
