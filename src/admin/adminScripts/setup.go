package adminScripts

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

var setupCompleted = SetupComleteRead() // Флаг установки

// Подключение к базе данных и выполнение начальной настройки
func setupDatabase(adminUsername, adminPassword, dbIp, dbUser, dbPassword, dbName string) error {
	// Подключение к MySQL
	dsn := fmt.Sprintf("%s:%s@tcp("+dbIp+":3306)/", dbUser, dbPassword)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	defer db.Close()

	// Создание базы данных
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return fmt.Errorf("ошибка при создании базы данных: %w", err)
	}

	// Подключение к созданной базе данных
	dsn = fmt.Sprintf("%s:%s@tcp("+dbIp+":3306)/%s", dbUser, dbPassword, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных %s: %w", dbName, err)
	}
	defer db.Close()

	// Создание таблиц
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    group_user INT
);`)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы users: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    keyName VARCHAR(50) NOT NULL UNIQUE,
    value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы settings: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS news (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title TEXT,
    description TEXT,
    dateNews DATE
);`)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы news: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE user_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    expiration DATETIME NOT NULL
);`)

	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы user_sessions: %w", err)
	}

	// Добавление администратора
	passwordHash := HashPassword(adminPassword)
	_, err = db.Exec("INSERT INTO users (username, password_hash, group_user) VALUES (?, ?, ?)", adminUsername, passwordHash, 0)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении администратора: %w", err)
	}

	return nil
}

// Обработчик установки
func SetupHandler(w http.ResponseWriter, r *http.Request) {
	if setupCompleted {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/setup/setup.html")
		if err != nil {
			http.Error(w, "Не удалось загрузить шаблон", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Получение данных из формы
		adminUsername := r.FormValue("username")
		adminPassword := r.FormValue("password")
		dbIp := r.FormValue("db_ip")
		dbUser := r.FormValue("db_user")
		dbPassword := r.FormValue("db_password")
		dbName := r.FormValue("db_name")

		// Настройка базы данных
		err := setupDatabase(adminUsername, adminPassword, dbIp, dbUser, dbPassword, dbName)
		if err != nil {
			log.Println("Ошибка настройки базы данных:", err)
			http.Error(w, "Не удалось выполнить настройку базы данных", http.StatusInternalServerError)
			return
		}

		// Пометка, что установка завершена
		SetupComleteWrite()

		// Удаление папки setup
		err = os.RemoveAll("templates/test")
		if err != nil {
			log.Println("Ошибка при удалении папки setup:", err)
		} else {
			log.Println("Папка setup успешно удалена.")
		}

		DatabaseSettingsWriter(dbIp, dbUser, dbPassword, dbName)

		// Перенаправление на основной сайт
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
