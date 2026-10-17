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

func (r *ExpenseRepository) Create(e *models.Expense) (int, error) {
	dateStr := e.Date.Format("2006-01-02")
	res, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)",
		e.Amount, e.Description, dateStr,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query("SELECT id, amount, description, date FROM expenses ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		var dateStr string
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &dateStr); err != nil {
			return nil, err
		}
		e.Date, _ = time.Parse("2006-01-02", dateStr)
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *ExpenseRepository) GetByID(id int) (*models.Expense, error) {
	var e models.Expense
	var dateStr string
	err := r.db.QueryRow("SELECT id, amount, description, date FROM expenses WHERE id = ?", id).
		Scan(&e.ID, &e.Amount, &e.Description, &dateStr)
	if err != nil {
		return nil, err
	}
	e.Date, _ = time.Parse("2006-01-02", dateStr)
	return &e, nil
}

func (r *ExpenseRepository) Update(e *models.Expense) error {
	dateStr := e.Date.Format("2006-01-02")
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?",
		e.Amount, e.Description, dateStr, e.ID,
	)
	return err
}

func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}