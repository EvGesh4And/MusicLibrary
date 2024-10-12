package models

/*
Song представляет модель песни в базе данных.

@Description Модель, содержащая информацию о песне, включая её название, группу, дату выпуска, текст и ссылку на видео.
@Properties
  - id (integer) "Уникальный идентификатор песни" example(1)
  - group (string) "Название группы, исполнившей песню" example("Группа 1")
  - song (string) "Название песни" example("Песня 1")
  - release_date (string) "Дата выпуска песни" example("2024-01-01")
  - text (string) "Текст песни" example("Это пример текста песни.")
  - link (string) "Ссылка на видео с песней" example("https://www.youtube.com/watch?v=example")
*/
type Song struct {
	ID          uint   `gorm:"primaryKey" json:"id"`      // Уникальный идентификатор песни
	Group       string `gorm:"column:group" json:"group"` // Название группы, исполнившей песню
	Song        string `json:"song"`                      // Название песни
	ReleaseDate string `json:"release_date"`              // Дата выпуска песни
	Text        string `json:"text"`                      // Текст песни
	Link        string `json:"link"`                      // Ссылка на видео с песней
}
