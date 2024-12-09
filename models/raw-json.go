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

type JSONTransaksiPengajuan struct {
	Request struct {
		OTR           float64 `json:"otr"`
		AdminFee      float64 `json:"adminFee"`
		JumlahCicilan float64 `json:"jumlahCicilan"`
		JumlahBunga   float64 `json:"jumlahBunga"`
		NamaAset      string  `json:"namaAset"`
	} `json:"request"`
}

type JSONTransaksiSimulasi struct {
	Request struct {
		OTR      float64 `json:"otr"`
		Tenor    int     `json:"tenor"`
		NamaAset string  `json:"namaAset"`
	} `json:"request"`
}

type JSONTransaksiPinjaman struct {
	Request struct {
		OTR        float64 `json:"otr"`
		Tenor      int     `json:"tenor"`
		NamaAset   string  `json:"namaAset"`
		IdMerchant int     `json:"idMerchant"`
	} `json:"request"`
}

type JSONTransaksiPayment struct {
	Request struct {
		NoKontrak      string   `json:"noKontrak"`
		DetailAngsuran []string `json:"detailAngsuran"`
	} `json:"request"`
}
