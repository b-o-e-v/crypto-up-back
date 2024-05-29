package main

import (
	"fmt"
	"log"
	"os"

	"gitlab.com/dev.pologov/crypto-up-back/pkg/db"
	"gitlab.com/dev.pologov/crypto-up-back/pkg/envs"
	"gitlab.com/dev.pologov/crypto-up-back/server"
)

func init() {
	// получаем переменные окружения
	if err := envs.LoadConf(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	// конфиг для подключения к базе
	conf := db.Config{
		User:     envs.Conf.DBUser,
		Password: envs.Conf.DBPassword,
		Database: envs.Conf.DBName,
		Port:     envs.Conf.DBPort, // Порт по умолчанию для PostgreSQL
		Host:     envs.Conf.DBHost,
		SSLMode:  "disable", // Отключение SSL (временно, если не требуется)
	}

	// запускаем базу
	if err := db.Up(conf); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// запускаем сервер
	if err := server.Up(fmt.Sprintf(":%s", envs.Conf.Port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
