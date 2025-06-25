package handlers

import (
	"fmt"
	"net/http"
	"perpustakaan-api/services"
	"perpustakaan-api/utils"
	"strconv"
	"time"
	"perpustakaan-api/models"
	"github.com/gin-gonic/gin"
)

// AdminDPKHandler menangani permintaan terkait Admin DPK
type AdminDPKHandler struct {
	AdminDPKService     *services.AdminDPKService
	PerpustakaanService *services.PerpustakaanService // Untuk merespon revisi
}

// NewAdminDPKHandler membuat instance AdminDPKHandler baru
func NewAdminDPKHandler(adminDPKService *services.AdminDPKService, perpustakaanService *services.PerpustakaanService) *AdminDPKHandler {
	return &AdminDPKHandler{
		AdminDPKService:     adminDPKService,
		PerpustakaanService: perpustakaanService,
	}
}

// GetAllPerpustakaanForDPK mengambil semua data perpustakaan untuk Admin DPK
func (h *AdminDPKHandler) GetAllPerpustakaanForDPK(c *gin.Context) {
	perpustakaans, err := h.AdminDPKService.GetAllPerpustakaanForDPK()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data perpustakaan", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data perpustakaan berhasil diambil", perpustakaans)
}

// VerifyPerpustakaanRequest merepresentasikan payload untuk verifikasi perpustakaan
type VerifyPerpustakaanRequest struct {
	Status string `json:"status" binding:"required,oneof=Disetujui Perlu Revisi"`
	Catatan string `json:"catatan"`
}

// VerifyPerpustakaan memverifikasi data perpustakaan oleh Admin DPK
func (h *AdminDPKHandler) VerifyPerpustakaan(c *gin.Context) {
	idStr := c.Param("id")
	perpustakaanID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	var req VerifyPerpustakaanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	adminDPKUserID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ID Admin DPK tidak ditemukan di konteks", nil)
		return
	}

	err = h.AdminDPKService.VerifyPerpustakaan(uint(perpustakaanID), adminDPKUserID.(uint), req.Status, req.Catatan)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memverifikasi perpustakaan", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Verifikasi perpustakaan berhasil", nil)
}

// RequestRevisi meminta revisi data perpustakaan oleh Admin DPK
func (h *AdminDPKHandler) RequestRevisi(c *gin.Context) {
	idStr := c.Param("id")
	perpustakaanID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	var req struct {
		Catatan string `json:"catatan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	adminDPKUserID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ID Admin DPK tidak ditemukan di konteks", nil)
		return
	}

	err = h.AdminDPKService.RequestRevisi(uint(perpustakaanID), adminDPKUserID.(uint), req.Catatan)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal meminta revisi", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permintaan revisi berhasil dikirim", nil)
}

// RespondToRevisi (oleh Admin Perpustakaan)
func (h *AdminDPKHandler) RespondToRevisi(c *gin.Context) {
	revisiIDStr := c.Param("revisi_id")
	revisiID, err := strconv.ParseUint(revisiIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID revisi tidak valid", nil)
		return
	}

	var req struct {
		Respon string `json:"respon" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	adminPerpustakaanUserID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ID Admin Perpustakaan tidak ditemukan di konteks", nil)
		return
	}

	// Dapatkan ID AdminPerpustakaan dari UserID
	var adminPerpustakaan models.AdminPerpustakaan
	if err := h.PerpustakaanService.DB.Where("user_id = ?", adminPerpustakaanUserID).First(&adminPerpustakaan).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Admin Perpustakaan tidak ditemukan", err.Error())
		return
	}

	var revisi models.Revisi
	if err := h.PerpustakaanService.DB.First(&revisi, uint(revisiID)).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Revisi tidak ditemukan", err.Error())
		return
	}

	revisi.ResponAdminPerpustakaan = req.Respon
	now := time.Now()
	revisi.TanggalResponAdminPerpustakaan = &now
	revisi.AdminPerpustakaanID = &adminPerpustakaan.ID
	revisi.StatusRevisi = "Selesai" // Atau "Ditolak" jika responnya negatif

	if err := h.PerpustakaanService.DB.Save(&revisi).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal merespon revisi", err.Error())
		return
	}

	// Opsional: Ubah status verifikasi perpustakaan kembali ke "Dalam Proses Verifikasi"
	// agar Admin DPK bisa meninjau ulang setelah revisi
	var perpustakaan models.Perpustakaan
	if err := h.PerpustakaanService.DB.First(&perpustakaan, revisi.IDData).Error; err == nil {
		perpustakaan.StatusVerifikasi = "Dalam Proses Verifikasi"
		h.PerpustakaanService.DB.Save(&perpustakaan)
	}

	utils.SuccessResponse(c, http.StatusOK, "Respon revisi berhasil dikirim", nil)
}

