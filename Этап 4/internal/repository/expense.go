package repository

import (
	"database/sql"
	"time"

	"expenses/internal/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// GetAllCategories возвращает список всех категорий
func (r *ExpenseRepository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

// Create сохраняет трату с категорией
func (r *ExpenseRepository) Create(e *models.Expense) (int, error) {
	dateStr := e.Date.Format("2006-01-02")
	res, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date, category_id) VALUES (?, ?, ?, ?)",
		e.Amount, e.Description, dateStr, e.CategoryID,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// GetAllWithFilters возвращает траты с названиями категорий и поддержкой фильтров
func (r *ExpenseRepository) GetAllWithFilters(categoryID int, dateFrom, dateTo string) ([]models.Expense, error) {
	// Базовый запрос с JOIN
	query := `
		SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name 
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE 1=1
	`
	var args []interface{}

	// Динамически добавляем условия фильтрации
	if categoryID > 0 {
		query += " AND e.category_id = ?"
		args = append(args, categoryID)
	}
	if dateFrom != "" {
		query += " AND e.date >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND e.date <= ?"
		args = append(args, dateTo)
	}

	query += " ORDER BY e.date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		var dateStr string
		var catID sql.NullInt64
		var catName sql.NullString

		// Сканируем данные (категория может быть NULL, если не указана)
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &dateStr, &catID, &catName); err != nil {
			return nil, err
		}
		
		e.Date, _ = time.Parse("2006-01-02", dateStr)
		if catID.Valid {
			e.CategoryID = int(catID.Int64)
		}
		if catName.Valid {
			e.CategoryName = catName.String
		}

		result = append(result, e)
	}
	return result, rows.Err()
}

// GetByID теперь тоже достает категорию
func (r *ExpenseRepository) GetByID(id int) (*models.Expense, error) {
	var e models.Expense
	var dateStr string
	var catID sql.NullInt64
	var catName sql.NullString

	err := r.db.QueryRow(`
		SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name 
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.id = ?`, id).
		Scan(&e.ID, &e.Amount, &e.Description, &dateStr, &catID, &catName)
	
	if err != nil {
		return nil, err
	}
	
	e.Date, _ = time.Parse("2006-01-02", dateStr)
	if catID.Valid { e.CategoryID = int(catID.Int64) }
	if catName.Valid { e.CategoryName = catName.String }
	
	return &e, nil
}

// Update обновляет трату вместе с категорией
func (r *ExpenseRepository) Update(e *models.Expense) error {
	dateStr := e.Date.Format("2006-01-02")
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ?",
		e.Amount, e.Description, dateStr, e.CategoryID, e.ID,
	)
	return err
}

func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}