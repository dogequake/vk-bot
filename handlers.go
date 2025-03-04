package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v3/events"
)

// CallbackHandler обрабатывает запросы от VK
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Ошибка чтения тела запроса:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	log.Println("Тело запроса от VK:", string(body))

	var req events.GroupEvent
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("Ошибка обработки JSON:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	switch req.Type {
	case events.EventConfirmation:
		log.Println("Отправлен confirmation_code:", confirmationCode)
		fmt.Fprint(w, confirmationCode)
		return
	case events.EventMessageNew:
		var msg events.MessageNewObject
		if err := json.Unmarshal(req.Object, &msg); err != nil {
			log.Println("Ошибка декодирования сообщения:", err)
			return
		}
		handleMessage(msg)
	}

	fmt.Fprint(w, "ok")
}

func handleMessage(msg events.MessageNewObject) {
	userID := msg.Message.PeerID
	text := msg.Message.Text

	log.Println("Получено сообщение:", text)

	// Проверяем, если введено число - это выбор класса или расы
	if classID, err := strconv.Atoi(text); err == nil {
		if classExists(classID) {
			setUserClass(userID, classID)
			sendMessage(userID, "Класс выбран! Теперь выберите расу.\n\n"+getRaceListText())
			return
		}
	}

	if raceID, err := strconv.Atoi(text); err == nil {
		if raceExists(raceID) {
			setUserRace(userID, raceID)
			finalizeRegistration(userID)
			sendMessage(userID, "Вы успешно зарегистрированы! Добро пожаловать в игру!")
		} else {
			sendMessage(userID, "Ошибка: такой расы нет. Введите число от 1 до 10.")
		}
		return
	}

	// Обработка остальных команд
	switch text {
	case "/start":
		if isUserRegistered(userID) {
			sendMessage(userID, "Вы уже зарегистрированы! Добро пожаловать обратно.")
		} else {
			sendRegistrationPrompt(userID)
		}
	default:
		sendMessage(userID, "Неизвестная команда. Используйте /start")
	}
}

func handleButtonClick(userID int, payload string) {
	log.Println("Получен payload (до обработки):", payload)

	// Убираем кавычки из payload, если они есть
	payload = strings.Trim(payload, "\"")

	log.Println("Получен payload (после обработки):", payload)

	switch payload {
	case "profile":
		sendMessage(userID, "Вот ваш профиль.")
	case "stats":
		sendMessage(userID, "Вот ваша статистика.")
	case "register":
		if isUserRegistered(userID) {
			sendMessage(userID, "Вы уже зарегистрированы.")
			return
		}

		registerUser(userID) // Пока просто добавляем в базу
		sendMessage(userID, "Вы успешно зарегистрированы! Теперь выберите класс.\n\n"+getClassListText())
	default:
		sendMessage(userID, "Неизвестная кнопка.")
	}
}

func getClassListText() string {
	classes := getClasses()

	if len(classes) == 0 {
		return "Ошибка: классы не найдены."
	}

	var classList string
	for _, class := range classes {
		classList += fmt.Sprintf("%d. %s\n", class.ID, class.Name)
	}

	return "Вот доступные классы:\n\n" + classList + "\nВведите номер класса, чтобы выбрать."
}
