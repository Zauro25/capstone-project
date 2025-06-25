package middlewares

import (
	"net/http"
	"time"

	"perpustakaan-api/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Simple rate limiting map (in-memory for simplicity)
var ipRequestCounts = make(map[string]int)
var ipLastRequestTime = make(map[string]time.Time)
const maxRequests = 5 // Max requests per minute
const window = 1 * time.Minute

// RateLimitMiddleware membatasi jumlah permintaan dari IP yang sama
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if lastTime, ok := ipLastRequestTime[ip]; ok {
			if time.Since(lastTime) < window {
				ipRequestCounts[ip]++
				if ipRequestCounts[ip] > maxRequests {
					utils.ErrorResponse(c, http.StatusTooManyRequests, "Terlalu banyak permintaan. Coba lagi nanti.", nil)
					c.Abort()
					return
				}
			} else {
				// Reset count if window passed
				ipRequestCounts[ip] = 1
				ipLastRequestTime[ip] = time.Now()
			}
		} else {
			ipRequestCounts[ip] = 1
			ipLastRequestTime[ip] = time.Now()
		}
		c.Next()
	}
}

// AuditTrailMiddleware mencatat aktivitas pengguna
func AuditTrailMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sebelum request diproses
		c.Next()

		// Setelah request diproses
		username, _ := c.Get("username")
		role, _ := c.Get("role")
		userID, _ := c.Get("user_id") // Asumsi user_id disimpan di JWT Subject

		// Contoh sederhana: log setiap request yang berhasil
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			// Di sini Anda bisa memanggil service untuk menyimpan log ke database
			// Untuk demo, kita hanya print ke konsol
			logMessage := fmt.Sprintf("Aktivitas: %s %s, User: %v (Role: %v, ID: %v), IP: %s, Status: %d",
				c.Request.Method, c.Request.URL.Path, username, role, userID, c.ClientIP(), c.Writer.Status())
			fmt.Println(logMessage)

			// TODO: Simpan ke models.LogAktivitas
			// logAktivitas := models.LogAktivitas{
			// 	UserID:        userID.(uint), // Pastikan tipe data sesuai
			// 	TipeAktivitas: c.Request.Method + " " + c.Request.URL.Path,
			// 	Deskripsi:     logMessage,
			// 	Timestamp:     time.Now(),
			// 	IPAddress:     c.ClientIP(),
			// }
			// db.Create(&logAktivitas) // Perlu akses ke instance DB
		}
	}
}

// XSSProtectionMiddleware menambahkan header X-XSS-Protection
func XSSProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}

// CSRFProtectionMiddleware (Placeholder - CSRF biasanya ditangani di sisi frontend dengan token)
// Untuk API stateless dengan JWT, CSRF kurang relevan karena tidak ada session cookie.
// Namun, jika ada session cookie atau jika API digunakan dengan form HTML, ini penting.
func CSRFProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementasi CSRF token validation di sini jika diperlukan
		// Untuk API berbasis JWT, ini seringkali tidak diperlukan karena JWT disimpan di localStorage/header
		// dan tidak rentan terhadap serangan CSRF klasik.
		c.Next()
	}
}
