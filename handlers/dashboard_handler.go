package handlers

import (
	"net/http"
	"perpustakaan-api/services"
	"perpustakaan-api/utils"

	"github.com/gin-gonic/gin"
)

// DashboardHandler menangani permintaan terkait dashboard
type DashboardHandler struct {
	DashboardService *services.DashboardService
}

// NewDashboardHandler membuat instance DashboardHandler baru
func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{DashboardService: dashboardService}
}

// GetDashboardStats mengambil statistik untuk dashboard
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.DashboardService.GetDashboardStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil statistik dashboard", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Statistik dashboard berhasil diambil", stats)
}
