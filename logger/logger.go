/*
Package logger предоставляет функциональность для инициализации и настройки логгера.
*/
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

/*
InitLogger инициализирует и настраивает логгер для записи в файл.
Файл создается, если его нет, или добавляются новые записи, если файл существует.
Логгер использует текстовый формат с включенными метками времени.

@Summary Инициализация логгера
@Description Инициализирует и настраивает логгер для записи логов в файл с текстовым форматом и метками времени.
@Tags logger
@Accept  json
@Produce  json
@Success 200 {object} logrus.Logger "Возвращает настроенный экземпляр логгера"
@Failure 500 {string} string "Ошибка инициализации логгера"
@Router /logger/init [get]
*/
func InitLogger() *logrus.Logger {
	logger := logrus.New()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %s", err)
	}

	// Настраиваем логгер для записи логов в файл
	logger.SetOutput(file)

	// Устанавливаем формат логов — текстовый с полными метками времени
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	return logger
}
