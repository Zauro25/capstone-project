package models

import "gorm.io/gorm"

// AdminDPK merepresentasikan admin di Dinas Perpustakaan dan Kearsipan
type AdminDPK struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"` // Foreign key ke tabel users
	User   User `gorm:"foreignKey:UserID"`
}

// TableName untuk model AdminDPK
func (AdminDPK) TableName() string {
	return "admin_dpk"
}
