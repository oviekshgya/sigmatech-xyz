package transaksi

import (
	"sigmatech-xyz/pkg"
	"time"
)

type Transaksi struct {
	IdTransaksi    int       `gorm:"column:idTransaksi;autoIncrement;primaryKey" json:"idTransaksi"`
	IdMerchant     int       `gorm:"column:idMerchant;foreignKey" json:"idMerchant"`
	IdUserCustomer int       `gorm:"column:idUserCustomer;foreignKey" json:"idUserCustomer"`
	NoKontrak      string    `gorm:"column:noKontrak;size:200" json:"noKontrak"`
	OTR            float64   `gorm:"column:otr;type:double" json:"otr"`
	AdminFee       float64   `gorm:"column:adminFee;type:double" json:"adminFee"`
	JumlahCicilan  float64   `gorm:"column:jumlahCicilan;type:double" json:"jumlahCicilan"`
	JumlahBunga    float64   `gorm:"column:jumlahBunga;type:double" json:"jumlahBunga"`
	NamaAset       string    `gorm:"column:namaAset;size:200" json:"namaAset"`
	Tenor          int       `gorm:"column:tenor;size:200" json:"tenor"`
	Status         string    `gorm:"column:status;size:200" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"updated_at;autoUpdateTime"`
	TglJatuhTempo  time.Time `gorm:"column:tglJatuhTempo;type:date"`
}

type PaymentTransaksi struct {
	IdPaymentTransaksi int        `gorm:"column:idPaymentTransaksi;autoIncrement;primaryKey" json:"idPaymentTransaksi"`
	IdTransaksi        int        `gorm:"column:idTransaksi;foreignKey" json:"idTransaksi"`
	Tanggal            time.Time  `gorm:"column:tanggal;type:datetime" json:"tanggal"`
	TanggalJatuhTempo  time.Time  `gorm:"column:tanggalJatuhTempo" json:"tanggalJatuhTempo"`
	AngsuranKe         int        `gorm:"column:angsuranKe" json:"angsuranKe"`
	Bunga              float64    `gorm:"column:bunga;type:double" json:"bunga"`
	JumlahCicilan      float64    `gorm:"column:jumlahCicilan;type:double" json:"jumlahCicilan"`
	TotalCicilan       float64    `gorm:"column:totalCicilan;type:double" json:"totalCicilan"`
	CreatedAt          time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time  `gorm:"updated_at;autoUpdateTime"`
	Status             string     `gorm:"column:status;size:200" json:"status"`
	TglBayar           *time.Time `gorm:"column:tglBayar;type:datetime" json:"tglBayar"`
}

func (Transaksi) TableName() string {
	return pkg.TRANSACTION
}

func (PaymentTransaksi) TableName() string {
	return pkg.PAYMENTTRANSAKSI
}

type DataTransaksi struct {
	IdTransaksi     int         `gorm:"column:idTransaksi" json:"-"`
	Status          string      `gorm:"column:status;size:200" json:"status"`
	NoKontrak       string      `gorm:"column:noKontrak;size:200" json:"noKontrak"`
	OTR             float64     `gorm:"column:otr;type:double" json:"otr"`
	TglJatuhTempo   time.Time   `gorm:"column:tglJatuhTempo;type:date"`
	JumlahCicilan   float64     `gorm:"column:jumlahCicilan;type:double" json:"jumlahCicilan"`
	PaymentSchedule interface{} `gorm:"-" json:"paymentSchedule,omitempty"`
}

type CheckPayment struct {
	IdPaymentTransaksi int       `gorm:"column:idPaymentTransaksi;autoIncrement;primaryKey" json:"idPaymentTransaksi,omitempty"`
	IdTransaksi        int       `gorm:"column:idTransaksi;foreignKey" json:"idTransaksi,omitempty"`
	Tanggal            time.Time `gorm:"column:tanggal;type:datetime" json:"tanggal,omitempty"`
	TanggalJatuhTempo  time.Time `gorm:"column:tanggalJatuhTempo" json:"tanggalJatuhTempo,omitempty"`
	AngsuranKe         int       `gorm:"column:angsuranKe" json:"angsuranKe,omitempty"`
	Bunga              float64   `gorm:"column:bunga;type:double" json:"bunga,omitempty"`
	JumlahCicilan      float64   `gorm:"column:jumlahCicilan;type:double" json:"jumlahCicilan,omitempty"`
	TotalCicilan       float64   `gorm:"column:totalCicilan;type:double" json:"totalCicilan,omitempty"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt          time.Time `gorm:"updated_at;autoUpdateTime" json:"-"`
	Status             string    `gorm:"column:status;size:200" json:"status,omitempty"`
	TglBayar           time.Time `gorm:"column:tglBayar;type:datetime" json:"tglBayar,omitempty"`
}
