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

type JSONLogin struct {
	Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"request"`
}

type JSONVerifikasi struct {
	Request struct {
		Nik          string  `json:"nik"`
		LegalName    string  `json:"legalName"`
		TempatLahir  string  `json:"tempatLahir"`
		TanggalLahir string  `json:"tanggalLahir"`
		Salary       float64 `json:"salary"`
		FotoKtp      string  `json:"fotoKtp"`
		FotoSelfie   string  `json:"fotoSelfie"`
		IsAktivasi   int     `json:"isAktivasi"`
	} `json:"request"`
}
