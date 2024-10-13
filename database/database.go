/*
Package database предоставляет функциональность для инициализации подключения к базе данных PostgreSQL.
Он использует GORM для работы с базой данных и управляет миграцией моделей, загружая параметры конфигурации из файла .env.
Этот пакет обеспечивает глобальный доступ к подключению к базе данных через переменную DB.
*/

package database

import (
	"MusicLibrary/models"
	"fmt"
	"os"

	// Библиотека для работы с файлами .env
	"github.com/sirupsen/logrus" // Логирование
	"gorm.io/driver/postgres"    // Драйвер для PostgreSQL
	"gorm.io/gorm"               // GORM — ORM-библиотека для Go
)

// DB является глобальной переменной для хранения подключения к базе данных
var DB *gorm.DB

// Init инициализирует подключение к базе данных PostgreSQL и выполняет миграцию моделей.
// @Summary Инициализация базы данных
// @Description Устанавливает соединение с PostgreSQL и загружает параметры из .env файла.
// @Tags database
func Init(logger *logrus.Logger) {
	// Формируем строку подключения к базе данных
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))

	// Открываем подключение к базе данных
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Could not connect to the database: %v", err)
	}

	// Включаем стандартное экранирование строк в PostgreSQL
	if err := db.Exec("SET standard_conforming_strings = on;").Error; err != nil {
		logger.Errorf("Failed to set standard_conforming_strings: %v", err)
	} else {
		logger.Infof("Successfully set standard_conforming_strings to on")
	}

	// Проводим автоматическую миграцию модели Song
	if err := db.AutoMigrate(&models.Song{}); err != nil {
		logger.Fatalf("Error during database migration: %v", err)
	} else {
		logger.Infof("Database migration completed successfully")
	}

	// Сохраняем подключение к базе данных в глобальную переменную DB
	DB = db
	logger.Infof("Database connection established successfully")
}
