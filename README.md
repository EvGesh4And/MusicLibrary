# MusicLibrary API

## Описание

MusicLibrary предоставляет API для управления песнями в музыкальной библиотеке. Пользователи могут создавать, получать, обновлять и удалять песни.

## Дополнительно

Простая реализация внешнее API для получения информации о песне: https://github.com/EvGesh4And/MusicInfo

## Версия

1.0

## Контактная информация

- **Имя**: Евгений
- **Email**: [i@evgesh4.ru](mailto:i@evgesh4.ru)

## Запуск приложения

Для запуска приложения выполните следующие шаги:

1. Клонируйте репозиторий:

    ```bash
    git clone https://github.com/EvGesh4And/MusicLibrary
    cd MusicLibrary
    ```

2. Установите необходимые зависимости:

    ```bash
    go mod tidy
    ```

3. Настройте переменные окружения. Создайте файл `.env` в корне проекта и добавьте следующие параметры:

    ```plaintext
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_username
    DB_NAME=your_database_name
    DB_PASSWORD=your_password
    API_PORT=8080  # Опционально, для настройки порта API
    EXTERNAL_API_URL=http://localhost:9090/info # Указать путь внешнего API для получения дополнительных данных о песне
    ```

4. Запустите приложение:

    ```bash
    go run main.go
    ```

Сервер будет запущен на порту, указанном в переменной окружения `API_PORT` (по умолчанию — 8080).

## Использование API

API MusicLibrary поддерживает следующие эндпоинты:

### Получение списка всех песен
- **URL**: `/songs`
- **Метод**: `GET`
- **Параметры запроса**:
  - `group` (опционально): название группы
  - `song` (опционально): название песни
  - `releaseDate` (опционально): дата выпуска (формат: DD.MM.YYYY)
  - `page` (опционально): номер страницы (по умолчанию: 1)
  - `limit` (опционально): количество записей на странице (по умолчанию: 5)
- **Ответ**:
  - `200 OK`: список песен
  - `400 Bad Request`: ошибка запроса
  - `500 Internal Server Error`: внутренняя ошибка сервера

### Получение информации о песне и её куплетах по ID
- **URL**: `/songs/:id/verses`
- **Метод**: `GET`
- **Параметры**:
  - `id` (обязательный): ID песни
  - `page` (опционально): номер страницы (по умолчанию 1)
  - `limit` (опционально): лимит куплетов на странице (по умолчанию 1)
- **Ответ**:
  - `200 OK`: информация о песне и запрашиваемые куплеты (может вернуть пустой список, если на запрашиваемой странице нет куплетов)
  - `400 Bad Request`: неверный параметр запроса (например, некорректные значения для `page` или `limit`)
  - `404 Not Found`: песня не найдена
  - `500 Internal Server Error`: внутренняя ошибка сервера

### Создание новой песни
- **URL**: `/songs`
- **Метод**: `POST`
- **Тело запроса**: JSON объект с данными песни
- **Ответ**:
  - `200 OK`: созданная песня
  - `400 Bad Request`: ошибка запроса
  - `409 Conflict`: песня уже существует
  - `500 Internal Server Error`: внутренняя ошибка сервера

### Обновление существующей песни
- **URL**: `/songs/:id`
- **Метод**: `PATCH`
- **Параметры**:
  - `id` (обязательный): ID песни
- **Тело запроса**: JSON объект с обновленными данными песни (Изменение ID песни не допускается)
- **Ответ**:
  - `200 OK`: обновленная песня
  - `400 Bad Request`: ошибка запроса
  - `404 Not Found`: песня не найдена
  - `500 Internal Server Error`: внутренняя ошибка сервера

### Удаление песни по ID
- **URL**: `/songs/:id`
- **Метод**: `DELETE`
- **Параметры**:
  - `id` (обязательный): ID песни
- **Ответ**:
  - `200 OK`: песня успешно удалена
  - `404 Not Found`: песня не найдена
  - `500 Internal Server Error`: внутренняя ошибка сервера

## Логирование
Приложение использует logrus для ведения логов. Логи можно настраивать и просматривать для отслеживания работы API и ошибок.

## Swagger Documentation
Документация по API доступна по следующему адресу:
- [Swagger Documentation](http://localhost:8080/swagger/index.html)
