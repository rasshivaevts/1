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

// Create теперь привязывает трату к userID
func (r *ExpenseRepository) Create(e *models.Expense) (int, error) {
	dateStr := e.Date.Format("2006-01-02")
	res, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date, category_id, user_id) VALUES (?, ?, ?, ?, ?)",
		e.Amount, e.Description, dateStr, e.CategoryID, e.UserID,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// GetAllWithFilters теперь фильтрует по userID
func (r *ExpenseRepository) GetAllWithFilters(userID, categoryID int, dateFrom, dateTo string) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name, e.user_id
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = ?
	`
	args := []interface{}{userID}

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

		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &dateStr, &catID, &catName, &e.UserID); err != nil {
			return nil, err
		}
		e.Date, _ = time.Parse("2006-01-02", dateStr)
		if catID.Valid { e.CategoryID = int(catID.Int64) }
		if catName.Valid { e.CategoryName = catName.String }

		result = append(result, e)
	}
	return result, rows.Err()
}

// GetByID (также с проверкой userID, чтобы нельзя было редактировать чужое)
func (r *ExpenseRepository) GetByID(id, userID int) (*models.Expense, error) {
	var e models.Expense
	var dateStr string
	var catID sql.NullInt64
	var catName sql.NullString

	err := r.db.QueryRow(`
		SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name, e.user_id
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.id = ? AND e.user_id = ?`, id, userID).
		Scan(&e.ID, &e.Amount, &e.Description, &dateStr, &catID, &catName, &e.UserID)
	
	if err != nil {
		return nil, err
	}
	e.Date, _ = time.Parse("2006-01-02", dateStr)
	if catID.Valid { e.CategoryID = int(catID.Int64) }
	if catName.Valid { e.CategoryName = catName.String }
	
	return &e, nil
}

func (r *ExpenseRepository) Update(e *models.Expense) error {
	dateStr := e.Date.Format("2006-01-02")
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ? AND user_id = ?",
		e.Amount, e.Description, dateStr, e.CategoryID, e.ID, e.UserID,
	)
	return err
}

func (r *ExpenseRepository) Delete(id, userID int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ? AND user_id = ?", id, userID)
	return err
}
