package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Структура для получения данных от VK
type VkMessage struct {
	Type   string `json:"type"`
	Object struct {
		Message struct {
			Text   string `json:"text"`
			FromID int    `json:"from_id"`
		} `json:"message"`
	} `json:"object"`
}

// Структура для отправки сообщения
type VkResponse struct {
	UserID   int    `json:"user_id"`
	Message  string `json:"message"`
	RandomID int    `json:"random_id"`
}

const token = "vk1.a.naO-VwqdpBkrsva5qd3zZ_aBLKDZwtXpI7xYMAKjx30zIIk1394tMmq0jTbMA8J42dKzpzYiDSRXPDv0WCE1aKwZxVE45sa8ZO9xv58worZiI59m78x-oVHWxopTShmsagkOdXGq6-5I9nktAW0VpuDoIvXmIA369bQwm8JLOYYDWVnD3LkwxoFPkmOz4rhifGDTF3fWbNDutOs8nnJVPw" // Замени на свой токен

func sendMessage(userID int, text string) {
	url := "https://api.vk.com/method/messages.send"
	data := VkResponse{
		UserID:   userID,
		Message:  text,
		RandomID: 0,
	}

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	q := req.URL.Query()
	q.Add("access_token", token)
	q.Add("v", "5.131")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
		return
	}
	defer resp.Body.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var message VkMessage
	json.Unmarshal(body, &message)

	if message.Type == "message_new" {
		userID := message.Object.Message.FromID
		responseText := fmt.Sprintf("Привет, пользователь %d!", userID)
		sendMessage(userID, responseText)
	}

	fmt.Fprint(w, "ok")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Бот запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
