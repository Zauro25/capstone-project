package services

import (
	"errors"
	"fmt"
	"perpustakaan-api/models"
	"perpustakaan-api/utils"

	"gorm.io/gorm"
)

// PerpustakaanService menyediakan fungsi terkait manajemen data perpustakaan
type PerpustakaanService struct {
	DB *gorm.DB
}

// NewPerpustakaanService membuat instance PerpustakaanService baru
func NewPerpustakaanService(db *gorm.DB) *PerpustakaanService {
	return &PerpustakaanService{DB: db}
}

// CreatePerpustakaan membuat entri perpustakaan baru
func (s *PerpustakaanService) CreatePerpustakaan(perpustakaan *models.Perpustakaan) error {
	if err := utils.ValidateStruct(perpustakaan); err != nil {
		return fmt.Errorf("validasi data gagal: %w", err)
	}
	return s.DB.Create(perpustakaan).Error
}

// GetAllPerpustakaan mengambil semua data perpustakaan
func (s *PerpustakaanService) GetAllPerpustakaan() ([]models.Perpustakaan, error) {
	var perpustakaans []models.Perpustakaan
	err := s.DB.Find(&perpustakaans).Error
	return perpustakaans, err
}

// GetPerpustakaanByID mengambil data perpustakaan berdasarkan ID
func (s *PerpustakaanService) GetPerpustakaanByID(id uint) (*models.Perpustakaan, error) {
	var perpustakaan models.Perpustakaan
	err := s.DB.First(&perpustakaan, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("perpustakaan tidak ditemukan")
	}
	return &perpustakaan, err
}

// UpdatePerpustakaan memperbarui data perpustakaan
func (s *PerpustakaanService) UpdatePerpustakaan(id uint, updatedPerpustakaan *models.Perpustakaan) error {
	var perpustakaan models.Perpustakaan
	if err := s.DB.First(&perpustakaan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("perpustakaan tidak ditemukan")
		}
		return err
	}

	// Update fields
	perpustakaan.Nama = updatedPerpustakaan.Nama
	perpustakaan.Alamat = updatedPerpustakaan.Alamat
	perpustakaan.Jenis = updatedPerpustakaan.Jenis
	perpustakaan.NomorInduk = updatedPerpustakaan.NomorInduk
	perpustakaan.JumlahSDM = updatedPerpustakaan.JumlahSDM
	perpustakaan.JumlahPengunjung = updatedPerpustakaan.JumlahPengunjung
	perpustakaan.JumlahAnggota = updatedPerpustakaan.JumlahAnggota
	// StatusVerifikasi tidak diupdate oleh Admin Perpustakaan

	if err := utils.ValidateStruct(perpustakaan); err != nil {
		return fmt.Errorf("validasi data gagal: %w", err)
	}

	return s.DB.Save(&perpustakaan).Error
}

// DeletePerpustakaan menghapus data perpustakaan
func (s *PerpustakaanService) DeletePerpustakaan(id uint) error {
	result := s.DB.Delete(&models.Perpustakaan{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("perpustakaan tidak ditemukan")
	}
	return nil
}

// --- Fungsi untuk sub-entitas (Koleksi, SDM, Pengunjung, Anggota) ---


// AddSDM menambahkan SDM ke perpustakaan tertentu
func (s *PerpustakaanService) AddSDM(sdm *models.SDM) error {
	if err := utils.ValidateStruct(sdm); err != nil {
		return fmt.Errorf("validasi data gagal: %w", err)
	}
	return s.DB.Create(sdm).Error
}

// GetSDMByPerpustakaanID mengambil semua SDM dari perpustakaan tertentu
func (s *PerpustakaanService) GetSDMByPerpustakaanID(perpustakaanID uint) ([]models.SDM, error) {
	var sdms []models.SDM
	err := s.DB.Where("perpustakaan_id = ?", perpustakaanID).Find(&sdms).Error
	return sdms, err
}

// UpdateSDM memperbarui data SDM
func (s *PerpustakaanService) UpdateSDM(id uint, updatedSDM *models.SDM) error {
	var sdm models.SDM
	if err := s.DB.First(&sdm, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("SDM tidak ditemukan")
		}
		return err
	}
	sdm.Nama = updatedSDM.Nama
	sdm.Jabatan = updatedSDM.Jabatan
	sdm.PendidikanTerakhir = updatedSDM.PendidikanTerakhir
	sdm.StatusKepegawaian = updatedSDM.StatusKepegawaian

	if err := utils.ValidateStruct(sdm); err != nil {
		return fmt.Errorf("validasi data gagal: %w", err)
	}
	return s.DB.Save(&sdm).Error
}

