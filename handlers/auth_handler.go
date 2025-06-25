package handlers

import (
	"net/http"
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
	// Dalam kasus JWT, logout di sisi server biasanya hanya mencatat aktivitas.
	// Token harus dihapus di sisi klien.
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
