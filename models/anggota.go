package models

import (
	"time"

	"gorm.io/gorm"
)

// Anggota merepresentasikan data anggota perpustakaan
type Anggota struct {
	gorm.Model
	PerpustakaanID uint      `gorm:"not null" json:"perpustakaan_id"`
	Nama           string    `gorm:"not null" json:"nama" validate:"required"`
	TanggalDaftar  time.Time `gorm:"not null" json:"tanggal_daftar" validate:"required"`
	StatusAktif    bool      `gorm:"default:true" json:"status_aktif"`
	JenisKelamin   string    `gorm:"type:varchar(10)" json:"jenis_kelamin"`
	Pekerjaan      string    `json:"pekerjaan"`
}

// TableName untuk model Anggota
func (Anggota) TableName() string {
	return "anggota"
}
