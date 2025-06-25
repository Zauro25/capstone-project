package routes

import (
	"perpustakaan-api/handlers"
	"perpustakaan-api/middlewares"
	"perpustakaan-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes menginisialisasi semua rute API
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Inisialisasi Services
	authService := services.NewAuthService(db)
	perpustakaanService := services.NewPerpustakaanService(db)
	adminDPKService := services.NewAdminDPKService(db)
	dashboardService := services.NewDashboardService(db)

	// Inisialisasi Handlers
	authHandler := handlers.NewAuthHandler(authService)
	perpustakaanHandler := handlers.NewPerpustakaanHandler(perpustakaanService)
	adminDPKHandler := handlers.NewAdminDPKHandler(adminDPKService, perpustakaanService) // Pass perpustakaanService for revisi response
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Middleware global
	router.Use(middlewares.RateLimitMiddleware())
	router.Use(middlewares.XSSProtectionMiddleware())
	router.Use(middlewares.CSRFProtectionMiddleware()) // Placeholder, lihat catatan di middleware
	router.Use(middlewares.AuditTrailMiddleware())

	v1 := router.Group("/api/v1")
	{
		// Rute Autentikasi (tanpa autentikasi JWT)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterUser)
			auth.POST("/login", authHandler.LoginUser)
			auth.POST("/logout", middlewares.AuthMiddleware(), authHandler.LogoutUser) // Logout memerlukan token untuk identifikasi user
		}

		// Rute untuk Admin Perpustakaan
		adminPerpustakaan := v1.Group("/admin-perpustakaan")
		adminPerpustakaan.Use(middlewares.AuthMiddleware(), middlewares.RBACMiddleware("AdminPerpustakaan"))
		{
			// Perpustakaan
			adminPerpustakaan.POST("/perpustakaan", perpustakaanHandler.CreatePerpustakaan)
			adminPerpustakaan.GET("/perpustakaan/:id", perpustakaanHandler.GetPerpustakaanByID)
			adminPerpustakaan.PUT("/perpustakaan/:id", perpustakaanHandler.UpdatePerpustakaan)
			adminPerpustakaan.DELETE("/perpustakaan/:id", perpustakaanHandler.DeletePerpustakaan)
			adminPerpustakaan.POST("/perpustakaan/:id/submit-verifikasi", perpustakaanHandler.SubmitDataForVerification)

			// SDM
			adminPerpustakaan.POST("/perpustakaan/:perpustakaan_id/sdm", perpustakaanHandler.AddSDM)
			adminPerpustakaan.GET("/perpustakaan/:perpustakaan_id/sdm", perpustakaanHandler.GetSDMByPerpustakaanID)
			adminPerpustakaan.PUT("/sdm/:id", perpustakaanHandler.UpdateSDM)
			adminPerpustakaan.DELETE("/sdm/:id", perpustakaanHandler.DeleteSDM)

			// Pengunjung
			adminPerpustakaan.POST("/perpustakaan/:perpustakaan_id/pengunjung", perpustakaanHandler.AddPengunjung)
			adminPerpustakaan.GET("/perpustakaan/:perpustakaan_id/pengunjung", perpustakaanHandler.GetPengunjungByPerpustakaanID)
			adminPerpustakaan.PUT("/pengunjung/:id", perpustakaanHandler.UpdatePengunjung)
			adminPerpustakaan.DELETE("/pengunjung/:id", perpustakaanHandler.DeletePengunjung)

			// Anggota
			adminPerpustakaan.POST("/perpustakaan/:perpustakaan_id/anggota", perpustakaanHandler.AddAnggota)
			adminPerpustakaan.GET("/perpustakaan/:perpustakaan_id/anggota", perpustakaanHandler.GetAnggotaByPerpustakaanID)
			adminPerpustakaan.PUT("/anggota/:id", perpustakaanHandler.UpdateAnggota)
			adminPerpustakaan.DELETE("/anggota/:id", perpustakaanHandler.DeleteAnggota)

			// Respon Revisi (Admin Perpustakaan merespon revisi dari DPK)
			adminPerpustakaan.POST("/revisi/:revisi_id/respon", adminDPKHandler.RespondToRevisi)
		}

		// Rute untuk Admin DPK
		adminDPK := v1.Group("/admin-dpk")
		adminDPK.Use(middlewares.AuthMiddleware(), middlewares.RBACMiddleware("AdminDPK"))
		{
			adminDPK.GET("/perpustakaan", adminDPKHandler.GetAllPerpustakaanForDPK)
			adminDPK.POST("/perpustakaan/:id/verify", adminDPKHandler.VerifyPerpustakaan)
			adminDPK.POST("/perpustakaan/:id/request-revisi", adminDPKHandler.RequestRevisi)

			// Laporan
			adminDPK.POST("/laporan/upload", adminDPKHandler.UploadLaporan)
			adminDPK.GET("/laporan", adminDPKHandler.GetAllLaporan)
			adminDPK.GET("/laporan/:id/download", adminDPKHandler.DownloadLaporan)
			adminDPK.GET("/perpustakaan/export-csv", adminDPKHandler.ExportDataPerpustakaan)

			// Riwayat Verifikasi/Revisi
			adminDPK.GET("/perpustakaan/:perpustakaan_id/verifikasi-history", adminDPKHandler.GetVerifikasiHistory)
			adminDPK.GET("/perpustakaan/:perpustakaan_id/revisi-history", adminDPKHandler.GetRevisiHistory)

			// Notifikasi (dummy)
			adminDPK.POST("/notifikasi/send", adminDPKHandler.SendNotification)
		}

		// Rute untuk Executive (read-only)
		executive := v1.Group("/executive")
		executive.Use(middlewares.AuthMiddleware(), middlewares.RBACMiddleware("Executive", "AdminDPK")) // AdminDPK juga bisa melihat dashboard
		{
			executive.GET("/dashboard-stats", dashboardHandler.GetDashboardStats)
			executive.GET("/perpustakaan", adminDPKHandler.GetAllPerpustakaanForDPK) // Executive bisa melihat semua perpustakaan
			executive.GET("/perpustakaan/:id", perpustakaanHandler.GetPerpustakaanByID) // Executive bisa melihat detail perpustakaan
			executive.GET("/laporan", adminDPKHandler.GetAllLaporan) // Executive bisa melihat laporan
		}
	}
}
