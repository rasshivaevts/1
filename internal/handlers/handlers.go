package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"expenses_tracker/internal/models"
)

// Handlers — контейнер для обработчиков с общим состоянием
type Handlers struct {
	tmpl  *template.Template
	store *models.ExpenseStore
}

// NewHandlers создаёт обработчики, загружает шаблоны
func NewHandlers(store *models.ExpenseStore) (*Handlers, error) {
	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки шаблонов: %w", err)
	}

	return &Handlers{
		tmpl:  tmpl,
		store: store,
	}, nil
}

// render — вспомогательный метод для рендеринга шаблона
func (h *Handlers) render(w http.ResponseWriter, name string, data map[string]any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("Ошибка рендеринга шаблона %s: %v", name, err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

// ============ Этап 1: базовые маршруты ============

// Home — главная страница
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	h.render(w, "layout", map[string]any{
		"Title":   "Главная",
		"Content": "home",
	})
}

// About — страница описания проекта
func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}
	h.render(w, "layout", map[string]any{
		"Title":   "О проекте",
		"Content": "about",
	})
}

// Ping — health-check, только GET, иначе 405
func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ping" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "pong")
}

// ============ Этап 2: траты ============

// ListExpenses — список всех трат (GET /expenses)
func (h *Handlers) ListExpenses(w http.ResponseWriter, r *http.Request) {
	expenses := h.store.GetAll()
	total := h.store.Total()

	h.render(w, "layout", map[string]any{
		"Title":    "Мои траты",
		"Content":  "expenses",
		"Expenses": expenses,
		"Total":    total,
	})
}

// NewExpense — форма добавления траты (GET /expenses/new)
func (h *Handlers) NewExpense(w http.ResponseWriter, r *http.Request) {
	h.render(w, "layout", map[string]any{
		"Title":   "Добавить трату",
		"Content": "form",
	})
}

// CreateExpense — обработка формы добавления (POST /expenses)
// Использует паттерн Post/Redirect/Get
func (h *Handlers) CreateExpense(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Не удалось разобрать форму", http.StatusBadRequest)
		return
	}

	amountStr := strings.TrimSpace(r.FormValue("amount"))
	description := strings.TrimSpace(r.FormValue("description"))
	date := strings.TrimSpace(r.FormValue("date"))

	// --- Серверная валидация ---
	var errors []string

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		errors = append(errors, "Сумма должна быть положительным числом")
	}
	if description == "" {
		errors = append(errors, "Описание не должно быть пустым")
	}

	// Если есть ошибки — повторный рендер формы с сообщениями
	if len(errors) > 0 {
		h.render(w, "layout", map[string]any{
			"Title":       "Добавить трату",
			"Content":     "form",
			"Errors":      errors,
			"AmountVal":   amountStr,
			"Description": description,
			"DateVal":     date,
		})
		return
	}

	// Сохраняем трату
	h.store.Add(amount, description, date)

	// Post/Redirect/Get
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
