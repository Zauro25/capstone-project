package middlewares

import (
	"net/http"
	"strings"

	"perpustakaan-api/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memverifikasi token JWT dari header Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token autentikasi tidak disediakan", nil)
			c.Abort()
			return
		}

		// Format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Format token tidak valid", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid atau kadaluarsa", err.Error())
			c.Abort()
			return
		}

		// Simpan klaim di konteks untuk digunakan di handler selanjutnya
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("user_id", claims.Subject) // Jika Anda menyimpan ID pengguna di Subject JWT

		c.Next()
	}
}
