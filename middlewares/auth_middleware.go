package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matheusmpz/auth-api/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token de autenticação necessário",
			})
			return
		}

		// Espera formato: Bearer TOKEN
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido",
			})
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido ou expirado",
			})
			return
		}

		// Salva userID no contexto
		c.Set("userID", claims.UserID)

		c.Next()
	}
}