// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Евгений",
            "email": "i@evgesh4.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "Возвращает список песен с возможностью фильтрации по группе, названию и дате выпуска, а также поддержкой пагинации.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Получение всех песен",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Название группы",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название песни",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Дата выпуска в формате DD.MM.YYYY",
                        "name": "releaseDate",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 5,
                        "description": "Количество песен на странице",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список песен",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseAllSongs"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Добавляет новую песню в библиотеку и обогащает её данные из внешнего API.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Создание новой песни",
                "parameters": [
                    {
                        "description": "Данные песни",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SongInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Созданная песня",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Песня уже существует",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/songs/{id}": {
            "delete": {
                "description": "Удаляет песню из библиотеки по её ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Удаление песни",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Песня успешно удалена",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Обновляет информацию о песне по её ID. Можно передавать только те поля модели, которые требуется изменить; остальные останутся без изменений. Изменение ID песни не допускается.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Обновление песни",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни. Изменение ID не допускается.",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Обновлённые данные песни (частичное обновление допускается). Формат даты releaseDate: DD.MM.YYYY",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Обновлённая песня",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/songs/{id}/verses": {
            "get": {
                "description": "Возвращает куплеты песни по указанному ID с поддержкой пагинации.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Получение куплетов песни",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Лимит на куплеты",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Информация о песне и ее куплеты, пустой список, если куплеты отсутствуют на запрашиваемой странице",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseSongVerses"
                        }
                    },
                    "400": {
                        "description": "Неверный параметр запроса",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ErrorResponse": {
            "description": "Структура, используемая для возврата сообщений об ошибках.",
            "type": "object",
            "properties": {
                "error": {
                    "description": "Сообщение об ошибке",
                    "type": "string"
                }
            }
        },
        "models.ResponseAllSongs": {
            "description": "Структура ответа для API, возвращающего все песни",
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "songs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Song"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "models.ResponseSongVerses": {
            "description": "Структура ответа для API, возвращающего куплеты песни",
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "limit": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "releaseDate": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                },
                "verses": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.Song": {
            "description": "Модель, содержащая информацию о песне, включая её название, группу, дату выпуска, текст и ссылку на видео.",
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "models.SongInput": {
            "description": "Структура, содержащая информацию о песне и группе для создания новой записи в библиотеке.",
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                }
            }
        },
        "models.SuccessResponse": {
            "description": "Структура содержит сообщение о том, что операция выполнена успешно.",
            "type": "object",
            "properties": {
                "message": {
                    "description": "Описание успешного выполнения операции",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "MusicLibrary API",
	Description:      "API для управления песнями в библиотеке. Позволяет пользователям получать информацию о песнях, добавлять новые, обновлять и удалять существующие.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
