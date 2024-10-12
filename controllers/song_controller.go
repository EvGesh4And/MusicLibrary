package controllers

import (
	"MusicLibrary/database"
	"MusicLibrary/models"
	"MusicLibrary/utils"
	"net/http"
	"strconv"
	"strings"

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
// @Param release_date query string false "Дата выпуска"
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
		releaseDate := c.Query("release_date")

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
			query = query.Where("release_date = ?", releaseDate)
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
// @Success 200 {object} models.ResponseSongVerses "Информация о песне и ее куплеты"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена или куплеты закончились"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id} [get]
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

		// Конвертация номера страницы и лимита в числа.
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 1 // По умолчанию один куплет на странице
		}

		// Разделяем текст песни на куплеты, используя \n\n как разделитель.
		verses := strings.Split(song.Text, "\n\n")

		// Вычисляем индексы для пагинации.
		start := (page - 1) * limit
		end := start + limit

		// Проверка границ.
		if start >= len(verses) {
			logger.Warnf("No more verses available for song ID: %s", id)
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "No more verses available"})
			return
		}
		if end > len(verses) {
			end = len(verses)
		}

		// Логируем успешное получение куплетов
		logger.Infof("Returning verses for song ID: %s, page: %d, limit: %d", id, page, limit)

		// Создаем ответ.
		response := models.ResponseSongVerses{
			Song:        song.Song,
			Group:       song.Group,
			ReleaseDate: song.ReleaseDate,
			Verses:      verses[start:end],
			Page:        page,
			Limit:       limit,
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
// @Description Обновляет информацию о песне по её ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.Song true "Обновлённые данные песни"
// @Success 200 {object} models.Song "Обновлённая песня"
// @Failure 400 {object} models.ErrorResponse "Ошибка запроса"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /songs/{id} [put]
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
