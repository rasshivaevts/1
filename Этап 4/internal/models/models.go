package models

import "time"

type Expense struct {
	ID           int
	Amount       float64
	Description  string
	Date         time.Time
	CategoryID   int    // ID выбранной категории
	CategoryName string // Для отображения в таблице (получаем через JOIN)
}

// Новая структура для категорий
type Category struct {
	ID   int
	Name string
}