package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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
	payload := msg.Message.Payload // Получаем payload

	log.Println("Получено сообщение:", text)
	log.Println("Получен payload:", payload)

	// ✅ Если пришел payload (нажатие кнопки), вызываем `handleButtonClick`
	if payload != "" {
		handleButtonClick(userID, payload)
		return
	}

	// Проверяем, если введено число - это выбор класса или расы
	step := getRegistrationStep(userID)

	if step == "choosing_class" {
		if classID, err := strconv.Atoi(text); err == nil && classExists(classID) {
			setUserClass(userID, classID)
			sendMessage(userID, "Класс выбран! Теперь выберите расу.\n\n"+getRaceListText())
			setRegistrationStep(userID, "choosing_race")
			return
		} else {
			sendMessage(userID, "Ошибка: такого класса нет. Введите число из списка.")
			return
		}
	}

	if step == "choosing_race" {
		if raceID, err := strconv.Atoi(text); err == nil && raceExists(raceID) {
			setUserRace(userID, raceID)
			finalizeRegistration(userID)
			sendMessage(userID, "Вы успешно зарегистрированы! Добро пожаловать в игру!")

			// Завершаем регистрацию, очищаем шаг
			setRegistrationStep(userID, "")
			return
		} else {
			sendMessage(userID, "Ошибка: такой расы нет. Введите число из списка.")
			return
		}
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

	// Пробуем распарсить JSON payload
	var cleanPayload string
	if err := json.Unmarshal([]byte(payload), &cleanPayload); err == nil {
		payload = cleanPayload
	}

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

		// ✅ Теперь функция используется
		registerUser(userID)

		sendMessage(userID, "Вы успешно зарегистрированы! Теперь выберите класс.\n\n"+getClassListText())

		// Устанавливаем этап регистрации
		setRegistrationStep(userID, "choosing_class")
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
