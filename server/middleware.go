package server

import (
	"fmt"
	"net/http"

	"github.com/b-o-e-v/crypto-up-back/pkg/auth"
	"github.com/b-o-e-v/crypto-up-back/pkg/envs"
	"github.com/gin-gonic/gin"
)

var allowedPaths = [3]string{
	"/ping",
	"/signin",
	"/signup",
}

func isAllowedPath(target string) bool {
	for _, str := range allowedPaths {
		if str == target {
			return true
		}
	}

	return false
}

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

func AUTHMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Получение пути запроса
		path := ctx.Request.URL.Path
		var token = getToken(ctx)

		claims, err := auth.ValidateToken(token, envs.Conf.SecretKey)

		if err != nil {
			if isAllowed := isAllowedPath(path); isAllowed == false {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Необходима авторизации")})
				ctx.Abort()
				return
			} else {
				claims = &auth.SignedDetails{Id: ""}
			}
		}

		ctx.Set("id", claims.Id)
		ctx.Next()
	}
}
