package main

import (
	"log"
	"os"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/object"
)

var vk *api.VK

// InitVK инициализирует VK API клиент
func InitVK() {
	vkToken := os.Getenv("VK_TOKEN")
	if vkToken == "" {
		log.Fatal("VK_TOKEN не установлен в переменных окружения")
	}
	vk = api.NewVK(vkToken)
}

func sendMessage(userID int, text string) {
	_, err := vk.MessagesSend(api.Params{
		"user_id":   userID,
		"message":   text,
		"random_id": 0,
	})
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

func sendMessageWithButtons(userID int, text string) {
	keyboard := object.NewMessagesKeyboardInline()
	keyboard.AddRow().AddTextButton("Профиль", "profile", "primary")
	keyboard.AddRow().AddTextButton("Статистика", "stats", "secondary")

	_, err := vk.MessagesSend(api.Params{
		"user_id":   userID,
		"message":   text,
		"random_id": 0,
		"keyboard":  keyboard.ToJSON(),
	})
	if err != nil {
		log.Println("Ошибка отправки сообщения с кнопками:", err)
	}
}