// UploadLaporan mengunggah file laporan
func (h *AdminDPKHandler) UploadLaporan(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal mendapatkan file", err.Error())
		return
	}

	periode := c.PostForm("periode")
	jenisLaporan := c.PostForm("jenis_laporan")

	if periode == "" || jenisLaporan == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Periode dan jenis laporan harus diisi", nil)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusInternalServerError, "ID pengguna tidak ditemukan di konteks", nil)
		return
	}

	laporan, err := h.AdminDPKService.UploadLaporan(file, periode, jenisLaporan, userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengunggah laporan", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Laporan berhasil diunggah", laporan)
}

// GetAllLaporan mengambil semua data laporan
func (h *AdminDPKHandler) GetAllLaporan(c *gin.Context) {
	laporans, err := h.AdminDPKService.GetAllLaporan()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil laporan", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Laporan berhasil diambil", laporans)
}

// DownloadLaporan mengunduh file laporan
func (h *AdminDPKHandler) DownloadLaporan(c *gin.Context) {
	idStr := c.Param("id")
	laporanID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID laporan tidak valid", nil)
		return
	}

	var laporan models.Laporan
	if err := h.AdminDPKService.DB.First(&laporan, uint(laporanID)).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Laporan tidak ditemukan", err.Error())
		return
	}

	c.FileAttachment(laporan.PathFile, laporan.NamaFile)
}

// ExportDataPerpustakaan mengekspor data perpustakaan
func (h *AdminDPKHandler) ExportDataPerpustakaan(c *gin.Context) {
	data, err := h.AdminDPKService.ExportDataPerpustakaan()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengekspor data perpustakaan", err.Error())
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=data_perpustakaan.csv")
	c.Data(http.StatusOK, "text/csv", data)
}

// GetVerifikasiHistory mengambil riwayat verifikasi
func (h *AdminDPKHandler) GetVerifikasiHistory(c *gin.Context) {
	perpustakaanIDStr := c.Param("perpustakaan_id")
	perpustakaanID, err := strconv.ParseUint(perpustakaanIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	history, err := h.AdminDPKService.GetVerifikasiHistory(uint(perpustakaanID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil riwayat verifikasi", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Riwayat verifikasi berhasil diambil", history)
}

// GetRevisiHistory mengambil riwayat revisi
func (h *AdminDPKHandler) GetRevisiHistory(c *gin.Context) {
	perpustakaanIDStr := c.Param("perpustakaan_id")
	perpustakaanID, err := strconv.ParseUint(perpustakaanIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	history, err := h.AdminDPKService.GetRevisiHistory(uint(perpustakaanID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil riwayat revisi", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Riwayat revisi berhasil diambil", history)
}

// SendNotification (Endpoint dummy untuk notifikasi)
func (h *AdminDPKHandler) SendNotification(c *gin.Context) {
	var req struct {
		TargetRole string `json:"target_role" binding:"required,oneof=AdminPerpustakaan AdminDPK Executive"`
		Message    string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	// Logika pengiriman notifikasi (misal: ke database, message queue, atau email)
	// Untuk demo, hanya log ke konsol
	fmt.Printf("Notifikasi dikirim ke %s: %s\n", req.TargetRole, req.Message)

	utils.SuccessResponse(c, http.StatusOK, "Notifikasi berhasil dikirim (dummy)", nil)
}
