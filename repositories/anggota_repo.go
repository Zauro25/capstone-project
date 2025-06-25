package repositories

import (
    "perpustakaan-api/models"

    "gorm.io/gorm"
)

type AnggotaRepository interface {
    CreateAnggota(anggota models.Anggota) (models.Anggota, error)
    GetAllAnggota() ([]models.Anggota, error)
    GetAnggotaByID(id uint) (models.Anggota, error)
    UpdateAnggota(anggota models.Anggota) (models.Anggota, error)
    DeleteAnggota(id uint) error
    GetAnggotaByPerpustakaan(perpustakaanID uint) ([]models.Anggota, error)
}

type anggotaRepository struct {
    db *gorm.DB
}

func NewAnggotaRepository(db *gorm.DB) AnggotaRepository {
    return &anggotaRepository{db}
}

func (r *anggotaRepository) CreateAnggota(anggota models.Anggota) (models.Anggota, error) {
    err := r.db.Create(&anggota).Error
    return anggota, err
}

func (r *anggotaRepository) GetAllAnggota() ([]models.Anggota, error) {
    var anggotas []models.Anggota
    err := r.db.Find(&anggotas).Error
    return anggotas, err
}

func (r *anggotaRepository) GetAnggotaByID(id uint) (models.Anggota, error) {
    var anggota models.Anggota
    err := r.db.First(&anggota, id).Error
    return anggota, err
}

func (r *anggotaRepository) UpdateAnggota(anggota models.Anggota) (models.Anggota, error) {
    err := r.db.Save(&anggota).Error
    return anggota, err
}

func (r *anggotaRepository) DeleteAnggota(id uint) error {
    err := r.db.Delete(&models.Anggota{}, id).Error
    return err
}

func (r *anggotaRepository) GetAnggotaByPerpustakaan(perpustakaanID uint) ([]models.Anggota, error) {
    var anggotas []models.Anggota
    err := r.db.Where("perpustakaan_id = ?", perpustakaanID).Find(&anggotas).Error
    return anggotas, err
}