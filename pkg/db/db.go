package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	User     string
	Password string
	Database string
	Port     string
	Host     string
	SSLMode  string
}

type Connect struct {
	DB *sql.DB
}

// соединение с БД
var Conn = &Connect{}

func Up(conf Config) error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Database,
		conf.SSLMode,
	))

	if err != nil {
		fmt.Println("Не удалось подключиться к базе")
		return err
	}

	Conn.DB = db

	if err = db.Ping(); err != nil {
		fmt.Println("Ошибка пинг-базы")
		return err
	}

	log.Printf("Успешное подключение к базе %s", conf.Database)

	return nil
}
