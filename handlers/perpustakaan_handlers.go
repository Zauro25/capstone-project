package handlers

import (
	"net/http"
	"perpustakaan-api/models"
	"perpustakaan-api/services"
	"perpustakaan-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PerpustakaanHandler menangani permintaan terkait data perpustakaan
type PerpustakaanHandler struct {
	PerpustakaanService *services.PerpustakaanService
}

// NewPerpustakaanHandler membuat instance PerpustakaanHandler baru
func NewPerpustakaanHandler(perpustakaanService *services.PerpustakaanService) *PerpustakaanHandler {
	return &PerpustakaanHandler{PerpustakaanService: perpustakaanService}
}

// CreatePerpustakaan membuat entri perpustakaan baru
func (h *PerpustakaanHandler) CreatePerpustakaan(c *gin.Context) {
	var perpustakaan models.Perpustakaan
	if err := c.ShouldBindJSON(&perpustakaan); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	if err := h.PerpustakaanService.CreatePerpustakaan(&perpustakaan); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat perpustakaan", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Perpustakaan berhasil dibuat", perpustakaan)
}

// GetAllPerpustakaan mengambil semua data perpustakaan
func (h *PerpustakaanHandler) GetAllPerpustakaan(c *gin.Context) {
	perpustakaans, err := h.PerpustakaanService.GetAllPerpustakaan()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data perpustakaan", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data perpustakaan berhasil diambil", perpustakaans)
}

// GetPerpustakaanByID mengambil data perpustakaan berdasarkan ID
func (h *PerpustakaanHandler) GetPerpustakaanByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	perpustakaan, err := h.PerpustakaanService.GetPerpustakaanByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Perpustakaan tidak ditemukan", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Perpustakaan berhasil diambil", perpustakaan)
}

// UpdatePerpustakaan memperbarui data perpustakaan
func (h *PerpustakaanHandler) UpdatePerpustakaan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	var perpustakaan models.Perpustakaan
	if err := c.ShouldBindJSON(&perpustakaan); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	if err := h.PerpustakaanService.UpdatePerpustakaan(uint(id), &perpustakaan); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui perpustakaan", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Perpustakaan berhasil diperbarui", nil)
}

// DeletePerpustakaan menghapus data perpustakaan
func (h *PerpustakaanHandler) DeletePerpustakaan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	if err := h.PerpustakaanService.DeletePerpustakaan(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus perpustakaan", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Perpustakaan berhasil dihapus", nil)
}

// SubmitDataForVerification mengajukan data perpustakaan untuk diverifikasi
func (h *PerpustakaanHandler) SubmitDataForVerification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	if err := h.PerpustakaanService.SubmitDataForVerification(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengajukan data untuk verifikasi", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data perpustakaan berhasil diajukan untuk verifikasi", nil)
}

// --- Handler untuk sub-entitas (Koleksi, SDM, Pengunjung, Anggota) ---

// AddSDM menambahkan SDM ke perpustakaan
func (h *PerpustakaanHandler) AddSDM(c *gin.Context) {
	var sdm models.SDM
	if err := c.ShouldBindJSON(&sdm); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	if err := h.PerpustakaanService.AddSDM(&sdm); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan SDM", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "SDM berhasil ditambahkan", sdm)
}

