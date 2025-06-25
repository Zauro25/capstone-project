package models

import (
	"time"

	"gorm.io/gorm"
)

// Pengunjung merepresentasikan data pengunjung perpustakaan
type Pengunjung struct {
	gorm.Model
	PerpustakaanID uint      `gorm:"not null" json:"perpustakaan_id"`
	TanggalKunjungan time.Time `gorm:"not null" json:"tanggal_kunjungan" validate:"required"`
	JumlahPengunjung int       `gorm:"not null" json:"jumlah_pengunjung" validate:"required,gte=0"`
	Keterangan       string    `json:"keterangan"`
}

// TableName untuk model Pengunjung
func (Pengunjung) TableName() string {
	return "pengunjung"
}
