# Путь до исполнимого файла
BUILD_DIR=./build

# Загрузка зависимостей и сборка
install-deps:
	@echo "Loading dependencies..."
	@go mod tidy   # Загружает все зависимости, очищает ненужные
	@go mod vendor # Добавляет зависимости в vendor (не обязательно, но полезно)

# Сборка проекта
build:
	@echo "Project assembly..."
	@go build -o $(BUILD_DIR)/your_project_name .

# Запуск сервера
run:
	@echo "Starting the server..."
	@go run main.go

# Запуск тестов
test:
	@echo "Running the tests..."
	@go test -v ./...

# Очистка скомпилированных файлов
clean:
	@echo "Clearing..."
	@rm -rf $(BUILD_DIR)

# Обновление зависимостей и запуск сервера
all: install-deps build run
