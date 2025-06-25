package services

import (
	"errors"
	"perpustakaan-api/models"
	"perpustakaan-api/utils"
	"time"

	"gorm.io/gorm"
)

// AuthService menyediakan fungsi terkait autentikasi pengguna
type AuthService struct {
	DB *gorm.DB
}

// NewAuthService membuat instance AuthService baru
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

// RegisterUser mendaftarkan pengguna baru
func (s *AuthService) RegisterUser(username, password, role string) (*models.User, error) {
	// Cek apakah username sudah ada
	var existingUser models.User
	if s.DB.Where("username = ?", username).First(&existingUser).Error == nil {
		return nil, errors.New("username sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	user := models.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, errors.New("gagal mendaftarkan pengguna")
	}

	// Jika role adalah AdminPerpustakaan, buat entri di tabel admin_perpustakaan
	if role == "AdminPerpustakaan" {
		// TODO: Perlu mekanisme untuk mengaitkan AdminPerpustakaan dengan PerpustakaanID yang ada
		// Untuk saat ini, kita asumsikan PerpustakaanID akan diisi secara manual atau melalui endpoint terpisah
		adminPerpustakaan := models.AdminPerpustakaan{
			UserID: user.ID,
			// PerpustakaanID: 0, // Ini harus diisi nanti
		}
		s.DB.Create(&adminPerpustakaan)
	} else if role == "AdminDPK" {
		adminDPK := models.AdminDPK{
			UserID: user.ID,
		}
		s.DB.Create(&adminDPK)
	}

	return &user, nil
}

// LoginUser mengautentikasi pengguna dan mengembalikan token JWT
func (s *AuthService) LoginUser(username, password string) (string, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("username atau password salah")
		}
		return "", errors.New("terjadi kesalahan saat mencari pengguna")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("username atau password salah")
	}

	// Update last login time
	now := time.Now()
	user.LastLogin = &now
	s.DB.Save(&user)

	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		return "", errors.New("gagal membuat token autentikasi")
	}

	return token, nil
}

// LogoutUser (Untuk API stateless dengan JWT, logout biasanya hanya berarti menghapus token di sisi klien)
// Namun, kita bisa mencatat aktivitas logout di audit trail.
func (s *AuthService) LogoutUser(userID uint, ipAddress string) error {
	// Catat aktivitas logout
	logAktivitas := models.LogAktivitas{
		UserID:        userID,
		TipeAktivitas: "Logout",
		Deskripsi:     "Pengguna berhasil logout",
		Timestamp:     time.Now(),
		IPAddress:     ipAddress,
	}
	if err := s.DB.Create(&logAktivitas).Error; err != nil {
		return errors.New("gagal mencatat log aktivitas logout")
	}
	return nil
}
