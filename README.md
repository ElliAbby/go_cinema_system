# go_cinema_system

Небольшой HTTP-сервис на Go.

## Требования

- Go 1.22+

## Запуск

Из корня проекта:

```bash
go run ./cmd/app
```

Сервер стартует на `:8080`.

Настройка через `.env`:

```env
APP_ADDR=:8080
APP_READ_TIMEOUT=10s
APP_WRITE_TIMEOUT=10s
APP_SHUTDOWN_TIMEOUT=15s
```

Примеры:

- `GET /`
- `GET /test`
- `GET /slow` (ответ с задержкой)
