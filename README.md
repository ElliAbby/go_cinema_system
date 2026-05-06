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

## Что делает приложение

- `GET /` — простая стартовая страница.
- `GET /test` — запрос через слой usecase/repository.
- `GET /slow` — имитация долгого запроса.
