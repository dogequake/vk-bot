package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := getEnvVar("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
}

func registerUser(vkID int, firstName, lastName string) {
	// Проверяем, зарегистрирован ли пользователь
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE vk_id = $1", vkID).Scan(&count)
	if err != nil {
		log.Println("Ошибка при проверке пользователя в базе данных:", err)
		return
	}

	if count == 0 {
		// Добавляем нового пользователя в базу данных
		_, err = db.Exec("INSERT INTO users (vk_id, first_name, last_name) VALUES ($1, $2, $3)", vkID, firstName, lastName)
		if err != nil {
			log.Println("Ошибка при добавлении пользователя в базу данных:", err)
			return
		}
		log.Printf("Пользователь %s %s успешно зарегистрирован.", firstName, lastName)
	}
}
