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

func isUserRegistered(vkID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM profiles WHERE vk_user_id = $1)", vkID).Scan(&exists)
	if err != nil {
		log.Println("Ошибка проверки регистрации:", err)
		return false
	}
	return exists
}

func registerUser(vkID int) {
	_, err := db.Exec("INSERT INTO profiles (vk_user_id) VALUES ($1) ON CONFLICT (vk_user_id) DO NOTHING", vkID)
	if err != nil {
		log.Println("Ошибка регистрации пользователя:", err)
	} else {
		log.Println("✅ Пользователь", vkID, "успешно зарегистрирован!")
	}
}

func getClasses() []Class {
	var classes []Class

	// Запрос в базу для получения классов
	rows, err := db.Query("SELECT id, name FROM classes")
	if err != nil {
		log.Println("Ошибка получения классов из БД:", err)
		return nil
	}
	defer rows.Close()

	// Считываем классы из результата запроса
	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.Name); err != nil {
			log.Println("Ошибка при разборе данных класса:", err)
			continue
		}
		classes = append(classes, class)
	}

	// Проверяем ошибки после завершения работы с результатом
	if err := rows.Err(); err != nil {
		log.Println("Ошибка при обработке списка классов:", err)
		return nil
	}

	return classes
}

// Структура класса
type Class struct {
	ID   int
	Name string
}

// func registerUser(vkID int, firstName, lastName string) {
// 	// Проверяем, зарегистрирован ли пользователь
// 	var count int
// 	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE vk_id = $1", vkID).Scan(&count)
// 	if err != nil {
// 		log.Println("Ошибка при проверке пользователя в базе данных:", err)
// 		return
// 	}

// 	if count == 0 {
// 		// Добавляем нового пользователя в базу данных
// 		_, err = db.Exec("INSERT INTO users (vk_id, first_name, last_name) VALUES ($1, $2, $3)", vkID, firstName, lastName)
// 		if err != nil {
// 			log.Println("Ошибка при добавлении пользователя в базу данных:", err)
// 			return
// 		}
// 		log.Printf("Пользователь %s %s успешно зарегистрирован.", firstName, lastName)
// 	}
// }
