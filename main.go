package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/events"
	"github.com/SevereCloud/vksdk/v3/object"
)

var confirmationCode string
var vk *api.VK

func main() {
	vk = api.NewVK(os.Getenv("VK_TOKEN"))

	confirmationCode = getConfirmationCode(os.Getenv("VK_GROUP_ID"), os.Getenv("VK_TOKEN"))

	http.HandleFunc("/callback", callbackHandler)

	fmt.Println("Бот запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getConfirmationCode(groupID string, token string) string {
	url := fmt.Sprintf("https://api.vk.com/method/groups.getCallbackConfirmationCode?group_id=%s&access_token=%s&v=5.131", groupID, token)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Ошибка запроса к VK API:", err)
	}
	defer resp.Body.Close()

	var data struct {
		Response struct {
			Code string `json:"code"`
		} `json:"response"`
		Error struct {
			ErrorCode int    `json:"error_code"`
			ErrorMsg  string `json:"error_msg"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal("Ошибка обработки ответа VK API:", err)
	}

	if data.Error.ErrorCode != 0 {
		log.Fatalf("Ошибка VK API: %d, %s", data.Error.ErrorCode, data.Error.ErrorMsg)
	}

	return data.Response.Code
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
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
	case events.EventCallbackQuery: // Обработчик нажатий на кнопки
		var callback events.CallbackQuery
		if err := json.Unmarshal(req.Object, &callback); err != nil {
			log.Println("Ошибка декодирования callback запроса:", err)
			return
		}
		handleCallbackQuery(callback)
	}

	fmt.Fprint(w, "ok")
}

func handleCallbackQuery(callback events.CallbackQuery) {
	payload := callback.Payload

	// Обработка нажатия кнопок
	switch payload {
	case "profile":
		sendMessage(callback.UserID, "Вот ваш профиль.")
	case "stats":
		sendMessage(callback.UserID, "Вот ваша статистика.")
	default:
		sendMessage(callback.UserID, "Неизвестная кнопка.")
	}
}

func handleMessage(msg events.MessageNewObject) {
	userID := msg.Message.PeerID
	text := msg.Message.Text

	switch text {
	case "/start":
		sendMessageWithButtons(userID, "Добро пожаловать в игру! Выберите действие:")
	default:
		sendMessage(userID, "Неизвестная команда. Используйте /start")
	}
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
