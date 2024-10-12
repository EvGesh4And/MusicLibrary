package database

import (
	"MusicLibrary/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"   // Библиотека для работы с файлами .env
	"github.com/sirupsen/logrus" // Импортируем библиотеку logrus для ведения логов
	"gorm.io/driver/postgres"    // Драйвер для работы с PostgreSQL через GORM
	"gorm.io/gorm"               // GORM — ORM-библиотека для Go
)

// DB является глобальной переменной для хранения подключения к базе данных
var DB *gorm.DB

/*
Init инициализирует подключение к базе данных PostgreSQL, загружает параметры из .env файла,
выполняет необходимые настройки и проводит миграцию моделей.
yj
@Summary Инициализация базы данных
@Description Подключение к базе данных PostgreSQL, настройка параметров и миграция моделей
@Tags database
@Accept json
@Produce json
@Success 200 {string} string "Успешное подключение к базе данных"
@Failure 500 {string} string "Ошибка подключения к базе данных"
@Router /database/init [get]

@Param DB_HOST envstring true "Хост базы данных"
@Param DB_PORT envstring true "Порт базы данных"
@Param DB_USER envstring true "Пользователь базы данных"
@Param DB_NAME envstring true "Имя базы данных"
@Param DB_PASSWORD envstring true "Пароль от базы данных"
*/
func Init(logger *logrus.Logger) {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}

	// Формируем строку подключения к базе данных, используя переменные окружения
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))

	// Открываем подключение к базе данных с использованием драйвера PostgreSQL и GORM
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Could not connect to the database: %v", err)
	}

	// Выполняем SQL-запрос для включения режима стандартного экранирования строк в PostgreSQL
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
