package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var confirmationCode string

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

func main() {
	// Инициализация VK API
	InitVK() // Теперь функция доступна

	// Получаем confirmationCode через GetConfirmationCode
	confirmationCode = getConfirmationCode(os.Getenv("VK_GROUP_ID"), os.Getenv("VK_TOKEN"))

	// Подключаем обработчик Callback API
	http.HandleFunc("/callback", CallbackHandler) // Теперь функция доступна

	fmt.Println("Бот запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
