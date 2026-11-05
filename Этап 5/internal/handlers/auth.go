package handlers

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"expenses/internal/models"   
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"expenses/internal/repository"
)

type AuthHandler struct {
	tmpl *template.Template
	repo *repository.UserRepository
}

func NewAuthHandler(tmpl *template.Template, repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{tmpl: tmpl, repo: repo}
}

// Показ форм
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "layout", map[string]any{"IsLogin": true})
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "layout", map[string]any{"IsRegister": true})
}

// Обработка регистрации
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	err := h.repo.CreateUser(&models.User{Email: email, PasswordHash: string(hash)})
	if err != nil {
		http.Error(w, "Ошибка: возможно, такой email уже существует", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Обработка входа
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.repo.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// Создаем сессию
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	h.repo.CreateSession(sessionID, user.ID, expiresAt)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: expiresAt,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// Выход
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.repo.DeleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// MIDDLEWARE: проверяет, залогинен ли пользователь
func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := h.repo.GetUserIDBySession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Передаем userID в контекст запроса
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}