// DeleteSDM menghapus data SDM
func (s *PerpustakaanService) DeleteSDM(id uint) error {
	result := s.DB.Delete(&models.SDM{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("SDM tidak ditemukan")
	}
	return nil
}

// AddPengunjung menambahkan data pengunjung ke perpustakaan tertentu
func (s *PerpustakaanService) AddPengunjung(pengunjung *models.Pengunjung) error {
	if err := utils.ValidateStruct(pengunjung); err != nil {
		return fmt.Errorf("validasi data gagal: %w", err)
	}
	return s.DB.Create(pengunjung).Error
}

// GetPengunjungByPerpustakaanID mengambil semua data pengunjung dari perpustakaan tertentu
func (s *PerpustakaanService) GetPengunjungByPerpustakaanID(perpustakaanID uint) ([]models.Pengunjung, error) {
	var pengunjung []models.Pengunjung
	err := s.DB.Where("perpustakaan_id = ?", perpustakaanID).Find(&pengunjung).Error
	return pengunjung, err
}

// UpdatePengunjung memperbarui data pengunjung
func (s *PerpustakaanService) UpdatePengunjung(id uint, updatedPengunjung *models.Pengunjung) error {
	var pengunjung models.Pengunjung
	if err := s.DB.First(&pengunjung, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("data pengunjung tidak ditemukan")
		}
		return err
	}
	pengunjung.TanggalKunjungan = updatedPengunjung.TanggalKunjungan
	pengunjung.JumlahPengunjung = updatedPengunjung.JumlahPengunjung
	pengunjung.Keterangan = updatedPengunjung.Keterangan

	if err := utils.ValidateStruct(pengunjung); err != nil {
		return fmt.Errorf("validasi data gagal: %w", err)
	}
	return s.DB.Save(&pengunjung).Error
}

// DeletePengunjung menghapus data pengunjung
func (s *PerpustakaanService) DeletePengunjung(id uint) error {
	result := s.DB.Delete(&models.Pengunjung{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("data pengunjung tidak ditemukan")
	}
	return nil
}

// AddAnggota menambahkan data anggota ke perpustakaan tertentu
func (s *PerpustakaanService) AddAnggota(anggota *models.Anggota) error {
	if err := utils.ValidateStruct(anggota); err != nil {
		return fmt.Errorf("validasi data gagal: %w", err)
	}
	return s.DB.Create(anggota).Error
}

// GetAnggotaByPerpustakaanID mengambil semua data anggota dari perpustakaan tertentu
func (s *PerpustakaanService) GetAnggotaByPerpustakaanID(perpustakaanID uint) ([]models.Anggota, error) {
	var anggota []models.Anggota
	err := s.DB.Where("perpustakaan_id = ?", perpustakaanID).Find(&anggota).Error
	return anggota, err
}

// UpdateAnggota memperbarui data anggota
func (s *PerpustakaanService) UpdateAnggota(id uint, updatedAnggota *models.Anggota) error {
	var anggota models.Anggota
	if err := s.DB.First(&anggota, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("data anggota tidak ditemukan")
		}
		return err
	}
	anggota.Nama = updatedAnggota.Nama
	anggota.TanggalDaftar = updatedAnggota.TanggalDaftar
	anggota.StatusAktif = updatedAnggota.StatusAktif
	anggota.JenisKelamin = updatedAnggota.JenisKelamin
	anggota.Pekerjaan = updatedAnggota.Pekerjaan

	if err := utils.ValidateStruct(anggota); err != nil {
		return fmt.Errorf("validasi data gagal: %w", err)
	}
	return s.DB.Save(&anggota).Error
}

// DeleteAnggota menghapus data anggota
func (s *PerpustakaanService) DeleteAnggota(id uint) error {
	result := s.DB.Delete(&models.Anggota{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("data anggota tidak ditemukan")
	}
	return nil
}

// SubmitDataForVerification mengubah status perpustakaan menjadi "Dalam Proses Verifikasi"
func (s *PerpustakaanService) SubmitDataForVerification(perpustakaanID uint) error {
	var perpustakaan models.Perpustakaan
	if err := s.DB.First(&perpustakaan, perpustakaanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("perpustakaan tidak ditemukan")
		}
		return err
	}

	// Cek kelengkapan data sebelum submit (contoh sederhana)
	if perpustakaan.Nama == "" || perpustakaan.Alamat == "" || perpustakaan.NomorInduk == "" {
		return errors.New("data perpustakaan belum lengkap, tidak dapat diajukan verifikasi")
	}
	// Anda bisa menambahkan validasi lebih lanjut untuk koleksi, SDM, dll.

	perpustakaan.StatusVerifikasi = "Dalam Proses Verifikasi"
	return s.DB.Save(&perpustakaan).Error
}
