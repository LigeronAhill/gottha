# https://just.systems
# Настройки
#

set shell := ["sh", "-c"]
set windows-shell := ["pwsh.exe", "-NoLogo", "-Command"]
set dotenv-load := true
set quiet := true

# Переменные

app_name := "app"
build_dir := "bin"
main_path := "cmd/server/main.go"
db_url := env("APP_DATABASE_URL")
export GOOSE_DRIVER := "postgres"
export GOOSE_DBSTRING := db_url
export GOOSE_MIGRATION_DIR := "./pkg/database/migrations"

# Основные задачи

# Список команд
default:
    @just --list

# Показать помощь
help:
    just --list

# Запуск приложения
[group("app")]
run: gen-style gen-templ
    go run {{ main_path }}

# Сборка приложения
[group("app")]
build: gen-style gen-templ
    mkdir -p {{ build_dir }}
    go build -o {{ build_dir }}/{{ app_name }} {{ main_path }}

# Запуск тестов
[group("app")]
test:
    go test ./...

# Запуск тестов с подробным выводом
[group("app")]
test-v:
    go test -v ./...

# Запуск тестов с детектором гонок
[group("app")]
test-race:
    go test -race ./...

# Запуск бенчмарков
[group("app")]
bench:
    go test -bench=. ./...

# Бенчмарки с профилированием
[group("app")]
bench-mem:
    go test -bench=. -benchmem ./...

# Покрытие кода тестами
[group("app")]
coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    echo "Отчет по покрытию тестами: coverage.html"

# Быстрая проверка покрытия
[group("app")]
cover-func:
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

# Очистка
[group("app")]
clean:
    rm -rf {{ build_dir }}
    rm -f coverage.out coverage.html

# Запуск линтера
[group("app")]
lint:
    go vet ./...
    go fmt ./...

# Полная проверка перед коммитом
[group("app")]
pre-commit: lint test-race

# Установка зависимостей для разработки
[group("app")]
install-deps:
    npm install tailwindcss @tailwindcss/cli
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    go get -tool github.com/a-h/templ/cmd/templ@latest
    go get -tool github.com/air-verse/air@latest
    go get -tool github.com/pressly/goose/v3/cmd/goose@latest
    go mod tidy
    echo "✅ Зависимости установлены"

# Запуск с hot-reload (air)
[group("app")]
watch:
    go tool air

# Генерировать tailwindcss в реальном времени
[group("style")]
watch-style:
    npx @tailwindcss/cli -i ./styles/input.css -o ./public/main.css --watch

# Генерировать tailwindcss
[group("style")]
gen-style:
    npx @tailwindcss/cli -i ./styles/input.css -o ./public/main.css

# Генерировать templ
[group("templ")]
gen-templ:
    go tool templ generate

# Генерировать templ в реальном времени
[group("templ")]
watch-templ:
    go tool templ generate --watch --proxy="http://localhost:42069" --cmd="go run cmd/server/main.go"

# Применить миграции
[group("database")]
migrate: validate
    go tool goose up

# Откатить последнюю миграцию
[group("database")]
down:
    go tool goose down

# Откатить и сразу применить последнюю миграцию
[group("database")]
redo:
    go tool goose redo

# Откатить все миграции
[group("database")]
reset:
    go tool goose reset

# Создать новую миграцию
[group("database")]
create NAME:
    go tool goose create {{ NAME }} sql

# Проверить файлы миграции без их применения
[group("database")]
validate:
    go tool goose validate

# Проверить файлы миграции без их применения
[group("database")]
gen-sql:
    sqlc generate

# Сгенерировать набор иконок
[group("favicon")]
gen-fav SOURCE:
    npx favpie {{ SOURCE }} -o ./public -ap {{ app_name }} -sn {{ app_name }}

# Собрать образ докер
[group("docker")]
docker-build:
    docker build --tag {{ app_name }} .

# Запустить докер контейнер
[group("docker")]
docker-run:
    docker compose up --build -d

# Остановить докер контейнер
[group("docker")]
docker-stop:
    docker compose down -v
    docker system prune -f

# Просмотр логов
[group("docker")]
docker-logs SERVICE="app":
    docker compose logs -f {{ SERVICE }}
