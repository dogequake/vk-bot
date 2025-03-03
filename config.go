package main

import (
	"log"
	"os"

	"github.com/SevereCloud/vksdk/v3/api"
)

func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Ошибка: переменная окружения %s не установлена", key)
	}
	return value
}

func initVK() {
	// Здесь добавь свою логику для инициализации VK
	// Например, присваиваем глобальную переменную vk с помощью getEnvVar
	vkToken := getEnvVar("VK_TOKEN")
	vk = api.NewVK(vkToken)
}
