package adminScripts

import (
	"database/sql"
	"fmt"
)

type Users struct {
	ID        int
	Name      string
	CreatedAt string
	GroupUser int
}

func ReaderUser() ([]Users, error) {
	dbIp, dbUser, dbPassword, dbName := DatabaseSettingsReader()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbIp, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, created_at, group_user FROM users")
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var users []Users
	for rows.Next() {
		var user Users
		if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.GroupUser); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}

	return users, nil
}

func WriterUser() {

}
