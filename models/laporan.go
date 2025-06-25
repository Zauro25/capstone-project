package models

import (
	"time"

	"gorm.io/gorm"
)

// Laporan merepresentasikan data laporan yang dihasilkan atau diunggah
type Laporan struct {
	gorm.Model
	Periode       string    `gorm:"type:varchar(50);not null" json:"periode" validate:"required"` // Contoh: "Semester 1 2023"
	JenisLaporan  string    `gorm:"type:varchar(100);not null" json:"jenis_laporan" validate:"required"`
	NamaFile      string    `gorm:"not null" json:"nama_file"`
	PathFile      string    `gorm:"not null" json:"path_file"` // Path penyimpanan file di server
	TanggalGenerate time.Time `gorm:"not null" json:"tanggal_generate"`
	UserID        uint      `gorm:"not null" json:"user_id"` // User yang mengunggah/mengenerate laporan
	User          User      `gorm:"foreignKey:UserID"`
}

// TableName untuk model Laporan
func (Laporan) TableName() string {
	return "laporans"
}
