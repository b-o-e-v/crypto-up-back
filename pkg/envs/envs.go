package envs

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type EnvConfig struct {
	Port       string
	DBPort     string
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
	SecretKey  string
}

var Conf *EnvConfig

func LoadConf() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Ошибка загрузки переменных окружения из .env файла")
		return err
	}

	Conf = &EnvConfig{
		Port:       os.Getenv("APP_PORT"),
		DBPort:     os.Getenv("DB_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		SecretKey:  os.Getenv("SECRET_KEY"),
	}

	return nil
}
