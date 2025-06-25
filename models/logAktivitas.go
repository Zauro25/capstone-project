package models

import (
	"time"

	"gorm.io/gorm"
)

// LogAktivitas merepresentasikan log aktivitas pengguna untuk audit trail
type LogAktivitas struct {
	gorm.Model
	UserID        uint      `gorm:"not null" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID"`
	TipeAktivitas string    `gorm:"type:varchar(100);not null" json:"tipe_aktivitas"` // Contoh: "Login", "Tambah Perpustakaan", "Verifikasi Data"
	Deskripsi     string    `gorm:"type:text" json:"deskripsi"`
	Timestamp     time.Time `gorm:"not null" json:"timestamp"`
	IPAddress     string    `gorm:"type:varchar(50)" json:"ip_address"`
}

// TableName untuk model LogAktivitas
func (LogAktivitas) TableName() string {
	return "log_aktivitas"
}
