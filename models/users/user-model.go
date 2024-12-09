package users

import (
	"sigmatech-xyz/pkg"
	"time"
)

type UserCustomer struct {
	IdUserCustomer int       `gorm:"column:idUserCustomer;primaryKey;autoIncrement" json:"idUserCustomer"`
	Nik            string    `gorm:"column:nik;size:100" json:"nik"`
	IdAkun         int       `gorm:"column:idAkun;foreignKey" json:"idAkun"`
	LegalName      string    `gorm:"column:legalName;size:255" json:"legalName"`
	TempatLahir    string    `gorm:"column:tempatLahir;size:100" json:"tempatLahir"`
	TanggalLahir   time.Time `gorm:"column:tanggalLahir;type:date" json:"tanggalLahir"`
	Salary         float64   `gorm:"column:salary;type:double" json:"salary"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"updated_at;autoUpdateTime"`
	FotoKtp        string    `gorm:"column:fotoKtp;size:255" json:"fotoKtp"`
	FotoSelfie     string    `gorm:"column:fotoSelfie;size:255" json:"fotoSelfie"`
	IsAktivasi     int       `gorm:"column:isAktivasi;" json:"isAktivasi"`
}

func (UserCustomer) TableName() string {
	return pkg.USERSCUSTOMER
}

type AkunCustomer struct {
	IdAkun      int       `gorm:"column:idAkun;primaryKey;autoIncrement" json:"idAkun"`
	NamaLengkap string    `gorm:"column:namaLengkap;size:255" json:"namaLengkap"`
	Hp          string    `gorm:"column:hp;size:20" json:"hp"`
	Password    string    `gorm:"column:password;size:255" json:"password"`
	Email       string    `gorm:"column:email;size:255" json:"email"`
	IsActive    int       `gorm:"column:isActive" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"updated_at;autoUpdateTime"`
}

func (AkunCustomer) TableName() string {
	return pkg.AKUNCUSTOMER
}

type OTPCustomer struct {
	IdOtp     int       `gorm:"column:idOtp;primaryKey;autoIncrement" json:"idOtp"`
	Email     string    `gorm:"column:email;size:255" json:"email"`
	IsUsed    int       `gorm:"column:isUsed" json:"isUsed"`
	Kode      string    `gorm:"column:kode;size:255" json:"kode"`
	ExpiredAt time.Time `gorm:"column:expiredAt;type:datetime" json:"expiredAt"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"updated_at;autoUpdateTime"`
}

func (OTPCustomer) TableName() string {
	return pkg.OTPCUSTOMERS
}

type UserCustomerData struct {
	IdAkun      int           `gorm:"column:idAkun;primaryKey;autoIncrement" json:"idAkun"`
	NamaLengkap string        `gorm:"column:namaLengkap;size:255" json:"namaLengkap"`
	Hp          string        `gorm:"column:hp;size:20" json:"hp"`
	Password    string        `gorm:"column:password;size:255" json:"password"`
	Email       string        `gorm:"column:email;size:255" json:"email"`
	IsActive    int           `gorm:"column:isActive" json:"isActive"`
	CreatedAt   time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"updated_at;autoUpdateTime"`
	Data        *UserCustomer `gorm:"foreignKey:idAkun;references:IdAkun" json:"data"`
}

func (UserCustomerData) TableName() string {
	return pkg.AKUNCUSTOMER
}

type UserLimits struct {
	IdUserLimit int       `gorm:"column:idUserLimit;primaryKey" json:"idUserLimit"`
	Limit       float64   `gorm:"column:limit;type:double" json:"limit"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"updated_at;autoUpdateTime"`
}

func (UserLimits) TableName() string {
	return pkg.USERLIMIT
}
