package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

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
		return fmt.Errorf("Не удалось подключиться к базе")
	}

	Conn.DB = db

	if err = db.Ping(); err != nil {
		return fmt.Errorf("Ошибка пинг-базы")
	}

	log.Printf("Успешное подключение к базе %s!", conf.Database)

	if err := startMigrations(db); err != nil {
		return err
	}

	return nil
}

func startMigrations(db *sql.DB) error {
	// чтение содержимого файла .sql
	sqlFile, err := ioutil.ReadFile("migrations/create_tables.up.sql")

	if err != nil {
		return fmt.Errorf("Ошибка при запуске миграций")
	}

	// Разделение содержимого файла на отдельные SQL-запросы
	sqlCommands := strings.Split(string(sqlFile), ";")

	// Выполнение каждого SQL-запроса
	for _, cmd := range sqlCommands {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" {
			_, err := db.Exec(cmd)
			if err != nil {
				log.Printf("Ошибка при выполнении SQL: %s\n%s\n", err, cmd)
			}
		}
	}

	log.Println("Таблицы успешно созданы!")

	return nil
}
