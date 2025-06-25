package models

import "gorm.io/gorm"

type Verifikasi struct {
    gorm.Model
    PerpustakaanID  uint   `gorm:"not null"`
    Perpustakaan    Perpustakaan `gorm:"foreignKey:PerpustakaanID"`
    Status          string `gorm:"type:varchar(50);not null"`
    CatatanRevisi   string
    AdminDPKID      uint   `gorm:"not null"`
    AdminDPK        AdminDPK `gorm:"foreignKey:AdminDPKID"`
}