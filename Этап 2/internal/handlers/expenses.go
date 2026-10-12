package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"expenses/internal/models"
)

type ExpenseHandler struct {
	tmpl     *template.Template
	expenses []models.Expense
	nextID   int
}

func NewExpenseHandler(tmpl *template.Template) *ExpenseHandler {
	return &ExpenseHandler{
		tmpl:   tmpl,
		nextID: 1,
	}
}

// List отображает список всех трат
func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"IsForm":   false, // Говорим layout, что это страница списка
		"Expenses": h.expenses,
	}
	if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// NewForm отображает форму добавления траты
func (h *ExpenseHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"IsForm": true, // Говорим layout, что это страница формы
	}
	if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// Create обрабатывает отправку формы и сохраняет трату в памяти
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	h.expenses = append(h.expenses, models.Expense{
		ID:          h.nextID,
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        date,
	})
	h.nextID++

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
