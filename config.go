package main

import (
	"log"
	"os"
)

func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Ошибка: переменная окружения %s не установлена", key)
	}
	return value
}
