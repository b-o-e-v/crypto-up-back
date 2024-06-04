package server

import (
	"fmt"
	"net/http"

	"github.com/b-o-e-v/crypto-up-back/pkg/auth"
	"github.com/b-o-e-v/crypto-up-back/pkg/db"
	"github.com/b-o-e-v/crypto-up-back/pkg/envs"
	"github.com/gin-gonic/gin"
)

const updateUserTokenSQL = `
	UPDATE users SET token = ($1) WHERE id = ($2)
`

const insertUserSQL = `
	INSERT INTO users (email, phone, login, display_name, image_url, password)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
`

const getUserByIdSQL = `
	SELECT id, created_at, email, phone, login, display_name, image_url, password
	FROM users WHERE id = ($1)
`

const getUserByEmailSQL = `
	SELECT id, created_at, email, phone, login, display_name, image_url, password
	FROM users
	WHERE email LIKE ($1)
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

	db.Conn.DB.QueryRow(updateUserTokenSQL, token, userID)

	ctx.JSON(http.StatusOK, gin.H{"data": token})
}

// авторизация
func signin(ctx *gin.Context) {
	var user User
	var token = getToken(ctx)
	userId := ctx.GetString("id")

	if userId != "" {
		row := db.Conn.DB.QueryRow(getUserByIdSQL, userId)

		if err := getUserFromDB(row, &user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		var data User

		if err := ctx.BindJSON(&data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if data.Email == "" {
			ctx.JSON(http.StatusOK, gin.H{"error": "Необходимо указать email"})
			return
		}

		row := db.Conn.DB.QueryRow(getUserByEmailSQL, data.Email)

		if err := getUserFromDB(row, &user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		passwordIsValid, msg := auth.VerifyPassword(data.Password, user.Password)

		if passwordIsValid != true {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		newToken, err := auth.GenerateToken(fmt.Sprintf("%s", user.Id), envs.Conf.SecretKey)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		db.Conn.DB.QueryRow(updateUserTokenSQL, newToken, user.Id)
		token = newToken
	}

	ctx.JSON(http.StatusOK, gin.H{"data": User{
		Id:          user.Id,
		CreatedAt:   user.CreatedAt,
		Email:       user.Email,
		Phone:       user.Phone,
		Login:       user.Login,
		DisplayName: user.DisplayName,
		ImageUrl:    user.ImageUrl,
		Token:       token,
	}})
}
