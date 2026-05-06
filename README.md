# go_cinema_system

Небольшой HTTP-сервис на Go.

## Требования

- Go 1.22+
- PostgreSQL 16+
- `.env` файл в корне проекта

## Как запустить приложение

### 1. Подготовь базу данных PostgreSQL

Запустить Postgres через Docker:

```bash
docker run --name cinema-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cinema -p 5432:5432 -d postgres:16-alpine
```

Если PostgreSQL уже установлен локально, просто убедись, что он запущен и доступен на `localhost:5432`.

#### Тест работы БД

```bash
docker compose exec postgres psql -U postgres -d cinema -c "select * from cinemas;"
```

### 2. Создай файл `.env`

Скопируйте [`.env.example`](.env.example) в `.env` и заполните реальные значения:

```bash
APP_ADDR=:8080
APP_READ_TIMEOUT=10s
APP_WRITE_TIMEOUT=10s
APP_SHUTDOWN_TIMEOUT=15s

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cinema
DB_SSL_MODE=disable
```

### 3. Установи зависимости Go

Из корня проекта:

```bash
go mod download
```

### 4. Запусти приложение

```bash
go run ./cmd/app
```

Если всё настроено правильно, в логах появится сообщение о запуске сервера и успешном подключении к БД.

### 5. Проверь приложение

Открой в браузере или через `curl`:

```bash
curl http://localhost:8080/
curl http://localhost:8080/test
curl http://localhost:8080/slow
```

### 6. Останови сервер

Нажмите `Ctrl + C` в терминале. Приложение выполнит graceful shutdown и корректно закроет HTTP-сервер и соединение с БД.

## Запуск через Docker Compose

Поднять сразу и приложение, и PostgreSQL можно с помощью Docker:

```bash
docker compose up --build
```

После старта:

- приложение будет доступно на `http://localhost:8080`
- PostgreSQL будет доступен на `localhost:5432`

Чтобы остановить и удалить контейнеры:

```bash
docker compose down
```

Чтобы удалить еще и данные БД:

```bash
docker compose down -v
```

### Автосоздание таблиц из deployment/db/init

Файлы `*.sql` из папки `deployment/db/init` автоматически выполняются Postgres-контейнером при первом старте, потому что эта папка подключена в `/docker-entrypoint-initdb.d`.

Важно: скрипты инициализации выполняются только если каталог данных БД пустой (новый volume).

Если изменился SQL и хочется запустить инициализацию заново:

```bash
docker compose down -v
docker compose up --build
```

Важно: если таблицы уже были созданы старой схемой, одного `docker compose up --build` недостаточно. Нужно удалить volume, чтобы Postgres выполнил init-скрипты заново.

## API Documentation

После запуска приложение предоставляет следующие API эндпоинты для работы с данными:

### System Endpoints

- `GET /` — стартовая страница с приветствием
- `GET /health` — проверка здоровья приложения (проверяет БД)

### Movies API

- `GET /movies` — получить все фильмы
- `POST /movies` — создать новый фильм
  ```json
  {
    "title": "Interstellar",
    "duration": 169,
    "rating": "PG-13",
    "description": "Epic sci-fi movie"
  }
  ```
- `GET /movies/{id}` — получить фильм по ID
- `PUT /movies/{id}` — обновить фильм
- `DELETE /movies/{id}` — удалить фильм

### Cinemas API

- `GET /cinemas` — получить все кинотеатры
- `POST /cinemas` — создать новый кинотеатр
  ```json
  {
    "name": "Киносинема",
    "address": "ул. Пушкина, 10"
  }
  ```
- `GET /cinemas/{id}` — получить кинотеатр по ID

### Halls API

- `GET /cinemas/{cinemaId}/halls` — получить все залы в кинотеатре

### Sessions API

- `GET /sessions` — получить все сеансы
- `POST /sessions` — создать новый сеанс
  ```json
  {
    "movie_id": 1,
    "hall_id": 1,
    "start_time": "2026-05-10T18:00:00Z",
    "price_base": 250.0
  }
  ```
- `GET /sessions/{id}` — получить сеанс по ID
- `GET /movies/{movieId}/sessions` — получить все сеансы для фильма

### Примеры использования

#### Создать фильм

```bash
curl -X POST http://localhost:8080/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix",
    "duration": 136,
    "rating": "R",
    "description": "A computer programmer discovers the true nature of his reality"
  }'
```

#### Получить все фильмы

```bash
curl http://localhost:8080/movies
```

#### Получить фильм по ID

```bash
curl http://localhost:8080/movies/1
```

#### Создать кинотеатр

```bash
curl -X POST http://localhost:8080/cinemas \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Киносинема Центр",
    "address": "Красная площадь, 1"
  }'
```

#### Создать сеанс фильма

```bash
curl -X POST http://localhost:8080/sessions \
  -H "Content-Type: application/json" \
  -d '{
    "movie_id": 1,
    "hall_id": 1,
    "start_time": "2026-05-10T19:00:00Z",
    "price_base": 300
  }'
```

#### Проверить здоровье приложения

```bash
curl http://localhost:8080/health
```

### Error Handling

Приложение возвращает структурированные JSON ошибки:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "movie not found",
    "status": 404
  }
}
```

Коды ошибок:

- `NOT_FOUND` (404) — ресурс не найден
- `INVALID_INPUT` (400) — неверные данные
- `CONFLICT` (409) — ресурс уже существует
- `DATABASE_ERROR` (500) — ошибка базы данных
- `INTERNAL_ERROR` (500) — внутренняя ошибка сервера

## Что делает приложение

- `GET /` — простая стартовая страница.
- `GET /test` — запрос через слой usecase/repository.
- `GET /slow` — имитация долгого запроса.
