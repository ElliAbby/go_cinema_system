# go_cinema_system

Небольшой HTTP-сервис на Go с вынесенной обработкой заказов в отдельный Kafka worker.

## Требования

- Go 1.22+
- PostgreSQL 16+
- `.env` файл в корне проекта

## Как запустить приложение

### 1. Создай файл `.env`

Скопируйте [`.env.example`](.env.example) в `.env` и заполните реальные значения:

```bash
APP_ADDR=:8080
APP_READ_TIMEOUT=10s
APP_WRITE_TIMEOUT=10s
APP_SHUTDOWN_TIMEOUT=15s
APP_METRICS_ADDR=:8081
APP_ENV=local
RESERVATION_HOLD_DURATION=15m

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cinema
DB_SSL_MODE=disable

JWT_SECRET_KEY=change-me

KAFKA_BROKERS=localhost:9092
KAFKA_BOOKING_PAYMENTS_TOPIC=booking.payments
KAFKA_BOOKING_PAYMENTS_GROUP=cinema-booking-api
```

### 2. Запуск через Docker Compose

Поднять сразу и приложение, и PostgreSQL можно с помощью Docker:

```bash
docker compose up --build
```

После старта:

- приложение будет доступно на `http://localhost:8080`
- PostgreSQL будет доступен на `localhost:5432`
- Kafka будет доступна на `localhost:9092`
- worker заказов будет поднят как отдельный контейнер

#### Тест работы БД

```bash
docker compose exec postgres psql -U postgres -d cinema -c "select * from cinemas;"
```

#### Веб-интерфейсы

После `docker compose up -d --build` доступны:

| Сервис         | URL                           | Описание                                 |
| -------------- | ----------------------------- | ---------------------------------------- |
| Приложение     | http://localhost:8080         | Само приложение                          |
| Метрики        | http://localhost:8080/metrics | Метрики приложения (Prometheus format)   |
| Worker метрики | http://localhost:8081/metrics | Метрики Kafka worker (Prometheus format) |
| Kafka UI       | http://localhost:8090         | Управление топиками и сообщениями        |
| Prometheus     | http://localhost:9090         | Запросы к метрикам (PromQL)              |
| Grafana        | http://localhost:3002         | Интерфейс Grafana дашбордов              |

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

## API Reference

Подробное описание API вынесено в отдельный файл: [docs/api.md](docs/api.md)
