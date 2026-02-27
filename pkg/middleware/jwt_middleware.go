package middleware

import (
	"hyperscal-go/internal/dto"
	"hyperscal-go/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse("Authorization header requires", "missing token"))
			ctx.Abort()
			return
		}

		// Extract token dari Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid Authorization", "format must be: Bearer <token>"))
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// validasi token
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid token", err.Error()))
			ctx.Abort()
			return
		}

		//Set user info ke context untuk digunakan handler
		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)

		//lanjut ke handler berikutnya
		ctx.Next()
	}
}
