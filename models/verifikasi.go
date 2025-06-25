package models

import (
	"time"

	"gorm.io/gorm"
)

// Verifikasi merepresentasikan proses verifikasi data oleh Admin DPK
type Verifikasi struct {
	gorm.Model
	IDData        uint      `gorm:"not null" json:"id_data"` // ID dari data yang diverifikasi (misal: PerpustakaanID)
	DataType      string    `gorm:"type:varchar(50);not null" json:"data_type"` // Tipe data yang diverifikasi (misal: "Perpustakaan")
	Status        string    `gorm:"type:varchar(50);not null" json:"status" validate:"required,oneof=Disetujui Perlu Revisi"` // Disetujui, Perlu Revisi
	CatatanRevisi string    `json:"catatan_revisi"`
	TanggalVerifikasi time.Time `gorm:"not null" json:"tanggal_verifikasi"`
	AdminDPKID    uint      `gorm:"not null" json:"admin_dpk_id"` // ID Admin DPK yang melakukan verifikasi
	AdminDPK      AdminDPK  `gorm:"foreignKey:AdminDPKID"`
}

// TableName untuk model Verifikasi
func (Verifikasi) TableName() string {
	return "verifikasis"
}
