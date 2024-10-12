package models

// SongInput представляет данные, необходимые для создания новой песни.
// @Description Структура, содержащая информацию о песне и группе для создания новой записи в библиотеке.
type SongInput struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

// Song представляет модель песни в базе данных.
// @Description Модель, содержащая информацию о песне, включая её название, группу, дату выпуска, текст и ссылку на видео.
type Song struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Group       string `gorm:"column:group" json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// ResponseAllSongs описывает структуру ответа для получения всех песен.
// @Description Структура ответа для API, возвращающего все песни
type ResponseAllSongs struct {
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Songs []Song `json:"songs"`
}

// ResponseSongVerses описывает структуру ответа для получения куплетов песни.
// @Description Структура ответа для API, возвращающего куплеты песни
type ResponseSongVerses struct {
	Song        string   `json:"song"`
	Group       string   `json:"group"`
	ReleaseDate string   `json:"release_date"`
	Verses      []string `json:"verses"`
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
	Total       int      `json:"total"`
}

// SongDetail представляет данные, полученные из внешнего API.
// @Description Модель, содержащая информацию о дате выпуска песни, тексте и ссылке на видео.
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"` // Дата выпуска песни
	Text        string `json:"text"`        // Текст песни
	Link        string `json:"link"`        // Ссылка на видео с песней
}

// SuccessResponse представляет ответ при успешном удалении песни.
// @Description Структура содержит сообщение о том, что операция выполнена успешно.
type SuccessResponse struct {
	Message string `json:"message"` // Описание успешного выполнения операции
}

// ErrorResponse представляет структуру для ошибок, возвращаемых API.
// @Description Структура, используемая для возврата сообщений об ошибках.
type ErrorResponse struct {
	Error string `json:"error"` // Сообщение об ошибке
}
