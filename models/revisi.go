package models

import (
	"time"
	"gorm.io/gorm"
)

type Revisi struct {
	gorm.Model
	PerpustakaanID  uint      `gorm:"not null" json:"perpustakaan_id"`
	Perpustakaan    Perpustakaan `gorm:"foreignKey:PerpustakaanID"`
	CatatanDPK      string    `gorm:"not null" json:"catatan_dpk" validate:"required"`
	TanggalRevisiDPK time.Time `gorm:"not null;default:now()" json:"tanggal_revisi_dpk"`
	AdminDPKID      uint      `gorm:"not null" json:"admin_dpk_id"`
	AdminDPK        AdminDPK  `gorm:"foreignKey:AdminDPKID"`

	ResponAdminPerpustakaan string     `json:"respon_admin_perpustakaan"`
	TanggalResponAdminPerpustakaan *time.Time `gorm:"type:date" json:"tanggal_respon_admin_perpustakaan"`
	AdminPerpustakaanID     *uint      `json:"admin_perpustakaan_id"`
	AdminPerpustakaan       *AdminPerpustakaan `gorm:"foreignKey:AdminPerpustakaanID"`

	StatusRevisi string `gorm:"type:varchar(50);default:'Menunggu Respon'" json:"status_revisi"`
}

func (Revisi) TableName() string {
	return "revisis"
}