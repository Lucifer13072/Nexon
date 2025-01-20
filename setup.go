package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // Драйвер для MySQL
)

var setupCompleted = false // Флаг установки

// Хэширование пароля
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Подключение к базе данных и выполнение начальной настройки
func setupDatabase(adminUsername, adminPassword, dbUser, dbPassword, dbName string) error {
	// Подключение к MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", dbUser, dbPassword)
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
	dsn = fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPassword, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных %s: %w", dbName, err)
	}
	defer db.Close()

	// Создание таблиц
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS settings (
		id INT AUTO_INCREMENT PRIMARY KEY,
		key_name VARCHAR(50) NOT NULL UNIQUE,
		value TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблиц: %w", err)
	}

	// Добавление администратора
	passwordHash := hashPassword(adminPassword)
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", adminUsername, passwordHash)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении администратора: %w", err)
	}

	return nil
}

// Обработчик установки
func setupHandler(w http.ResponseWriter, r *http.Request) {
	if setupCompleted {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/setup.html")
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
		dbUser := r.FormValue("db_user")
		dbPassword := r.FormValue("db_password")
		dbName := r.FormValue("db_name")

		// Настройка базы данных
		err := setupDatabase(adminUsername, adminPassword, dbUser, dbPassword, dbName)
		if err != nil {
			log.Println("Ошибка настройки базы данных:", err)
			http.Error(w, "Не удалось выполнить настройку базы данных", http.StatusInternalServerError)
			return
		}

		setupCompleted = true
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
