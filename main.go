// main.go
package main

import (
	"fitnesshub/server/config"
	"fitnesshub/server/database"
	"fitnesshub/server/router"
	"log"
	"net/http"
)

func main() {
	// Загрузка конфигурации
	config.LoadConfig()

	// Инициализация подключения к базе данных
	database.ConnectDatabase()

	// Создание маршрутизатора с подключением к базе данных
	r := router.NewRouter(database.GetDB())

	// Запуск сервера
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
