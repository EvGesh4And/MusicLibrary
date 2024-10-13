// Модуль MusicLibrary предоставляет API для управления песнями в музыкальной библиотеке.
// Пользователи могут создавать, получать, обновлять и удалять песни.

// @title MusicLibrary API
// @version 1.0
// @description API для управления песнями в библиотеке. Позволяет пользователям получать информацию о песнях, добавлять новые, обновлять и удалять существующие.
// @contact.name Евгений
// @contact.email i@evgesh4.ru
// @host localhost:8080
// @BasePath /
package main

import (
	"MusicLibrary/database"
	_ "MusicLibrary/docs"
	"MusicLibrary/logger"
	"MusicLibrary/routes"
	"os"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Инициализируем логгер
	log := logger.InitLogger()

	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Инициализация базы данных с логгером
	database.Init(log)

	// Настройка маршрутов с логгером
	router := routes.SetupRouter(log)

	// Регистрация Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Infof("Starting server on port %s", port) // Используем логгер для записи информации
	router.Run(":" + port)                        // Запуск сервера на указанном порту.
}
