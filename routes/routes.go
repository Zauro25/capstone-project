package routes

import (
	"perpustakaan-api/handlers"
	"perpustakaan-api/middlewares"
	"perpustakaan-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Inisialisasi Services dan Handlers (tetap sama)
	authService := services.NewAuthService(db)
	perpustakaanService := services.NewPerpustakaanService(db)
	adminDPKService := services.NewAdminDPKService(db)
	dashboardService := services.NewDashboardService(db)

	authHandler := handlers.NewAuthHandler(authService)
	perpustakaanHandler := handlers.NewPerpustakaanHandler(perpustakaanService)
	adminDPKHandler := handlers.NewAdminDPKHandler(adminDPKService, perpustakaanService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Middleware global (tetap sama)
	router.Use(middlewares.RateLimitMiddleware())
	router.Use(middlewares.XSSProtectionMiddleware())
	router.Use(middlewares.CSRFProtectionMiddleware())
	router.Use(middlewares.AuditTrailMiddleware())

	v1 := router.Group("/api/v1")
	{
		// Rute Autentikasi (tetap sama)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterUser)
			auth.POST("/register/admin-perpustakaan", authHandler.RegisterAdminPerpustakaan)
    		auth.POST("/register/admin-dpk", authHandler.RegisterAdminDPK)
    		auth.POST("/register/executive", authHandler.RegisterExecutive)
			auth.POST("/login", authHandler.LoginUser)
			auth.POST("/logout", middlewares.AuthMiddleware(), authHandler.LogoutUser)
		}

		// Rute untuk Admin Perpustakaan - PERBAIKAN DI SINI
		adminPerpustakaan := v1.Group("/admin-perpustakaan")
		adminPerpustakaan.Use(middlewares.AuthMiddleware(), middlewares.RBACMiddleware("AdminPerpustakaan"))
		{
			// Group khusus untuk operasi perpustakaan
			perpustakaanGroup := adminPerpustakaan.Group("/perpustakaan")
			{
				perpustakaanGroup.POST("", perpustakaanHandler.CreatePerpustakaan)
				perpustakaanGroup.GET("/:id", perpustakaanHandler.GetPerpustakaanByID)
				perpustakaanGroup.PUT("/:id", perpustakaanHandler.UpdatePerpustakaan)
				perpustakaanGroup.DELETE("/:id", perpustakaanHandler.DeletePerpustakaan)
				perpustakaanGroup.POST("/:id/submit-verifikasi", perpustakaanHandler.SubmitDataForVerification)
				
				// Sub-group untuk SDM, Pengunjung, dan Anggota
				sdmGroup := perpustakaanGroup.Group("/:id/sdm")
				{
					sdmGroup.POST("", perpustakaanHandler.AddSDM)
					sdmGroup.GET("", perpustakaanHandler.GetSDMByPerpustakaanID)
					sdmGroup.PUT("/:sdm_id", perpustakaanHandler.UpdateSDM)
					sdmGroup.DELETE("/:sdm_id", perpustakaanHandler.DeleteSDM)
				}

				pengunjungGroup := perpustakaanGroup.Group("/:id/pengunjung")
				{
					pengunjungGroup.POST("", perpustakaanHandler.AddPengunjung)
					pengunjungGroup.GET("", perpustakaanHandler.GetPengunjungByPerpustakaanID)
					pengunjungGroup.PUT("/:pengunjung_id", perpustakaanHandler.UpdatePengunjung)
					pengunjungGroup.DELETE("/:pengunjung_id", perpustakaanHandler.DeletePengunjung)
				}

				anggotaGroup := perpustakaanGroup.Group("/:id/anggota")
				{
					anggotaGroup.POST("", perpustakaanHandler.AddAnggota)
					anggotaGroup.GET("", perpustakaanHandler.GetAnggotaByPerpustakaanID)
					anggotaGroup.PUT("/:anggota_id", perpustakaanHandler.UpdateAnggota)
					anggotaGroup.DELETE("/:anggota_id", perpustakaanHandler.DeleteAnggota)
				}
			}

			// Respon Revisi (dipindahkan ke luar group perpustakaan)
			adminPerpustakaan.POST("/revisi/:revisi_id/respon", adminDPKHandler.RespondToRevisi)
		}

		// Rute untuk Admin DPK - TETAP SAMA
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

			// Notifikasi
			adminDPK.POST("/notifikasi/send", adminDPKHandler.SendNotification)
		}

		// Rute untuk Executive - TETAP SAMA
		executive := v1.Group("/executive")
		executive.Use(middlewares.AuthMiddleware(), middlewares.RBACMiddleware("Executive", "AdminDPK"))
		{
			executive.GET("/dashboard-stats", dashboardHandler.GetDashboardStats)
			executive.GET("/perpustakaan", adminDPKHandler.GetAllPerpustakaanForDPK)
			executive.GET("/perpustakaan/:id", perpustakaanHandler.GetPerpustakaanByID)
			executive.GET("/laporan", adminDPKHandler.GetAllLaporan)
		}
	}
}