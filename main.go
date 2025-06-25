package main

import (
	"log"
	"os"
	"perpustakaan-api/config"
	"perpustakaan-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi database
	config.InitDB() // Call InitDB without assignment if it returns no value



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
	routes.SetupRoutes(router, config.DB)

	// Ambil port dari environment variable sistem (bukan dari .env)
	port := os.Getenv("PORT") // Heroku/DigitalOcean biasanya pakai PORT
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server berjalan di port %s", port)
	log.Fatal(router.Run(":" + port))
}