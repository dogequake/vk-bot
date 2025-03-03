package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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
			sendMessageWithButtons(userID, "Вы не зарегистрированы. Хотите зарегистрироваться?")
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

	log.Println("Получен payload:", payload) // Логируем, что приходит от VK

	// Обработка нажатия на кнопки
	switch payload {
	case "profile":
		sendMessage(userID, "Вот ваш профиль.")
	case "stats":
		sendMessage(userID, "Вот ваша статистика.")
	case "register":
		// Регистрация нового пользователя
		if isUserRegistered(userID) {
			sendMessage(userID, "Вы уже зарегистрированы.")
			return
		}

		registerUser(userID)
		sendMessage(userID, "Вы успешно зарегистрированы! Теперь выберите класс.")
	default:
		sendMessage(userID, "Неизвестная кнопка.")
	}
}
