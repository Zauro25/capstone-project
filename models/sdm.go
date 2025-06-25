package models

import "gorm.io/gorm"

// SDM merepresentasikan data Sumber Daya Manusia di perpustakaan
type SDM struct {
	gorm.Model
	PerpustakaanID uint   `gorm:"not null" json:"perpustakaan_id"`
	Nama           string `gorm:"not null" json:"nama" validate:"required"`
	Jabatan        string `gorm:"type:varchar(100)" json:"jabatan"`
	PendidikanTerakhir string `gorm:"type:varchar(50)" json:"pendidikan_terakhir"`
	StatusKepegawaian string `gorm:"type:varchar(50)" json:"status_kepegawaian"`
}

// TableName untuk model SDM
func (SDM) TableName() string {
	return "sdm"
}
