package services

import (
	"perpustakaan-api/models"

	"gorm.io/gorm"
)

// DashboardService menyediakan fungsi untuk data dashboard
type DashboardService struct {
	DB *gorm.DB
}

// NewDashboardService membuat instance DashboardService baru
func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{DB: db}
}

// GetDashboardStats mengambil statistik umum untuk dashboard
func (s *DashboardService) GetDashboardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total Perpustakaan
	var totalPerpustakaan int64
	s.DB.Model(&models.Perpustakaan{}).Count(&totalPerpustakaan)
	stats["total_perpustakaan"] = totalPerpustakaan

	// Status Verifikasi
	var totalDisetujui int64
	s.DB.Model(&models.Perpustakaan{}).Where("status_verifikasi = ?", "Disetujui").Count(&totalDisetujui)
	stats["total_disetujui"] = totalDisetujui

	var totalDalamProses int64
	s.DB.Model(&models.Perpustakaan{}).Where("status_verifikasi = ?", "Dalam Proses Verifikasi").Count(&totalDalamProses)
	stats["total_dalam_proses_verifikasi"] = totalDalamProses

	var totalPerluRevisi int64
	s.DB.Model(&models.Perpustakaan{}).Where("status_verifikasi = ?", "Perlu Revisi").Count(&totalPerluRevisi)
	stats["total_perlu_revisi"] = totalPerluRevisi

	var totalBelumDiverifikasi int64
	s.DB.Model(&models.Perpustakaan{}).Where("status_verifikasi = ?", "Belum Diverifikasi").Count(&totalBelumDiverifikasi)
	stats["total_belum_diverifikasi"] = totalBelumDiverifikasi

	// Total Koleksi

	// Total SDM
	var totalSDM int64
	s.DB.Model(&models.SDM{}).Count(&totalSDM)
	stats["total_sdm"] = totalSDM

	// Total Pengunjung (akumulasi)
	var totalPengunjung int64
	s.DB.Model(&models.Pengunjung{}).Select("SUM(jumlah_pengunjung)").Scan(&totalPengunjung)
	stats["total_pengunjung_akumulasi"] = totalPengunjung

	// Total Anggota
	var totalAnggota int64
	s.DB.Model(&models.Anggota{}).Count(&totalAnggota)
	stats["total_anggota"] = totalAnggota

	return stats, nil
}
