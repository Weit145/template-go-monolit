# template-go-monolit

Go-монолит с HTTP API, PostgreSQL и миграциями Goose. Все команды ниже запускаются из корня репозитория.

## Архитектура

Запрос проходит через `internal/transport/api/v1` (маршруты, middleware, обработчики и DTO), затем через `internal/usecase` к `internal/repositories/storage/postgres`. Бизнес-модели и ошибки находятся в `internal/domain`. Точка входа `cmd/app/main.go` загружает `internal/config`, создаёт подключение к БД и запускает HTTP-сервер.

| Каталог | Назначение |
| --- | --- |
| `cmd/app` | Запуск приложения и сборка зависимостей |
| `internal/config` | Настройки из переменных окружения |
| `internal/domain` | Бизнес-модели и ошибки |
| `internal/usecase` | Сценарии работы и интерфейс репозитория |
| `internal/repositories/storage/postgres` | Работа с PostgreSQL через pgxpool |
| `internal/transport/api/v1` | HTTP API версии v1 |
| `migrations` | SQL-миграции Goose |
| `docs/api/v1` | Спецификация OpenAPI |

## Запуск через Docker

Нужны Docker и Docker Compose. Файл `.env.example` показывает доступные настройки; при необходимости скопируйте его в `.env` и измените значения. Если `.env` отсутствует, Compose использует значения по умолчанию из `docker-compose.yml`.

```bash
make build
```

Команда собирает образы и запускает Compose в текущем терминале. Порядок запуска: PostgreSQL становится доступным → одноразовый контейнер `migrate` применяет миграции Goose → стартует `service`. Если миграция завершится с ошибкой, сервис не запустится. Уже применённые миграции Goose повторно не выполняет.

После первой сборки используйте `make up`. Для запуска в фоне: `docker compose up -d --build`. HTTP API доступен на `localhost:8080`, PostgreSQL — на `localhost:5432`. Внутри сети Compose приложение подключается к БД по имени `db`.

## Команды

| Команда | Что делает |
| --- | --- |
| `make build` | Собирает образы и запускает Compose |
| `make up` | Запускает Compose без принудительной пересборки |
| `make down` | Останавливает и удаляет контейнеры, сохраняя том БД |
| `make delete` | Останавливает Compose **и удаляет том БД со всеми данными** |
| `task goose:install` | Устанавливает Goose в `bin/goose`, если он ещё не установлен |
| `task migrate:up` | Применяет ожидающие миграции вручную |
| `task migrate:down` | Откатывает одну миграцию вручную |
| `task migrate:status` | Показывает состояние миграций |
| `task migrate:create -- имя` | Создаёт SQL-файл миграции |
| `go test ./...` | Проверяет сборку пакетов и запускает тесты |

Задачи `task migrate:*` запускаются на хосте и по умолчанию подключаются к `localhost:5432`. Для другого адреса укажите `PG_DSN`, например `PG_DSN='postgres://postgres:postgres@localhost:5432/postgres' task migrate:status`. При запуске через Compose вызывать `task migrate:up` отдельно не требуется.

Для запуска приложения без Docker сначала запустите PostgreSQL и примените миграции, затем выполните:

```bash
DATABASE_URL='postgres://postgres:postgres@localhost:5432/postgres' go run ./cmd/app
```

Приложение читает переменные окружения процесса; файл `.env` напрямую оно не загружает.

## API

| Метод и путь | Назначение |
| --- | --- |
| `GET /_health` | Проверка готовности приложения и БД |
| `POST /api/v1/users` | Создание пользователя; JSON-тело `{"name":"Alice"}` |
| `GET /swagger/` | Swagger UI |
| `GET /swagger/doc.yaml` | Спецификация OpenAPI |

Исходный документ API: [`docs/api/v1/swagger.yaml`](docs/api/v1/swagger.yaml).

## Конфигурация

Основные переменные: `HTTP_ADDRESS` (по умолчанию `0.0.0.0:8080`), `DATABASE_URL` (в Compose — `postgres://postgres:postgres@db:5432/postgres`), `LOG_LEVEL` и `HTTP_TRUSTED_PROXIES`. HTTP-таймауты задаются переменными `HTTP_READ_HEADER_TIMEOUT`, `HTTP_READ_TIMEOUT`, `HTTP_WRITE_TIMEOUT`, `HTTP_IDLE_TIMEOUT`, `HTTP_SHUTDOWN_TIMEOUT` в формате длительности Go, например `10s`.

Параметры пула PostgreSQL `MAX_CONNS`, `MIN_CONNS`, `MAX_CONN_LIFETIME`, `MAX_CONN_IDLE_TIME` необязательны. Если они не заданы, используются настройки pgxpool по умолчанию.
