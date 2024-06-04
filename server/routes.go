package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Up(port string) error {
	server := gin.Default()
	// устанавливаем заголовки
	server.Use(CORSMiddleware())
	// добавляем проверку на авторизацию
	server.Use(AUTHMiddleware())

	// проверяем работу сервера
	server.GET("/ping", ping)
	// регистрация
	server.POST("/signup", signup)
	// авторизация
	server.POST("/signin", signin)

	// запускаем сервер
	if err := server.Run(port); err != nil {
		fmt.Println("Не удалось запустить сервер")
		return err
	}

	return nil
}
