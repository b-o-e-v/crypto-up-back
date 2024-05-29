package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// пинг сервера
func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// регистрация
func signup(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "signup"})
}

// авторизация
func signin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "signin"})
}
