package repositories

import (
	"perpustakaan-api/models"

	"gorm.io/gorm"
)

type PerpustakaanRepository interface {
	CreatePerpustakaan(perpustakaan models.Perpustakaan) (models.Perpustakaan, error)
	GetAllPerpustakaan() ([]models.Perpustakaan, error)
	GetPerpustakaanByID(id uint) (models.Perpustakaan, error)
	UpdatePerpustakaan(perpustakaan models.Perpustakaan) (models.Perpustakaan, error)
	DeletePerpustakaan(id uint) error
	GetPerpustakaanByStatus(status string) ([]models.Perpustakaan, error)
}

type perpustakaanRepository struct {
	db *gorm.DB
}

func NewPerpustakaanRepository(db *gorm.DB) PerpustakaanRepository {
	return &perpustakaanRepository{db}
}

func (r *perpustakaanRepository) CreatePerpustakaan(perpustakaan models.Perpustakaan) (models.Perpustakaan, error) {
	err := r.db.Create(&perpustakaan).Error
	return perpustakaan, err
}

func (r *perpustakaanRepository) GetAllPerpustakaan() ([]models.Perpustakaan, error) {
	var perpustakaans []models.Perpustakaan
	err := r.db.Preload("Koleksi").Preload("SDM").Preload("Pengunjung").Preload("Anggota").Find(&perpustakaans).Error
	return perpustakaans, err
}

func (r *perpustakaanRepository) GetPerpustakaanByID(id uint) (models.Perpustakaan, error) {
	var perpustakaan models.Perpustakaan
	err := r.db.Preload("Koleksi").Preload("SDM").Preload("Pengunjung").Preload("Anggota").First(&perpustakaan, id).Error
	return perpustakaan, err
}

func (r *perpustakaanRepository) UpdatePerpustakaan(perpustakaan models.Perpustakaan) (models.Perpustakaan, error) {
	err := r.db.Save(&perpustakaan).Error
	return perpustakaan, err
}

func (r *perpustakaanRepository) DeletePerpustakaan(id uint) error {
	err := r.db.Delete(&models.Perpustakaan{}, id).Error
	return err
}

func (r *perpustakaanRepository) GetPerpustakaanByStatus(status string) ([]models.Perpustakaan, error) {
	var perpustakaans []models.Perpustakaan
	err := r.db.Where("status_verifikasi = ?", status).Find(&perpustakaans).Error
	return perpustakaans, err
}