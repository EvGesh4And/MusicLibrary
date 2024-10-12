package routes

import (
	"MusicLibrary/controllers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*
SetupRouter создает маршруты для приложения и регистрирует обработчики запросов для работы с песнями.

@Summary Настройка маршрутов для работы с песнями
@Description Определение маршрутов для получения, создания, обновления и удаления песен.
@Tags songs
@Accept json
@Produce json
@Router /songs [get]  // Получение списка всех песен
@Router /songs/:id [get]  // Получение информации о песне по ID
@Router /songs [post]  // Создание новой песни
@Router /songs/:id [put]  // Обновление существующей песни по ID
@Router /songs/:id [delete]  // Удаление песни по ID
*/
func SetupRouter(logger *logrus.Logger) *gin.Engine {
	r := gin.Default() // Создаем экземпляр роутера Gin

	// Группа маршрутов для работы с песнями
	songRoutes := r.Group("/songs")
	{
		// GET /songs — маршрут для получения всех песен
		logger.Infof("Setting up route: GET /songs")         // Логируем настройку маршрута
		songRoutes.GET("/", controllers.GetAllSongs(logger)) // Передаем логгер в контроллер

		// GET /songs/:id — маршрут для получения текста песни по ID
		logger.Infof("Setting up route: GET /songs/:id") // Логируем настройку маршрута
		songRoutes.GET("/:id", controllers.GetSongVerses(logger))

		// POST /songs — маршрут для создания новой песни
		logger.Infof("Setting up route: POST /songs") // Логируем настройку маршрута
		songRoutes.POST("/", controllers.CreateSong(logger))

		// PUT /songs/:id — маршрут для обновления данных о песне по ID
		logger.Infof("Setting up route: PUT /songs/:id") // Логируем настройку маршрута
		songRoutes.PUT("/:id", controllers.UpdateSong(logger))

		// DELETE /songs/:id — маршрут для удаления песни по ID
		logger.Infof("Setting up route: DELETE /songs/:id") // Логируем настройку маршрута
		songRoutes.DELETE("/:id", controllers.DeleteSong(logger))
	}

	return r
}
