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
    return &ExpenseHandler{tmpl: tmpl, repo: repo}
}

// List с поддержкой фильтров
func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
    // Получаем параметры из URL (если их нет, вернется пустая строка/0)
    catStr := r.URL.Query().Get("category")
    fromDate := r.URL.Query().Get("from")
    toDate := r.URL.Query().Get("to")

    catID, _ := strconv.Atoi(catStr)

    // Получаем отфильтрованные траты
    expenses, err := h.repo.GetAllWithFilters(catID, fromDate, toDate)
    if err != nil {
        http.Error(w, "Ошибка БД: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Получаем категории для выпадающего списка
    categories, _ := h.repo.GetAllCategories()

    data := map[string]any{
        "IsForm":      false,
        "Expenses":    expenses,
        "Categories":  categories,
        "SelectedCat": catID,
        "FromDate":    fromDate,
        "ToDate":      toDate,
    }

    if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
    }
}

func (h *ExpenseHandler) NewForm(w http.ResponseWriter, r *http.Request) {
    categories, _ := h.repo.GetAllCategories()
    data := map[string]any{
        "IsForm":     true,
        "Categories": categories,
    }
    if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
    }
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
    dateParsed, _ := time.Parse("2006-01-02", r.FormValue("date"))
    catID, _ := strconv.Atoi(r.FormValue("category_id"))

    expense := models.Expense{
        Amount:      amount,
        Description: r.FormValue("description"),
        Date:        dateParsed,
        CategoryID:  catID,
    }

    _, err := h.repo.Create(&expense)
    if err != nil {
        http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func (h *ExpenseHandler) EditForm(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.PathValue("id"))
    expense, err := h.repo.GetByID(id)
    if errors.Is(err, sql.ErrNoRows) {
        http.NotFound(w, r)
        return
    }
    
    categories, _ := h.repo.GetAllCategories()

    data := map[string]any{
        "IsForm":     true,
        "Expense":    expense,
        "Categories": categories,
    }
    h.tmpl.ExecuteTemplate(w, "layout", data)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.PathValue("id"))
    r.ParseForm()

    amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
    dateParsed, _ := time.Parse("2006-01-02", r.FormValue("date"))
    catID, _ := strconv.Atoi(r.FormValue("category_id"))

    expense := models.Expense{
        ID:          id,
        Amount:      amount,
        Description: r.FormValue("description"),
        Date:        dateParsed,
        CategoryID:  catID,
    }

    if err := h.repo.Update(&expense); err != nil {
        http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.PathValue("id"))
    h.repo.Delete(id)
    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}