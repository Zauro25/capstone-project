package main

import (
	"log"
	"os"
	"perpustakaan-api/config"
	"perpustakaan-api/models"
	"perpustakaan-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi database
	db := config.InitDB()

	// Auto-migrate semua model
	err := db.AutoMigrate(
		&models.Perpustakaan{},
		&models.SDM{},
		&models.Pengunjung{},
		&models.Anggota{},
		&models.AdminPerpustakaan{},
		&models.AdminDPK{},
		&models.Verifikasi{},
		&models.Revisi{},
		&models.Laporan{},
		&models.LogAktivitas{},
		&models.User{},
	)
	if err != nil {
		log.Fatalf("Gagal auto-migrate database: %v", err)
	}
	log.Println("Database berhasil di-migrate.")

	// Inisialisasi Gin router
	router := gin.Default()

	// Setup CORS sederhana
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(router, db)

	// Ambil port dari environment variable sistem (bukan dari .env)
	port := os.Getenv("PORT") // Heroku/DigitalOcean biasanya pakai PORT
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server berjalan di port %s", port)
	log.Fatal(router.Run(":" + port))
}