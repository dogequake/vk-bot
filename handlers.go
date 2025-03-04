package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/events"
	"github.com/SevereCloud/vksdk/v3/object"
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
	payload := msg.Message.Payload

	// Выводим текст сообщения для отладки
	log.Println("Получено сообщение:", text)

	// Если это нажатие на кнопку (Payload)
	if payload != "" {
		log.Println("Получен payload:", payload)
		handleButtonClick(userID, payload)
		return
	}

	// Обработка обычных сообщений
	switch text {
	case "/start", "\\/start":
		if !isUserRegistered(userID) {
			sendRegistrationPrompt(userID)
		} else {
			sendMessage(userID, "Вы уже зарегистрированы!")
		}
		return
		//sendMessageWithButtons(userID, "Добро пожаловать в игру! Выберите действие:")
	default:
		sendMessage(userID, "Неизвестная команда. Используйте /start")
	}
}

func handleButtonClick(userID int, payload string) {
	log.Println("Получен payload (до обработки):", payload)

	// Убираем лишние кавычки, если они есть
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

		registerUser(userID)
		sendMessage(userID, "Вы успешно зарегистрированы! Теперь выберите класс.")

		// Отправляем кнопки с классами
		sendClassChoice(userID)
	default:
		sendMessage(userID, "Неизвестная кнопка.")
	}
}

func sendClassChoice(userID int) {
	// Получаем список доступных классов из базы
	classes := getClasses()

	// Если классы не найдены, отправляем ошибку
	if len(classes) == 0 {
		sendMessage(userID, "Ошибка: не найдено ни одного класса.")
		return
	}

	// Создаем клавиатуру для выбора класса
	keyboard := object.NewMessagesKeyboardInline()

	// Добавляем кнопку для каждого класса
	for _, class := range classes {
		keyboard.AddRow().AddTextButton(class.Name, fmt.Sprintf("class_%d", class.ID), "primary")
	}

	// Отправляем сообщение с клавиатурой выбора класса
	_, err := vk.MessagesSend(api.Params{
		"user_id":   userID,
		"message":   "Выберите класс:",
		"random_id": 0,
		"keyboard":  keyboard.ToJSON(),
	})
	if err != nil {
		log.Println("Ошибка отправки кнопок выбора класса:", err)
	}
}
