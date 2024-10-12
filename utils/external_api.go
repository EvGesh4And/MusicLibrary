package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

/*
SongDetail представляет данные, полученные из внешнего API.

@Description Модель, содержащая информацию о дате выпуска песни, тексте и ссылке на видео.
@Properties
  - release_date (string) "Дата выпуска песни" example("2024-01-01")
  - text (string) "Текст песни" example("Это пример текста песни.")
  - link (string) "Ссылка на видео" example("https://www.youtube.com/watch?v=example")
*/
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"` // Дата выпуска песни
	Text        string `json:"text"`        // Текст песни
	Link        string `json:"link"`        // Ссылка на видео с песней
}

/*
FetchSongDetails отправляет запрос к внешнему API для получения дополнительных данных о песне.

@Summary Запрос к внешнему API для обогащения данных песни
@Description Эта функция отправляет GET-запрос к внешнему API для получения информации о песне,
включая дату выпуска, текст и ссылку на видео. Пользователи могут передать название группы и название песни
в качестве параметров запроса для получения соответствующих данных.
@Tags utils
@Accept json
@Produce json
@Param group query string true "Название группы"
@Param song query string true "Название песни"
@Success 200 {object} SongDetail "Детали песни, включая дату выпуска, текст и ссылку на видео"
@Failure 400 {string} string "Ошибка получения данных, возможно, некорректные параметры запроса"
@Failure 500 {string} string "Внутренняя ошибка сервера, возникшая при обращении к внешнему API"
@Router /songs/details [get]
*/
func FetchSongDetails(group, song string) (*SongDetail, error) {
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
	var songDetail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, err
	}

	return &songDetail, nil
}
