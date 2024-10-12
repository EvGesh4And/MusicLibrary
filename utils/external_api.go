package utils

import (
	"MusicLibrary/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// FetchSongDetails отправляет запрос к внешнему API для получения дополнительных данных о песне.
// @Summary Запрос к внешнему API для обогащения данных песни
// @Description Эта функция отправляет GET-запрос к внешнему API для получения информации о песне, включая дату выпуска, текст и ссылку на видео.
func FetchSongDetails(group, song string) (*models.SongDetail, error) {
	// Формируем URL запроса к внешнему API с экранированием параметров группы и песни
	apiURL := fmt.Sprintf("%s?group=%s&song=%s", os.Getenv("EXTERNAL_API_URL"), url.QueryEscape(group), url.QueryEscape(song))

	// Отправляем GET-запрос к API
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем успешность запроса по статус-коду
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get song details: %v", resp.Status)
	}

	// Декодируем JSON-ответ в структуру SongDetail
	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, err
	}

	return &songDetail, nil
}
