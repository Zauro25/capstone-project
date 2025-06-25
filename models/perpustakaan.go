package models

import "gorm.io/gorm"

type Perpustakaan struct {
    gorm.Model
    Nama            string `gorm:"not null"`
    Alamat          string `gorm:"not null"`
    Jenis           string `gorm:"type:varchar(50);not null"`
    NomorInduk      string `gorm:"uniqueIndex;not null"`
    JumlahSDM       int    
    JumlahPengunjung int    
    JumlahAnggota   int    
    StatusVerifikasi string `gorm:"type:varchar(50);default:'Belum Diverifikasi'"`

    // Relasi yang diperbaiki
    SDM         []SDM         `gorm:"foreignKey:PerpustakaanID"`
    Pengunjung  []Pengunjung  `gorm:"foreignKey:PerpustakaanID"`
    Anggota     []Anggota     `gorm:"foreignKey:PerpustakaanID"`
    Verifikasi  []Verifikasi  `gorm:"foreignKey:PerpustakaanID"`
    Revisi      []Revisi      `gorm:"foreignKey:PerpustakaanID"`
    AdminPerpustakaan AdminPerpustakaan `gorm:"foreignKey:PerpustakaanID"`
}