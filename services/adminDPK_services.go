package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"perpustakaan-api/models"
	"perpustakaan-api/utils"
	"time"

	"gorm.io/gorm"
)

// AdminDPKService menyediakan fungsi terkait manajemen Admin DPK
type AdminDPKService struct {
	DB *gorm.DB
}

// NewAdminDPKService membuat instance AdminDPKService baru
func NewAdminDPKService(db *gorm.DB) *AdminDPKService {
	return &AdminDPKService{DB: db}
}

// GetAllPerpustakaanForDPK mengambil semua data perpustakaan (termasuk status verifikasi)
func (s *AdminDPKService) GetAllPerpustakaanForDPK() ([]models.Perpustakaan, error) {
	var perpustakaans []models.Perpustakaan
	err := s.DB.Find(&perpustakaans).Error
	return perpustakaans, err
}

// VerifyPerpustakaan memverifikasi data perpustakaan
func (s *AdminDPKService) VerifyPerpustakaan(perpustakaanID, adminDPKUserID uint, status, catatan string) error {
	var perpustakaan models.Perpustakaan
	if err := s.DB.First(&perpustakaan, perpustakaanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("perpustakaan tidak ditemukan")
		}
		return err
	}

	// Dapatkan ID AdminDPK dari UserID
	var adminDPK models.AdminDPK
	if err := s.DB.Where("user_id = ?", adminDPKUserID).First(&adminDPK).Error; err != nil {
		return errors.New("admin DPK tidak ditemukan")
	}

	perpustakaan.StatusVerifikasi = status
	if err := s.DB.Save(&perpustakaan).Error; err != nil {
		return err
	}

	// Catat verifikasi
	verifikasi := models.Verifikasi{
		IDData:            perpustakaanID,
		DataType:          "Perpustakaan",
		Status:            status,
		CatatanRevisi:     catatan,
		TanggalVerifikasi: time.Now(),
		AdminDPKID:        adminDPK.ID,
	}
	return s.DB.Create(&verifikasi).Error
}

// RequestRevisi meminta revisi data perpustakaan
func (s *AdminDPKService) RequestRevisi(perpustakaanID, adminDPKUserID uint, catatan string) error {
	var perpustakaan models.Perpustakaan
	if err := s.DB.First(&perpustakaan, perpustakaanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("perpustakaan tidak ditemukan")
		}
		return err
	}

	// Dapatkan ID AdminDPK dari UserID
	var adminDPK models.AdminDPK
	if err := s.DB.Where("user_id = ?", adminDPKUserID).First(&adminDPK).Error; err != nil {
		return errors.New("admin DPK tidak ditemukan")
	}

	perpustakaan.StatusVerifikasi = "Perlu Revisi"
	if err := s.DB.Save(&perpustakaan).Error; err != nil {
		return err
	}

	revisi := models.Revisi{
		IDData:           perpustakaanID,
		DataType:         "Perpustakaan",
		CatatanDPK:       catatan,
		TanggalRevisiDPK: time.Now(),
		AdminDPKID:       adminDPK.ID,
		StatusRevisi:     "Menunggu Respon",
	}
	return s.DB.Create(&revisi).Error
}

// UploadLaporan mengunggah file laporan
func (s *AdminDPKService) UploadLaporan(file *multipart.FileHeader, periode, jenisLaporan string, userID uint) (*models.Laporan, error) {
	// Buat direktori penyimpanan jika belum ada
	uploadDir := "./uploads/laporan"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori upload: %w", err)
	}

	// Simpan file
	filename := fmt.Sprintf("%d_%s_%s%s", userID, periode, jenisLaporan, filepath.Ext(file.Filename))
	filePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat file tujuan: %w", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("gagal menulis file: %w", err)
	}

	laporan := models.Laporan{
		Periode:       periode,
		JenisLaporan:  jenisLaporan,
		NamaFile:      filename,
		PathFile:      filePath,
		TanggalGenerate: time.Now(),
		UserID:        userID,
	}

	if err := utils.ValidateStruct(laporan); err != nil {
		// Hapus file yang sudah diupload jika validasi gagal
		os.Remove(filePath)
		return nil, fmt.Errorf("validasi data laporan gagal: %w", err)
	}

	if err := s.DB.Create(&laporan).Error; err != nil {
		os.Remove(filePath) // Hapus file jika gagal disimpan ke DB
		return nil, fmt.Errorf("gagal menyimpan data laporan ke database: %w", err)
	}

	return &laporan, nil
}

// GetAllLaporan mengambil semua data laporan
func (s *AdminDPKService) GetAllLaporan() ([]models.Laporan, error) {
	var laporans []models.Laporan
	err := s.DB.Find(&laporans).Error
	return laporans, err
}

// ExportDataPerpustakaan mengekspor data perpustakaan ke format CSV (contoh sederhana)
func (s *AdminDPKService) ExportDataPerpustakaan() ([]byte, error) {
	var perpustakaans []models.Perpustakaan
	if err := s.DB.Find(&perpustakaans).Error; err != nil {
		return nil, err
	}

	// Buat header CSV
	csvContent := "ID,Nama,Alamat,Jenis,Nomor Induk,Jumlah SDM,Jumlah Pengunjung,Jumlah Anggota,Status Verifikasi\n"

	// Isi data
	for _, p := range perpustakaans {
		csvContent += fmt.Sprintf("%d,%s,%s,%s,%s,%d,%d,%d,%s\n",
			p.ID, p.Nama, p.Alamat, p.Jenis, p.NomorInduk, p.JumlahSDM, p.JumlahPengunjung, p.JumlahAnggota, p.StatusVerifikasi)
	}

	return []byte(csvContent), nil
}

// GetVerifikasiHistory mengambil riwayat verifikasi untuk perpustakaan tertentu
func (s *AdminDPKService) GetVerifikasiHistory(perpustakaanID uint) ([]models.Verifikasi, error) {
	var verifikasis []models.Verifikasi
	err := s.DB.Where("id_data = ? AND data_type = ?", perpustakaanID, "Perpustakaan").Order("tanggal_verifikasi desc").Find(&verifikasis).Error
	return verifikasis, err
}

// GetRevisiHistory mengambil riwayat revisi untuk perpustakaan tertentu
func (s *AdminDPKService) GetRevisiHistory(perpustakaanID uint) ([]models.Revisi, error) {
	var revisis []models.Revisi
	err := s.DB.Where("id_data = ? AND data_type = ?", perpustakaanID, "Perpustakaan").Order("tanggal_revisi_dpk desc").Find(&revisis).Error
	return revisis, err
}
