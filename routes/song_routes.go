/*
Package routes содержит настройки маршрутов для приложения MusicLibrary.
В этом пакете определяются маршруты для обработки запросов, связанных с песнями,
такие как получение, создание, обновление и удаление песен.
Каждый маршрут регистрирует соответствующий обработчик, обеспечивая необходимую функциональность.
*/

package routes

import (
	"MusicLibrary/controllers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRouter создает маршруты для приложения и регистрирует обработчики запросов для работы с песнями.
// @Summary Настройка маршрутов для работы с песнями
// @Description Определение маршрутов для получения, создания, обновления и удаления песен.
// @Tags songs
func SetupRouter(logger *logrus.Logger) *gin.Engine {
	r := gin.Default() // Создаем экземпляр роутера Gin

	// Группа маршрутов для работы с песнями
	songRoutes := r.Group("/songs")
	{
		// GET /songs — маршрут для получения всех песен
		logger.Infof("Setting up route: GET /songs")
		songRoutes.GET("", controllers.GetAllSongs(logger)) // Убедитесь, что используете "" вместо "/"

		// GET /songs/{id} — маршрут для получения текста песни по ID
		logger.Infof("Setting up route: GET /songs/{id}")
		songRoutes.GET("/:id", controllers.GetSongVerses(logger))

		// POST /songs — маршрут для создания новой песни
		logger.Infof("Setting up route: POST /songs")
		songRoutes.POST("", controllers.CreateSong(logger)) // Убедитесь, что используете "" вместо "/"

		// PATCH /songs/{id} — маршрут для обновления данных о песне по ID
		logger.Infof("Setting up route: PU /songs/{id}")
		songRoutes.PATCH("/:id", controllers.UpdateSong(logger))

		// DELETE /songs/{id} — маршрут для удаления песни по ID
		logger.Infof("Setting up route: DELETE /songs/{id}")
		songRoutes.DELETE("/:id", controllers.DeleteSong(logger))
	}

	return r
}