// GetSDMByPerpustakaanID mengambil SDM berdasarkan ID perpustakaan
func (h *PerpustakaanHandler) GetSDMByPerpustakaanID(c *gin.Context) {
	perpustakaanIDStr := c.Param("perpustakaan_id")
	perpustakaanID, err := strconv.ParseUint(perpustakaanIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	sdms, err := h.PerpustakaanService.GetSDMByPerpustakaanID(uint(perpustakaanID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil SDM", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "SDM berhasil diambil", sdms)
}

// UpdateSDM memperbarui data SDM
func (h *PerpustakaanHandler) UpdateSDM(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID SDM tidak valid", nil)
		return
	}

	var sdm models.SDM
	if err := c.ShouldBindJSON(&sdm); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	if err := h.PerpustakaanService.UpdateSDM(uint(id), &sdm); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui SDM", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "SDM berhasil diperbarui", nil)
}

// DeleteSDM menghapus data SDM
func (h *PerpustakaanHandler) DeleteSDM(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID SDM tidak valid", nil)
		return
	}

	if err := h.PerpustakaanService.DeleteSDM(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus SDM", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "SDM berhasil dihapus", nil)
}

// AddPengunjung menambahkan data pengunjung ke perpustakaan
func (h *PerpustakaanHandler) AddPengunjung(c *gin.Context) {
	var pengunjung models.Pengunjung
	if err := c.ShouldBindJSON(&pengunjung); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	if err := h.PerpustakaanService.AddPengunjung(&pengunjung); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan data pengunjung", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Data pengunjung berhasil ditambahkan", pengunjung)
}

// GetPengunjungByPerpustakaanID mengambil data pengunjung berdasarkan ID perpustakaan
func (h *PerpustakaanHandler) GetPengunjungByPerpustakaanID(c *gin.Context) {
	perpustakaanIDStr := c.Param("perpustakaan_id")
	perpustakaanID, err := strconv.ParseUint(perpustakaanIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	pengunjung, err := h.PerpustakaanService.GetPengunjungByPerpustakaanID(uint(perpustakaanID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data pengunjung", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data pengunjung berhasil diambil", pengunjung)
}

// UpdatePengunjung memperbarui data pengunjung
func (h *PerpustakaanHandler) UpdatePengunjung(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID pengunjung tidak valid", nil)
		return
	}

	var pengunjung models.Pengunjung
	if err := c.ShouldBindJSON(&pengunjung); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	if err := h.PerpustakaanService.UpdatePengunjung(uint(id), &pengunjung); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data pengunjung", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data pengunjung berhasil diperbarui", nil)
}

// DeletePengunjung menghapus data pengunjung
func (h *PerpustakaanHandler) DeletePengunjung(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID pengunjung tidak valid", nil)
		return
	}

	if err := h.PerpustakaanService.DeletePengunjung(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus data pengunjung", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data pengunjung berhasil dihapus", nil)
}

// AddAnggota menambahkan data anggota ke perpustakaan
func (h *PerpustakaanHandler) AddAnggota(c *gin.Context) {
	var anggota models.Anggota
	if err := c.ShouldBindJSON(&anggota); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	if err := h.PerpustakaanService.AddAnggota(&anggota); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan data anggota", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Data anggota berhasil ditambahkan", anggota)
}

// GetAnggotaByPerpustakaanID mengambil data anggota berdasarkan ID perpustakaan
func (h *PerpustakaanHandler) GetAnggotaByPerpustakaanID(c *gin.Context) {
	perpustakaanIDStr := c.Param("perpustakaan_id")
	perpustakaanID, err := strconv.ParseUint(perpustakaanIDStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID perpustakaan tidak valid", nil)
		return
	}

	anggota, err := h.PerpustakaanService.GetAnggotaByPerpustakaanID(uint(perpustakaanID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data anggota", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data anggota berhasil diambil", anggota)
}

// UpdateAnggota memperbarui data anggota
func (h *PerpustakaanHandler) UpdateAnggota(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID anggota tidak valid", nil)
		return
	}

	var anggota models.Anggota
	if err := c.ShouldBindJSON(&anggota); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	if err := h.PerpustakaanService.UpdateAnggota(uint(id), &anggota); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data anggota", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data anggota berhasil diperbarui", nil)
}

// DeleteAnggota menghapus data anggota
func (h *PerpustakaanHandler) DeleteAnggota(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID anggota tidak valid", nil)
		return
	}

	if err := h.PerpustakaanService.DeleteAnggota(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus data anggota", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data anggota berhasil dihapus", nil)
}
