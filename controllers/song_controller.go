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

/*
GetAllSongs возвращает список всех песен с фильтрацией и пагинацией.

@Summary Получение всех песен
@Description Возвращает список песен с возможностью фильтрации по группе, названию и дате выпуска, а также поддержкой пагинации.
@Tags songs
@Accept json
@Produce json
@Param group query string false "Название группы"
@Param song query string false "Название песни"
@Param release_date query string false "Дата выпуска"
@Param page query int false "Номер страницы" default(1)
@Param limit query int false "Количество записей на странице" default(10)
@Success 200 {object} gin.H{"total":int,"page":int,"limit":int,"songs":[]models.Song} "Список песен"
@Failure 400 {string} string "Ошибка запроса"
@Failure 500 {string} string "Внутренняя ошибка сервера"
@Router /songs [get]
*/
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
		limit := c.DefaultQuery("limit", "10")

		// Конвертация параметров пагинации в числа
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			pageInt = 1
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt < 1 {
			limitInt = 10
		}

		// Фильтрация
		query := database.DB.Model(&songs)
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
		query.Count(&total)

		// Пагинация
		offset := (pageInt - 1) * limitInt
		query = query.Offset(offset).Limit(limitInt)

		// Выполнение запроса
		query.Find(&songs)

		// Логируем полученные данные
		if len(songs) == 0 {
			logger.Warn("No songs found matching the provided filters")
		} else {
			logger.Infof("Retrieved %d songs", len(songs))
		}

		// Подготовка ответа с метаданными
		response := gin.H{
			"total": total,
			"page":  pageInt,
			"limit": limitInt,
			"songs": songs,
		}
		c.JSON(http.StatusOK, response)
	}
}

/*
GetSongVerses возвращает куплеты песни по ID.

@Summary Получение куплетов песни
@Description Возвращает куплеты песни по указанному ID с поддержкой пагинации.
@Tags songs
@Accept json
@Produce json
@Param id path int true "ID песни"
@Param page query int false "Номер страницы" default(1)
@Param limit query int false "Лимит на страницу" default(1)
@Success 200 {object} gin.H{"song":string,"group":string,"release_date":string,"verses":[]string,"page":int,"limit":int,"total":int} "Куплеты песни"
@Failure 404 {string} string "Песня не найдена"
@Failure 500 {string} string "Внутренняя ошибка сервера"
@Router /songs/{id} [get]
*/
func GetSongVerses(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var song models.Song
		id := c.Param("id")

		// Проверяем, существует ли песня с данным ID.
		if err := database.DB.First(&song, id).Error; err != nil {
			logger.Warnf("Song not found with ID: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
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
			c.JSON(http.StatusNotFound, gin.H{"error": "No more verses available"})
			return
		}
		if end > len(verses) {
			end = len(verses)
		}

		// Логируем успешное получение куплетов
		logger.Infof("Returning verses for song ID: %s, page: %d, limit: %d", id, page, limit)

		// Возвращаем запрашиваемые куплеты.
		c.JSON(http.StatusOK, gin.H{
			"song":         song.Song,
			"group":        song.Group,
			"release_date": song.ReleaseDate,
			"verses":       verses[start:end],
			"page":         page,
			"limit":        limit,
			"total":        len(verses),
		})
	}
}

/*
CreateSong добавляет новую песню и обогащает её данные из внешнего API.

@Summary Создание новой песни
@Description Добавляет новую песню в библиотеку и обогащает её данные из внешнего API.
@Tags songs
@Accept json
@Produce json
@Param song body models.Song true "Данные песни"
@Success 200 {object} models.Song "Созданная песня"
@Failure 400 {string} string "Ошибка запроса"
@Failure 409 {string} string "Песня уже существует"
@Failure 500 {string} string "Внутренняя ошибка сервера"
@Router /songs [post]
*/
func CreateSong(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.Song

		// Получаем данные из запроса.
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Warnf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверяем, существует ли песня с таким же названием и группой.
		var existingSong models.Song
		if err := database.DB.Where("song = ? AND \"group\" = ?", input.Song, input.Group).First(&existingSong).Error; err == nil {
			logger.Warnf("Song already exists: %s by %s", input.Song, input.Group)
			c.JSON(http.StatusConflict, gin.H{"error": "Song already exists in the library"})
			return
		}

		// Запрос обогащенной информации из внешнего API.
		enrichedData, err := utils.FetchSongDetails(input.Group, input.Song)
		if err != nil {
			logger.Errorf("Failed to fetch song details: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song details"})
			return
		}

		// Обновление данных песни.
		input.ReleaseDate = enrichedData.ReleaseDate
		input.Text = enrichedData.Text
		input.Link = enrichedData.Link

		// Сохранение в базу данных.
		if err := database.DB.Create(&input).Error; err != nil {
			logger.Errorf("Failed to save the song: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the song"})
			return
		}

		logger.Infof("Created song: %s by %s", input.Song, input.Group)
		c.JSON(http.StatusOK, input)
	}
}

/*
UpdateSong обновляет данные песни по ID.

@Summary Обновление песни
@Description Обновляет информацию о песне по её ID.
@Tags songs
@Accept json
@Produce json
@Param id path int true "ID песни"
@Param song body models.Song true "Обновлённые данные песни"
@Success 200 {object} models.Song "Обновлённая песня"
@Failure 400 {string} string "Ошибка запроса"
@Failure 404 {string} string "Песня не найдена"
@Failure 500 {string} string "Внутренняя ошибка сервера"
@Router /songs/{id} [put]
*/
func UpdateSong(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var song models.Song
		id := c.Param("id")

		if err := database.DB.First(&song, id).Error; err != nil {
			logger.Warnf("Song not found with ID: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}

		// Обновление данных.
		var input models.Song
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Warnf("Failed to bind JSON for updating song ID: %s, error: %v", id, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.DB.Model(&song).Updates(input).Error; err != nil {
			logger.Errorf("Failed to update song ID: %s, error: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the song"})
			return
		}

		logger.Infof("Updated song: %s by %s with ID: %s", song.Song, song.Group, id)
		c.JSON(http.StatusOK, song)
	}
}

/*
DeleteSong удаляет песню по ID.

@Summary Удаление песни
@Description Удаляет песню из библиотеки по её ID.
@Tags songs
@Produce json
@Param id path int true "ID песни"
@Success 200 {object} gin.H{"data":bool} "Песня успешно удалена"
@Failure 404 {string} string "Песня не найдена"
@Failure 500 {string} string "Внутренняя ошибка сервера"
@Router /songs/{id} [delete]
*/
func DeleteSong(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var song models.Song
		id := c.Param("id")

		if err := database.DB.First(&song, id).Error; err != nil {
			logger.Warnf("Song not found with ID: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}

		if err := database.DB.Delete(&song).Error; err != nil {
			logger.Errorf("Failed to delete song ID: %s, error: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the song"})
			return
		}

		logger.Infof("Deleted song: %s by %s with ID: %s", song.Song, song.Group, id)
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}
