package models

import (
	"time"

	"gorm.io/gorm"
)

// Revisi merepresentasikan permintaan revisi dari Admin DPK dan respons dari Admin Perpustakaan
type Revisi struct {
	gorm.Model
	IDData        uint      `gorm:"not null" json:"id_data"` // ID dari data yang direvisi (misal: PerpustakaanID)
	DataType      string    `gorm:"type:varchar(50);not null" json:"data_type"` // Tipe data yang direvisi (misal: "Perpustakaan")
	CatatanDPK    string    `gorm:"not null" json:"catatan_dpk" validate:"required"` // Catatan dari Admin DPK
	TanggalRevisiDPK time.Time `gorm:"not null" json:"tanggal_revisi_dpk"`
	AdminDPKID    uint      `gorm:"not null" json:"admin_dpk_id"`
	AdminDPK      AdminDPK  `gorm:"foreignKey:AdminDPKID"`

	ResponAdminPerpustakaan string    `json:"respon_admin_perpustakaan"`
	TanggalResponAdminPerpustakaan *time.Time `json:"tanggal_respon_admin_perpustakaan"`
	AdminPerpustakaanID     *uint     `json:"admin_perpustakaan_id"`
	AdminPerpustakaan       *AdminPerpustakaan `gorm:"foreignKey:AdminPerpustakaanID"`

	StatusRevisi string `gorm:"type:varchar(50);default:'Menunggu Respon'" json:"status_revisi"` // Menunggu Respon, Selesai, Ditolak
}

// TableName untuk model Revisi
func (Revisi) TableName() string {
	return "revisis"
}
