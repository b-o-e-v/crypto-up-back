package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	Id string
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println("Ошибка при генерации хэша пароля")
		return "", err
	}

	return string(bytes), nil
}

func VerifyPassword(password string, hashedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Phone or Password is incorrect")
		check = false
	}

	return check, msg
}

func GenerateToken(id string, secretKey string) (string, error) {
	// Получаем текущее время с учетом часового пояса и добавляем сутки
	date := time.Now().Local().Add(time.Hour * time.Duration(24))

	claims := &SignedDetails{
		Id:               id,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{date}},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken проверяет действительность JWT токена.
// openssl rand -base64 256
func ValidateToken(signedToken string, secretKey string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, fmt.Errorf("The token is invalid")
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return nil, fmt.Errorf("Token is expired")
	}

	return claims, nil
}
