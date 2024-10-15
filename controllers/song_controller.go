/*
Package controllers содержит обработчики HTTP-запросов для управления песнями в библиотеке.
Каждый контроллер предоставляет функциональность для получения, создания, обновления и удаления песен.
Этот пакет использует модель данных и базу данных, чтобы обрабатывать запросы и возвращать соответствующие ответы клиенту.
*/

package controllers

import (
	"MusicLibrary/database"
	"MusicLibrary/models"
	"MusicLibrary/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus" // Импортируем библиотеку logrus
)

// GetAllSongs возвращает список всех песен с фильтрацией и пагинацией.
// @Summary Получение всех песен
// @Description Возвращает список песен с возможностью фильтрации по группе, названию и дате выпуска, а также поддержкой пагинации.
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Название группы"
// @Param song query string false "Название песни"
// @Param releaseDate query string false "Дата выпуска в формате DD.MM.YYYY"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество песен на странице" default(5)
// @Success 200 {object} models.ResponseAllSongs "Список песен"
// @Failure 400 {object} models.ErrorResponse "Ошибка запроса"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs [get]
func GetAllSongs(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var songs []models.Song
		var total int64

		// Получение параметров фильтрации
		group := c.Query("group")
		song := c.Query("song")
		releaseDate := c.Query("releaseDate")

		// Получение параметров пагинации
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "5")

		// Конвертация параметров пагинации в числа
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			logger.Warnf("Invalid page parameter: %s", page)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid page parameter"})
			return
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt < 1 {
			logger.Warnf("Invalid limit parameter: %s", limit)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid limit parameter"})
			return
		}

		// Фильтрация
		query := database.DB.Model(&models.Song{})
		if group != "" {
			query = query.Where("\"group\" ILIKE ?", "%"+group+"%")
		}
		if song != "" {
			query = query.Where("song ILIKE ?", "%"+song+"%")
		}
		if releaseDate != "" {
			// Проверка формата даты
			parsedDate, err := time.Parse("02.01.2006", releaseDate)
			if err != nil {
				logger.Warnf("Invalid date format for releaseDate: %s, error: %v", releaseDate, err)
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid date format. Expected format: DD.MM.YYYY"})
				return
			}

			// Проверка, что дата не позднее сегодняшнего дня
			if parsedDate.After(time.Now().Truncate(24 * time.Hour)) {
				logger.Warnf("Release date cannot be in the future: %s", releaseDate)
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Release date cannot be in the future"})
				return
			}

			query = query.Where("\"releaseDate\" = ?", releaseDate)
		}

		// Получение общего количества записей
		if err := query.Count(&total).Error; err != nil {
			logger.Errorf("Failed to count songs: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve total count"})
			return
		}

		// Пагинация
		offset := (pageInt - 1) * limitInt
		if err := query.Offset(offset).Limit(limitInt).Find(&songs).Error; err != nil {
			logger.Errorf("Failed to retrieve songs: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve songs"})
			return
		}

		// Логируем полученные данные
		if len(songs) == 0 {
			logger.Warn("No songs found matching the provided filters")
		} else {
			logger.Infof("Retrieved %d songs", len(songs))
		}

		// Подготовка ответа с метаданными
		response := models.ResponseAllSongs{
			Total: total,
			Page:  pageInt,
			Limit: limitInt,
			Songs: songs,
		}
		c.JSON(http.StatusOK, response)
	}
}

// GetSongVerses возвращает куплеты песни по ID.
// @Summary Получение куплетов песни
// @Description Возвращает куплеты песни по указанному ID с поддержкой пагинации.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит на куплеты" default(1)
// @Success 200 {object} models.ResponseSongVerses "Информация о песне и ее куплеты, пустой список, если куплеты отсутствуют на запрашиваемой странице"
// @Failure 400 {object} models.ErrorResponse "Неверный параметр запроса"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id}/verses [get]
func GetSongVerses(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var song models.Song
		id := c.Param("id")

		// Проверяем, существует ли песня с данным ID.
		if err := database.DB.First(&song, id).Error; err != nil {
			logger.Warnf("Song not found with ID: %s", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Song not found"})
			return
		}

		// Получаем номер страницы и лимит из запроса.
		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "1")

		// Конвертация номера страницы и лимита в числа с проверками.
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil || pageInt < 1 {
			logger.Warnf("Invalid page parameter: %s", pageStr)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid page parameter"})
			return
		}

		limitInt, err := strconv.Atoi(limitStr)
		if err != nil || limitInt < 1 {
			logger.Warnf("Invalid limit parameter: %s", limitStr)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid limit parameter"})
			return
		}

		// Разделяем текст песни на куплеты, используя \n\n как разделитель.
		verses := strings.Split(song.Text, "\n\n")

		// Вычисляем индексы для пагинации.
		start := (pageInt - 1) * limitInt
		end := start + limitInt

		// Проверка границ.
		if start >= len(verses) {
			logger.Warnf("No more verses available for song ID: %s", id)
			// Возвращаем пустой список, соответствующий REST.
			c.JSON(http.StatusOK, models.ResponseSongVerses{
				Song:        song.Song,
				Group:       song.Group,
				ReleaseDate: song.ReleaseDate,
				Verses:      []string{},
				Page:        pageInt,
				Limit:       limitInt,
				Total:       len(verses),
			})
			return
		}
		if end > len(verses) {
			end = len(verses)
		}

		// Логируем успешное получение куплетов.
		logger.Infof("Returning verses for song ID: %s, page: %d, limit: %d", id, pageInt, limitInt)

		// Создаем ответ.
		response := models.ResponseSongVerses{
			Song:        song.Song,
			Group:       song.Group,
			ReleaseDate: song.ReleaseDate,
			Verses:      verses[start:end],
			Page:        pageInt,
			Limit:       limitInt,
			Total:       len(verses),
		}

		// Возвращаем запрашиваемые куплеты.
		c.JSON(http.StatusOK, response)
	}
}

