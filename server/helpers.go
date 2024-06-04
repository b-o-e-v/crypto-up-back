package server

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func getToken(ctx *gin.Context) string {
	return strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
}

func getUserFromDB(row *sql.Row, user *User) error {
	if err := row.Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Email,
		&user.Phone,
		&user.Login,
		&user.DisplayName,
		&user.ImageUrl,
		&user.Password,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return fmt.Errorf("Пользователь не найден")
		default:
			return err
		}
	}
	return nil
}
