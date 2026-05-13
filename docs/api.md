# API Documentation

Документация описывает HTTP API сервиса кинотеатра и поток обработки заказов. Сервис разделён на два процесса:

- `cinema-service` отвечает за HTTP API
- `order-service` обрабатывает события оплаты из Kafka

## Общие правила

- Base URL для локального запуска: `http://localhost:8080`
- Формат запросов и ответов: `application/json`
- Защищённые методы требуют заголовок `Authorization: Bearer <jwt>`
- При успешной регистрации и входе возвращается JWT токен с временем жизни 1 час

## Аутентификация

### Register

`POST /auth/register`

Создаёт пользователя и сразу возвращает JWT.

Пример запроса:

```json
{
  "email": "user@example.com",
  "password": "password123",
  "phone": "+7 999 123 45 67"
}
```

Пример ответа:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1,
  "email": "user@example.com",
  "expires_at": 1715461234
}
```

### Login

`POST /auth/login`

Аутентифицирует существующего пользователя и возвращает JWT.

Пример запроса:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

## Системные эндпоинты

### Root

`GET /`

Проверка, что сервис запущен.

### Health

`GET /health`

Проверка здоровья приложения и доступности зависимостей.

### Metrics

`GET /metrics`

Метрики в формате Prometheus.

## Public API

### Movies

- `GET /movies` — список фильмов
- `GET /movies/{id}` — фильм по ID
- `POST /movies` — создать фильм (только для авторизавнных пользователей)
- `PUT /movies/{id}` — обновить фильм (только для авторизавнных пользователей)
- `DELETE /movies/{id}` — удалить фильм (только для авторизавнных пользователей)

Пример создания фильма:

```json
{
  "title": "Interstellar",
  "duration": 169,
  "rating": "PG-13",
  "description": "Epic sci-fi movie"
}
```

### Cinemas

- `GET /cinemas` — список кинотеатров
- `GET /cinemas/{id}` — кинотеатр по ID
- `POST /cinemas` — создать кинотеатр (только для авторизавнных пользователей)
- `GET /cinemas/{cinemaId}/halls` — залы кинотеатра

Пример создания кинотеатра:

```json
{
  "name": "Киносинема",
  "address": "ул. Пушкина, 10"
}
```

### Sessions

- `GET /sessions` — список сеансов
- `GET /sessions/{id}` — сеанс по ID
- `GET /movies/{movieId}/sessions` — сеансы конкретного фильма
- `POST /sessions` — создать сеанс (только для авторизавнных пользователей)

Пример создания сеанса:

```json
{
  "movie_id": 1,
  "hall_id": 1,
  "start_time": "2026-05-10T18:00:00Z",
  "price_base": 250.0
}
```

## Protected API

### Bookings

Все методы ниже требуют JWT в заголовке `Authorization`.

- `GET /bookings` — список своих бронирований
- `GET /bookings/{id}` — бронирование по ID
- `POST /bookings` — создать бронь мест
- `POST /bookings/{id}/purchase` — отправить бронирование в Kafka на оплату

Пример создания брони:

```json
{
  "session_id": 1,
  "seats": [1, 2, 3]
}
```

Важно: `POST /bookings/{id}/purchase` возвращает `202 Accepted`. Финальный статус заказа можно проверить через `GET /bookings/{id}`.

### Users

- `GET /users` — список пользователей
- `GET /users/{id}` — профиль текущего пользователя

### Tickets

- `GET /tickets` — список билетов текущего пользователя
- `GET /tickets/{id}` — билет по ID

## Kafka Flow

1. HTTP API создаёт бронирование.
2. `POST /bookings/{id}/purchase` публикует событие `booking.payments`.
3. `order-service` читает событие из Kafka.
4. `order-service` завершает покупку, создаёт билеты и переводит бронь в статус `paid`.

## Ошибки

Ошибки возвращаются в структурированном JSON формате:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "movie not found",
    "status": 404
  }
}
```

Основные коды:

- `NOT_FOUND` — ресурс не найден
- `INVALID_INPUT` — неверные данные запроса
- `CONFLICT` — конфликт состояния или уникальности
- `DATABASE_ERROR` — ошибка доступа к БД
- `INTERNAL_ERROR` — внутренняя ошибка сервера

## Примеры

### Проверить здоровье

```bash
curl http://localhost:8080/health
```

### Получить все фильмы

```bash
curl http://localhost:8080/movies
```

### Получить токен

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```