// CreateSong добавляет новую песню и обогащает её данные из внешнего API.
// @Summary Создание новой песни
// @Description Добавляет новую песню в библиотеку и обогащает её данные из внешнего API.
// @Tags songs
// @Accept json
// @Produce json
// @Param input body models.SongInput true "Данные песни"
// @Success 200 {object} models.Song "Созданная песня"
// @Failure 400 {object} models.ErrorResponse "Ошибка запроса"
// @Failure 409 {object} models.ErrorResponse "Песня уже существует"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs [post]
func CreateSong(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.SongInput

		// Получаем данные из запроса.
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Warnf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// Проверяем, существует ли песня с таким же названием и группой.
		var existingSong models.Song
		if err := database.DB.Where("song = ? AND \"group\" = ?", input.Song, input.Group).First(&existingSong).Error; err == nil {
			logger.Warnf("Song already exists: %s by %s", input.Song, input.Group)
			c.JSON(http.StatusConflict, models.ErrorResponse{Error: "Song already exists in the library"})
			return
		}

		// Запрос обогащенной информации из внешнего API.
		enrichedData, err := utils.FetchSongDetails(input.Group, input.Song)
		if err != nil {
			logger.Errorf("Failed to fetch song details: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch song details"})
			return
		}

		// Создаем новую песню из данных запроса.
		newSong := models.Song{
			Group:       input.Group,
			Song:        input.Song,
			ReleaseDate: enrichedData.ReleaseDate,
			Text:        enrichedData.Text,
			Link:        enrichedData.Link,
		}

		// Сохранение в базу данных.
		if err := database.DB.Create(&newSong).Error; err != nil {
			logger.Errorf("Failed to save the song: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save the song"})
			return
		}

		logger.Infof("Created song: %s by %s", newSong.Song, newSong.Group)
		c.JSON(http.StatusOK, newSong)
	}
}

// UpdateSong обновляет данные песни по ID.
// @Summary Обновление песни
// @Description Обновляет информацию о песне по её ID. Можно передавать только те поля модели, которые требуется изменить; остальные останутся без изменений. Изменение ID песни не допускается.
// Ожидаемый формат даты: DD.MM.YYYY
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни. Изменение ID не допускается."
// @Param song body models.Song true "Обновлённые данные песни (частичное обновление допускается). Формат даты releaseDate: DD.MM.YYYY"
// @Success 200 {object} models.Song "Обновлённая песня"
// @Failure 400 {object} models.ErrorResponse "Ошибка запроса"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id} [patch]
func UpdateSong(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var song models.Song
		id := c.Param("id")

		// Проверка на существование песни по ID
		if err := database.DB.First(&song, id).Error; err != nil {
			logger.Warnf("Song not found with ID: %s", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Song not found"})
			return
		}

		// Обновление данных
		var input models.Song
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Warnf("Failed to bind JSON for updating song ID: %s, error: %v", id, err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// Проверка на изменение ID
		if input.ID != 0 && input.ID != song.ID {
			logger.Warnf("Attempt to change ID for song ID: %s, new ID: %d", id, input.ID)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Changing the song ID is not allowed"})
			return
		}

		// Проверка поля ReleaseDate на соответствие формату DD.MM.YYYY и на то, что дата не позднее сегодняшнего дня
		if input.ReleaseDate != "" {
			parsedDate, err := time.Parse("02.01.2006", input.ReleaseDate)
			if err != nil {
				logger.Warnf("Invalid date format for song ID: %s, error: %v", id, err)
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid date format. Expected format: DD.MM.YYYY"})
				return
			}

			// Проверка, что дата не позднее сегодняшнего дня
			if parsedDate.After(time.Now().Truncate(24 * time.Hour)) {
				logger.Warnf("Release date cannot be in the future for song ID: %s", id)
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Release date cannot be in the future"})
				return
			}
		}

		// Применение изменений к базе данных
		if err := database.DB.Model(&song).Updates(input).Error; err != nil {
			logger.Errorf("Failed to update song ID: %s, error: %v", id, err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update the song"})
			return
		}

		logger.Infof("Updated song: %s by %s with ID: %s", song.Song, song.Group, id)
		c.JSON(http.StatusOK, song)
	}
}

// DeleteSong удаляет песню по ID.
// @Summary Удаление песни
// @Description Удаляет песню из библиотеки по её ID.
// @Tags songs
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} models.SuccessResponse "Песня успешно удалена"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id} [delete]
func DeleteSong(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var song models.Song
		id := c.Param("id")

		if err := database.DB.First(&song, id).Error; err != nil {
			logger.Warnf("Song not found with ID: %s", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Song not found"})
			return
		}

		if err := database.DB.Delete(&song).Error; err != nil {
			logger.Errorf("Failed to delete song ID: %s, error: %v", id, err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete the song"})
			return
		}

		logger.Infof("Deleted song: %s by %s with ID: %s", song.Song, song.Group, id)
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Song deleted successfully"})
	}
}
