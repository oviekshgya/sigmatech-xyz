package logging

import (
	"sigmatech-xyz/pkg"
	"time"
)

type MasterLog struct {
	Apps      string    `gorm:"column:apps" json:"apps"`
	IdLog     int       `gorm:"column:idLog" json:"id_log"`
	Tanggal   time.Time `gorm:"column:tanggal" json:"tanggal"`
	Method    string    `gorm:"column:method" json:"method"`
	Endpoint  string    `gorm:"column:endpoint" json:"endpoint"`
	Header    string    `gorm:"column:header" json:"header"`
	Body      string    `gorm:"column:body" json:"body"`
	Ip        string    `gorm:"column:ip" json:"ip"`
	MacAddr   string    `gorm:"column:macAddr" json:"macAddr"`
	IdUser    int       `gorm:"column:idUser" json:"idUser"`
	CreatedAt time.Time `gorm:"column:created_at;auto_now_add;type(date)"`
	UpdatedAt time.Time `gorm:"column:updated_at;auto_now_add;type(date)" json:"updatedAt"`
	Response  string    `gorm:"column:response" json:"response"`
}

func (MasterLog) TableName() string {
	return pkg.LOG
}
