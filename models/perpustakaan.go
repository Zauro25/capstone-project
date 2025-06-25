package models

import (
	"gorm.io/gorm"
)

// Perpustakaan merepresentasikan entitas perpustakaan
type Perpustakaan struct {
	gorm.Model
	Nama            string `gorm:"not null" json:"nama" validate:"required"`
	Alamat          string `gorm:"not null" json:"alamat" validate:"required"`
	Jenis           string `gorm:"type:varchar(50);not null" json:"jenis" validate:"required,oneof=Umum Sekolah Khusus"`
	NomorInduk      string `gorm:"uniqueIndex;not null" json:"nomor_induk" validate:"required"`
	JumlahSDM       int    `json:"jumlah_sdm" validate:"gte=0"`
	JumlahPengunjung int    `json:"jumlah_pengunjung" validate:"gte=0"`
	JumlahAnggota   int    `json:"jumlah_anggota" validate:"gte=0"`
	StatusVerifikasi string `gorm:"type:varchar(50);default:'Belum Diverifikasi'" json:"status_verifikasi"` // Belum Diverifikasi, Dalam Proses Verifikasi, Disetujui, Perlu Revisi

	// Relasi
	SDM         []SDM         `gorm:"foreignKey:PerpustakaanID"`
	Pengunjung  []Pengunjung  `gorm:"foreignKey:PerpustakaanID"`
	Anggota     []Anggota     `gorm:"foreignKey:PerpustakaanID"`
	Verifikasi  []Verifikasi  `gorm:"foreignKey:IDData;polymorphic:DataType"` // Polymorphic untuk verifikasi data
	Revisi      []Revisi      `gorm:"foreignKey:IDData;polymorphic:DataType"`   // Polymorphic untuk revisi data
	AdminPerpustakaan AdminPerpustakaan `gorm:"foreignKey:PerpustakaanID"` // Satu perpustakaan memiliki satu admin perpustakaan
}

// TableName untuk model Perpustakaan
func (Perpustakaan) TableName() string {
	return "perpustakaans"
}
