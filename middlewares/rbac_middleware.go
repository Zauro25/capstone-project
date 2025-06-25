package middlewares

import (
	"net/http"
	"perpustakaan-api/utils"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware memeriksa apakah pengguna memiliki peran yang diizinkan
func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Informasi peran pengguna tidak ditemukan", nil)
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Tipe peran pengguna tidak valid", nil)
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki izin untuk mengakses sumber daya ini", nil)
		c.Abort()
	}
}
