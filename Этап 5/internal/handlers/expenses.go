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

// List отображает список трат текущего пользователя с учетом фильтров
func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
    // 1. Получаем ID пользователя (благодаря middleware)
    userID := r.Context().Value("userID").(int)

    // 2. Получаем фильтры из URL
    catStr := r.URL.Query().Get("category")
    fromDate := r.URL.Query().Get("from")
    toDate := r.URL.Query().Get("to")
    catID, _ := strconv.Atoi(catStr)

    // 3. Получаем траты (только этого пользователя!)
    expenses, err := h.repo.GetAllWithFilters(userID, catID, fromDate, toDate)
    if err != nil {
        http.Error(w, "Ошибка БД: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // 4. Получаем категории для фильтра
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

// NewForm отображает форму добавления траты
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

// Create сохраняет новую трату за текущим пользователем
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Получаем ID пользователя из контекста
    userID := r.Context().Value("userID").(int)

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

    catID, _ := strconv.Atoi(r.FormValue("category_id"))

    // Создаем модель и ПРИВЯЗЫВАЕМ к пользователю
    expense := models.Expense{
        Amount:      amount,
        Description: r.FormValue("description"),
        Date:        dateParsed,
        CategoryID:  catID,
        UserID:      userID, // <-- Важно: сохраняем ID пользователя
    }

    _, err = h.repo.Create(&expense)
    if err != nil {
        http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// EditForm отображает форму редактирования (только для своей траты)
func (h *ExpenseHandler) EditForm(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)

    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        http.NotFound(w, r)
        return
    }

    // Ищем трату, проверяя, что она принадлежит этому пользователю
    expense, err := h.repo.GetByID(id, userID)
    if errors.Is(err, sql.ErrNoRows) {
        http.NotFound(w, r)
        return
    }
    if err != nil {
        http.Error(w, "Ошибка БД", http.StatusInternalServerError)
        return
    }

    categories, _ := h.repo.GetAllCategories()

    data := map[string]any{
        "IsForm":     true,
        "Expense":    expense,
        "Categories": categories,
    }

    if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
    }
}

// Update сохраняет изменения
func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)

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

    catID, _ := strconv.Atoi(r.FormValue("category_id"))

    expense := models.Expense{
        ID:          id,
        Amount:      amount,
        Description: r.FormValue("description"),
        Date:        dateParsed,
        CategoryID:  catID,
        UserID:      userID, // <-- Передаем ID для проверки при обновлении
    }

    if err := h.repo.Update(&expense); err != nil {
        http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// Delete удаляет трату
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(int)

    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        http.NotFound(w, r)
        return
    }

    // Удаляем, проверяя userID (чтобы нельзя было удалить чужое)
    if err := h.repo.Delete(id, userID); err != nil {
        http.Error(w, "Ошибка удаления: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}