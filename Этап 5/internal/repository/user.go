package repository

import (
	"database/sql"
	"expenses/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser сохраняет нового пользователя
func (r *UserRepository) CreateUser(u *models.User) error {
	res, err := r.db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", u.Email, u.PasswordHash)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	u.ID = int(id)
	return nil
}

// GetUserByEmail ищет пользователя для входа
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateSession сохраняет сессию
func (r *UserRepository) CreateSession(sessionID string, userID int, expiresAt interface{}) error {
	_, err := r.db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	return err
}

// GetUserIDBySession проверяет сессию и возвращает ID пользователя
func (r *UserRepository) GetUserIDBySession(sessionID string) (int, error) {
	var userID int
	err := r.db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ? AND expires_at > CURRENT_TIMESTAMP", sessionID).Scan(&userID)
	return userID, err
}

// DeleteSession удаляет сессию (выход из системы)
func (r *UserRepository) DeleteSession(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}