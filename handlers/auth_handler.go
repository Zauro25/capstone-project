package handlers

import (
	"net/http"
	"perpustakaan-api/models"
	"perpustakaan-api/services"
	"perpustakaan-api/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler menangani permintaan terkait autentikasi
type AuthHandler struct {
	AuthService *services.AuthService
}

// NewAuthHandler membuat instance AuthHandler baru
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// RegisterRequest merepresentasikan payload untuk pendaftaran pengguna
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=AdminPerpustakaan AdminDPK Executive"`
}

// RegisterUser menangani pendaftaran pengguna baru
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Password, req.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mendaftarkan pengguna", err.Error())
		return
	}

	// Buat entitas spesifik berdasarkan role
	switch req.Role {
	case "AdminPerpustakaan":
		// Admin Perpustakaan memerlukan perpustakaan_id
		var adminReq struct {
			PerpustakaanID uint `json:"perpustakaan_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&adminReq); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Perpustakaan ID diperlukan untuk Admin Perpustakaan", nil)
			return
		}
		admin := models.AdminPerpustakaan{
			UserID:         user.ID,
			PerpustakaanID: adminReq.PerpustakaanID,
		}
		if err := h.AuthService.DB.Create(&admin).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat admin perpustakaan", err.Error())
			return
		}
	case "AdminDPK":
		admin := models.AdminDPK{
			UserID: user.ID,
		}
		if err := h.AuthService.DB.Create(&admin).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat admin DPK", err.Error())
			return
		}
	// Executive tidak memerlukan entitas tambahan
	}

	utils.SuccessResponse(c, http.StatusCreated, "Pengguna berhasil didaftarkan", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

// LoginRequest merepresentasikan payload untuk login pengguna
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginUser menangani login pengguna
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	token, err := h.AuthService.LoginUser(req.Username, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login gagal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", gin.H{"token": token})
}

// LogoutUser menangani logout pengguna
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ID pengguna tidak ditemukan di konteks", nil)
		return
	}

	err := h.AuthService.LogoutUser(userID.(uint), c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mencatat logout", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logout berhasil", nil)
}

// RegisterAdminPerpustakaan khusus untuk pendaftaran admin perpustakaan
func (h *AuthHandler) RegisterAdminPerpustakaan(c *gin.Context) {
	var req struct {
		Username       string `json:"username" binding:"required"`
		Password       string `json:"password" binding:"required"`
		PerpustakaanID uint   `json:"perpustakaan_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Password, "AdminPerpustakaan")
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mendaftarkan admin perpustakaan", err.Error())
		return
	}

	admin := models.AdminPerpustakaan{
		UserID:         user.ID,
		PerpustakaanID: req.PerpustakaanID,
	}
	if err := h.AuthService.DB.Create(&admin).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat admin perpustakaan", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Admin perpustakaan berhasil didaftarkan", gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"role":           user.Role,
		"perpustakaan_id": req.PerpustakaanID,
	})
}

// RegisterAdminDPK khusus untuk pendaftaran admin DPK
func (h *AuthHandler) RegisterAdminDPK(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Password, "AdminDPK")
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mendaftarkan admin DPK", err.Error())
		return
	}

	admin := models.AdminDPK{
		UserID: user.ID,
	}
	if err := h.AuthService.DB.Create(&admin).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat admin DPK", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Admin DPK berhasil didaftarkan", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

// RegisterExecutive khusus untuk pendaftaran executive
func (h *AuthHandler) RegisterExecutive(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Password, "Executive")
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mendaftarkan executive", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Executive berhasil didaftarkan", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}