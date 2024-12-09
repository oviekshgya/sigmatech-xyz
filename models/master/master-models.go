package master

import (
	"sigmatech-xyz/pkg"
	"time"
)

type MasterMerchants struct {
	IdMerchant   int       `gorm:"column:idMerchant;autoIncrement;primaryKey" json:"idMerchant"`
	NamaMerchant string    `gorm:"column:namaMerchant;size:200" json:"namaMerchant"`
	IsActive     int       `gorm:"column:isActive" json:"isActive"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"updated_at;autoUpdateTime"`
}

func (MasterMerchants) TableName() string {
	return pkg.MERCHANT
}

type MasterRates struct {
	IdRate    int       `gorm:"column:idRate;autoIncrement;primaryKey" json:"idRate"`
	Rate      float64   `gorm:"column:rate;type:double" json:"rate"`
	Admin     float64   `gorm:"column:admin;type:double" json:"admin"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"updated_at;autoUpdateTime"`
}

func (MasterRates) TableName() string {
	return pkg.RATE
}
