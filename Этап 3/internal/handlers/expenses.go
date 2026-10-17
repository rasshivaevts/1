package handlers

import (
    "database/sql"
    "errors"
    "html/template"
    "net/http"
    "strconv"
    "time"

    "expenses/internal/models"
    "expenses/internal/repository"
)

type ExpenseHandler struct {
    tmpl *template.Template
    repo *repository.ExpenseRepository
}

func NewExpenseHandler(tmpl *template.Template, repo *repository.ExpenseRepository) *ExpenseHandler {
    return &ExpenseHandler{
        tmpl: tmpl,
        repo: repo,
    }
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
    expenses, err := h.repo.GetAll()
    if err != nil {
        http.Error(w, "Ошибка БД: "+err.Error(), http.StatusInternalServerError)
        return
    }

    data := map[string]any{
        "IsForm":   false,
        "Expenses": expenses,
    }

    if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
    }
}

func (h *ExpenseHandler) NewForm(w http.ResponseWriter, r *http.Request) {
    data := map[string]any{
        "IsForm": true,
    }
    if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
    }
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "bad form", http.StatusBadRequest)
        return
    }

    amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
    if err != nil {
        http.Error(w, "Неверная сумма", http.StatusBadRequest)
        return
    }

    dateParsed, err := time.Parse("2006-01-02", r.FormValue("date"))
    if err != nil {
        http.Error(w, "Неверный формат даты", http.StatusBadRequest)
        return
    }

    expense := models.Expense{
        Amount:      amount,
        Description: r.FormValue("description"),
        Date:        dateParsed,
    }

    _, err = h.repo.Create(&expense)
    if err != nil {
        http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func (h *ExpenseHandler) EditForm(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        http.NotFound(w, r)
        return
    }

    expense, err := h.repo.GetByID(id)
    if errors.Is(err, sql.ErrNoRows) {
        http.NotFound(w, r)
        return
    }
    if err != nil {
        http.Error(w, "Ошибка БД", http.StatusInternalServerError)
        return
    }

    // Важно: передаем и IsForm, и сам объект Expense для заполнения формы
    data := map[string]any{
        "IsForm":  true,
        "Expense": expense, 
    }

    if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
    }
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        http.NotFound(w, r)
        return
    }

    r.ParseForm()

    amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
    dateParsed, err := time.Parse("2006-01-02", r.FormValue("date"))
    if err != nil {
        http.Error(w, "Неверный формат даты", http.StatusBadRequest)
        return
    }

    expense := models.Expense{
        ID:          id,
        Amount:      amount,
        Description: r.FormValue("description"),
        Date:        dateParsed,
    }

    if err := h.repo.Update(&expense); err != nil {
        http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        http.NotFound(w, r)
        return
    }

    if err := h.repo.Delete(id); err != nil {
        http.Error(w, "Ошибка удаления: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}