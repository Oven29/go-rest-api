# Go rest api

## Сервис назначения ревьюеров для Pull Request’ов
Внутри команды требуется единый микросервис, который автоматически назначает ревьюеров на Pull Request’ы (PR), а также позволяет управлять командами и участниками. Взаимодействие происходит исключительно через HTTP API.

### Запуск напрямую:

1. Установка зависимостей
```bash
make install
```
или
```bash
go mod download
```

2. Запуск
```bash
CONFIG="путь до конфига" make run
```
или 
```bash
CONFIG="путь до конфига" go run cmd/go-rest-api/main.go
```
Пример конфига `./config/dev.yaml`

### Сборка

```bash
make build
```
или
```bash
go build -o bin/go-rest-api cmd/go-rest-api/main.go
```

### Запуск через docker

```bash
make docker-up
```
или 
```bash
docker compose up --build -d
```
Для остановки
```bash
make docker-down
```
или 
```bash
docker compose down
```

### Swagger

После запуска доступен сваггер (если не отключен в конфиге):
/swagger/index.html


### Тесты

Для тестов нужен Python c установленными зависимостями `test/requirements.txt`
```bash
make test
```

### Миграции

Для выполнения миграций выполните команду
```bash
DATABASE_URL="postgresql://user:pas@host:port/db_name" make migrate
```
или
```bash
goose -dir migrations postgres "${}" up
```
