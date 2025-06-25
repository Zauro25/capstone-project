package models

import "gorm.io/gorm"

// AdminPerpustakaan merepresentasikan admin lokal di setiap perpustakaan
type AdminPerpustakaan struct {
	gorm.Model
	UserID         uint `gorm:"uniqueIndex;not null" json:"user_id"` // Foreign key ke tabel users
	User           User `gorm:"foreignKey:UserID"`
	PerpustakaanID uint `gorm:"uniqueIndex;not null" json:"perpustakaan_id"` // Foreign key ke tabel perpustakaans
	Perpustakaan   *Perpustakaan `gorm:"foreignKey:PerpustakaanID"`
}

// TableName untuk model AdminPerpustakaan
func (AdminPerpustakaan) TableName() string {
	return "admin_perpustakaan"
}
