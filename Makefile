# Путь до исполнимого файла
BUILD_DIR=./build

# Загрузка зависимостей и сборка
install-deps:
	@echo "Загрузка зависимостей..."
	@go mod tidy   # Загружает все зависимости, очищает ненужные
	@go mod vendor # Добавляет зависимости в vendor (не обязательно, но полезно)

# Сборка проекта
build:
	@echo "Сборка проекта..."
	@go build -o $(BUILD_DIR)/your_project_name .

# Запуск сервера
run:
	@echo "Запуск сервера..."
	@go run main.go

# Запуск тестов
test:
	@echo "Запуск тестов..."
	@go test -v ./...

# Очистка скомпилированных файлов
clean:
	@echo "Очистка..."
	@rm -rf $(BUILD_DIR)

# Обновление зависимостей и запуск сервера
all: install-deps build run
