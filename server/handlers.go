package server

import (
	"fmt"
	"net/http"

	"github.com/b-o-e-v/crypto-up-back/pkg/auth"
	"github.com/b-o-e-v/crypto-up-back/pkg/db"
	"github.com/b-o-e-v/crypto-up-back/pkg/envs"
	"github.com/gin-gonic/gin"
)

const insertUserSQL = `
	INSERT INTO users (email, phone, login, display_name, image_url, password)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING user_id
`

// пинг сервера
func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// регистрация
func signup(ctx *gin.Context) {
	var user = &User{}

	// создаем структуру
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обработать регистрационные данные"})
		return
	}

	// проверяем передан ли пароль
	if user.Password == "" {
		ctx.JSON(http.StatusOK, gin.H{"error": "Пароль не передан"})
		return
	}

	// генерируем хэш
	hashPassword, err := auth.HashPassword(user.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// записываем в структуру хэшированный пароль
	user.Password = hashPassword

	// Выполнение SQL-запроса
	var userID int

	err = db.Conn.DB.QueryRow(
		insertUserSQL,
		user.Email,
		user.Phone,
		user.Login,
		user.DisplayName,
		user.ImageUrl,
		user.Password,
	).Scan(&userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(fmt.Sprintf("%d", userID), envs.Conf.SecretKey)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": token})
}

// авторизация
func signin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "signin"})
}
